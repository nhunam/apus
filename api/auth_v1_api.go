package api

import (
	"apuscorp.com/core/auth-service/dto/apperror"
	"apuscorp.com/core/auth-service/dto/response"
	"apuscorp.com/core/auth-service/router"
	"apuscorp.com/core/auth-service/util"
	"apuscorp.com/core/auth-service/util/errorcode"
	jwt "github.com/appleboy/gin-jwt/v2"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/apus-backend/base-service/config"
	"gitlab.com/apus-backend/base-service/rabbitmq"
	"net/http"
	"time"
)

type AuthV1Api struct {
	Router   *router.Router
	RabbitMQ *rabbitmq.RabbitMQ
	Config   config.Config
}

func (r *AuthV1Api) GetRoutingKey(routingKeys []string, data interface{}) string {
	return routingKeys[0]
}

// @Summary Login với username/password
// @Description Login với username/password, see: github.com/appleboy/gin-jwt/v2, auth_jwt, function LoginHandler
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.UserLoginDto true "JSON body"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response "{"error_code": "<Mã lỗi>", "error_msg": "<Nội dung lỗi>"}"
// @Router /api/v1/auth/login [post]
func (r *AuthV1Api) Login(c *gin.Context) {
	if r.Router.AuthMiddleware.Authenticator == nil {
		panic(jwt.ErrMissingAuthenticatorFunc)
	}

	data, err := r.Router.AuthMiddleware.Authenticator(c)
	util.Must(err)

	// Create the token
	token := jwtGo.New(jwtGo.GetSigningMethod(r.Router.AuthMiddleware.SigningAlgorithm))
	claims := token.Claims.(jwtGo.MapClaims)

	if r.Router.AuthMiddleware.PayloadFunc != nil {
		for key, value := range r.Router.AuthMiddleware.PayloadFunc(data) {
			claims[key] = value
		}
	}

	// HaiMT: nếu login thường thì dùng Timeout. Nếu login có remember me thì dùng Long expire time.
	var expire time.Time
	if loginDto, ok := data.(*response.UserLoginRes); ok && loginDto.RememberMe != nil && *loginDto.RememberMe {
		expire = r.Router.AuthMiddleware.TimeFunc().Add(r.Router.AuthMiddleware.Timeout)
	} else {
		expire = r.Router.AuthMiddleware.TimeFunc().Add(r.Router.LongRefreshExpTime)
	}
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = r.Router.AuthMiddleware.TimeFunc().Unix()
	tokenString, err := r.signedString(token)
	util.Must(err)

	// set cookie
	if r.Router.AuthMiddleware.SendCookie {
		expireCookie := r.Router.AuthMiddleware.TimeFunc().Add(r.Router.AuthMiddleware.CookieMaxAge)
		maxage := int(expireCookie.Unix() - r.Router.AuthMiddleware.TimeFunc().Unix())

		if r.Router.AuthMiddleware.CookieSameSite != 0 {
			c.SetSameSite(r.Router.AuthMiddleware.CookieSameSite)
		}

		c.SetCookie(
			r.Router.AuthMiddleware.CookieName,
			tokenString,
			maxage,
			"/",
			r.Router.AuthMiddleware.CookieDomain,
			r.Router.AuthMiddleware.SecureCookie,
			r.Router.AuthMiddleware.CookieHTTPOnly,
		)
	}

	r.Router.AuthMiddleware.LoginResponse(c, http.StatusOK, tokenString, expire)
}

func (r *AuthV1Api) unauthorized(c *gin.Context, code int, err error) {
	c.Header("WWW-Authenticate", "JWT realm="+r.Router.AuthMiddleware.Realm)
	if !r.Router.AuthMiddleware.DisabledAbort {
		c.Abort()
	}

	if appError, ok := err.(apperror.AppError); ok {
		c.AbortWithStatusJSON(code, response.Response{
			ErrorCode:    appError.ErrorCode,
			ErrorMessage: appError.ErrorMessage,
		})
	} else {
		c.AbortWithStatusJSON(code, response.Response{
			ErrorCode:    errorcode.AUT1006,
			ErrorMessage: err.Error(),
		})
	}
}

func (r *AuthV1Api) signedString(token *jwtGo.Token) (string, error) {
	var tokenString string
	var err error
	if r.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(r.Router.PrivateKey)
	} else {
		tokenString, err = token.SignedString(r.Router.AuthMiddleware.Key)
	}
	return tokenString, err
}

func (r *AuthV1Api) usingPublicKeyAlgo() bool {
	switch r.Router.AuthMiddleware.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// @Summary Refresh access token
// @Description Refresh access token
// @Tags auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response "{"error_code": "<Mã lỗi>", "error_msg": "<Nội dung lỗi>"}"
// @Router /api/v1/auth/refresh_token [post]
func (r *AuthV1Api) RefreshToken(c *gin.Context) {
	// NOTE empty function in order to render swagger
}
