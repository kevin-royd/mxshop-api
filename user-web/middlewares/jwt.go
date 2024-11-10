package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 token 信息
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// 解析 token
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, TokenExpired) {
				zap.S().Infof("expired:%s", TokenExpired)
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "未登录",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConf.JWTInfo.SigningKey),
	}
}

// 创建一个 token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().In(global.TimeZone).Add(1 * time.Hour))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, parseTokenError(err)
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		// 校验过期时间是否已过
		if claims.ExpiresAt.Before(time.Now().In(global.TimeZone)) {
			return nil, TokenExpired
		}
		return claims, nil
	}
	return nil, TokenInvalid
}

// 刷新 token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().In(global.TimeZone).Add(1 * time.Hour)) // 更新过期时间
		return j.CreateToken(*claims)
	}

	return "", TokenInvalid
}

// parseTokenError 解析 Token 错误
func parseTokenError(err error) error {
	if err != nil {
		// 如果错误是 jwt.ErrSignatureInvalid，表示签名验证失败
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return TokenInvalid
		}

		// 如果错误是 jwt.ErrTokenExpired，表示 Token 已过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			return TokenExpired
		}

		// 如果错误是 jwt.ErrTokenNotValidYet，表示 Token 尚未生效
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return TokenNotValidYet
		}
	}

	return TokenInvalid
}
