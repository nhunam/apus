package router

import (
	"apuscorp.com/core/auth-service/dao"
	docs "apuscorp.com/core/auth-service/docs"
	"apuscorp.com/core/auth-service/dto/apperror"
	"apuscorp.com/core/auth-service/dto/request"
	"apuscorp.com/core/auth-service/dto/response"
	"apuscorp.com/core/auth-service/i18n"
	"apuscorp.com/core/auth-service/util"
	"apuscorp.com/core/auth-service/util/errorcode"
	"crypto/rsa"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gitlab.com/apus-backend/base-service/config"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	reset            = "\033[0m"
	JWT_IDENTITY_KEY = "id"
	JWT_USER_ID      = "user_id"
	JWT_COMPANY_ID   = "company_id"
	JWT_ENVIRONMENT  = "env"
)

type Router struct {
	Engine             *gin.Engine
	AuthMiddleware     *jwt.GinJWTMiddleware
	I18n               *i18n.I18n
	DfRedisStore       *persistence.RedisStore
	PrivateKey         *rsa.PrivateKey
	UserCredentialsDao *dao.UserCredentialsDao
	LongRefreshExpTime time.Duration
}

func NewRouter(c config.Config, userCredentialsDao *dao.UserCredentialsDao, i18n *i18n.I18n) (*Router, error) {
	e := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	//e.Use(gin.Logger())
	e.Use(logger.SetLogger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	e.Use(CustomRecovery(CustomRecoveryFunc(i18n)))

	// CORS
	e.Use(initCorsMiddleware())

	// the jwt middleware
	authMiddleware, err := initJwtMiddleware(c, userCredentialsDao)

	if err != nil {
		log.Fatal().Err(err).Msg("JWT Error:" + err.Error())
		return nil, err
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal().Err(errInit).Msg("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		return nil, errInit
	}

	refreshExpTime, err := time.ParseDuration(c.Jwt.RefreshExpTime)
	if err != nil {
		return nil, err
	}
	res := &Router{
		Engine:             e,
		AuthMiddleware:     authMiddleware,
		I18n:               i18n,
		DfRedisStore:       persistence.NewRedisCache("localhost:6379", "", time.Hour*2), // cache middleware
		UserCredentialsDao: userCredentialsDao,
		LongRefreshExpTime: refreshExpTime,
		PrivateKey:         privateKey(authMiddleware),
	}

	return res, nil
}

func CustomRecoveryFunc(i18n *i18n.I18n) func(c *gin.Context, recovered interface{}) {
	return func(c *gin.Context, recovered interface{}) {
		if v, ok := recovered.(apperror.AppError); ok {
			util.HandleAppError(c, i18n, v.Status, v)
		} else if v, ok := recovered.(*apperror.AppError); ok {
			util.HandleAppError(c, i18n, (*v).Status, *v)
		} else if v, ok := recovered.(validator.ValidationErrors); ok {
			util.HandleValidationErrors(c, i18n, http.StatusBadRequest, v)
		} else if v, ok := recovered.(*validator.ValidationErrors); ok {
			util.HandleValidationErrors(c, i18n, http.StatusBadRequest, *v)
		} else if v, ok := recovered.(error); ok {
			util.HandleError(c, i18n, http.StatusInternalServerError, v)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func initJwtMiddleware(c config.Config, userCredentialsDao *dao.UserCredentialsDao) (*jwt.GinJWTMiddleware, error) {
	expiredTime, err := time.ParseDuration(c.Jwt.ExpiredTime)
	if err != nil {
		return nil, err
	}
	refreshExpTime, err := time.ParseDuration(c.Jwt.RefreshExpTime)
	if err != nil {
		return nil, err
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            c.Jwt.Realm,
		SigningAlgorithm: c.Jwt.SigningAlg,
		Key:              []byte(c.Jwt.Secret),
		Timeout:          expiredTime,
		MaxRefresh:       refreshExpTime,
		IdentityKey:      JWT_IDENTITY_KEY,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginDto request.UserLoginDto
			if err := c.ShouldBind(&loginDto); err != nil {
				return "", err
			}

			return userCredentialsDao.Login(c, &loginDto)
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*response.UserLoginRes); ok {
				return jwt.MapClaims{
					JWT_IDENTITY_KEY: uuid.NewString(),
					JWT_USER_ID:      v.UserId,
					JWT_COMPANY_ID:   v.CompanyId,
					JWT_ENVIRONMENT:  c.Env,
				}
			} else {
				util.CheckError(apperror.NewAppError(errorcode.AUT0004, "", "", map[string]string{
					"Type": "response.UserLoginRes",
				}))
			}
			log.Fatal().Msgf("data must be instance of model.User: %s", data)
			return nil
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Info().Msgf("IdentityHandler, userId: %s", claims[JWT_USER_ID])
			// TODO return identity key
			return claims[JWT_USER_ID]
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Info().Msgf("IdentityHandler data: %s", data)
			// TODO authorize function/APIs
			// TODO test , always return true in order to authorize function
			return true
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}

func initCorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	})
}

func privateKey(privKeyFile *jwt.GinJWTMiddleware) *rsa.PrivateKey {
	value := util.GetUnexportedField(reflect.ValueOf(privKeyFile).Elem().FieldByName("privKey"))
	if value == nil {
		return nil
	}
	return value.(*rsa.PrivateKey)
}

func (r *Router) InitSwagger(c config.Config) {
	docs.SwaggerInfo.Host = fmt.Sprintf("%v:%v", c.Server.Host, c.Server.Port)
	url := ginSwagger.URL(fmt.Sprintf("http://%v:%v/swagger/doc.json", c.Server.Host, c.Server.Port))
	r.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
