module apuscorp.com/core/auth-service

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/logger v0.0.2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/cache/v8 v8.3.1
	github.com/go-redis/redis/v8 v8.6.0
	github.com/go-resty/resty/v2 v2.5.0
	github.com/google/uuid v1.2.0
	github.com/google/wire v0.5.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/rs/zerolog v1.20.0
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	gitlab.com/apus-backend/base-service v1.0.1
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/text v0.3.5
	golang.org/x/tools v0.1.0 // indirect
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

//replace gitlab.com/apus-backend/base-service v1.0.0 => ../base-service
