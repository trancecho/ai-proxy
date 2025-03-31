package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

var MundoSecret []byte
var OffercatSecret []byte

func InitSecret() {
	MundoSecret = []byte(viper.GetString("jwt.jwt_sec") + "mundo")
	//log.Println("MundoSecret: ", MundoSecret)
	//OffercatSecret = []byte(config.GetConfig().Jwt.Offercat + "offercat")
	//log.Println("OffercatSecret: ", OffercatSecret)
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	// StandardClaims 已经弃用，使用 RegisteredClaims
	jwt.RegisteredClaims
}

// 生成 JWT Token
//func GenerateToken(userID int64, username, role string, service string) (string, error) {
//	now := time.Now()
//	expireTime := now.Add(24 * 7 * time.Hour) // Token 有效期 一周
//
//	claims := Claims{
//		UserID:   userID,
//		Username: username,
//		Role:     role,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expireTime), // 转换为 *NumericDate
//			IssuedAt:  jwt.NewNumericDate(now),        // 转换为 *NumericDate
//			Issuer:    viper.GetString("jwt.issuer"),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	if service == "mundo" {
//		return token.SignedString(MundoSecret)
//	} else {
//		return token.SignedString(OffercatSecret)
//	}
//}

// 验证 JWT Token
func ParseToken(service string, tokenString string) (*Claims, error) {
	if service == "mundo" {
		return parseToken(MundoSecret, tokenString)
	} else {
		return parseToken(OffercatSecret, tokenString)
	}
}
func parseToken(secret []byte, tokenString string) (*Claims, error) {
	log.Println("tokenString: ", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println("Token 解析失败:", err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("Token 无效")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// JWTAuthMiddleware 是一个Gin中间件，用于验证JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		log.Println("Request Header: Authorization:", token) // 在此处打印请求的 Authorization 头

		// 如果请求头里没有 Token，则尝试从 URL 查询参数获取
		if token == "" {
			token = c.Query("token")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "需要 Token"})
				c.Abort()
				return
			}
		} else {
			// 处理 "Bearer token_value" 格式
			parts := strings.Split(token, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token 格式"})
				c.Abort()
				return
			}
		}

		// 解析 Token
		claims, err := ParseToken("mundo", token) // 你的 Token 默认使用 "mundo"
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的或过期的 Token"})
			c.Abort()
			return
		}

		// 将用户信息存储在 Gin 的上下文中，后续处理可以直接获取
		c.Set("uid", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 继续处理请求
		c.Next()
	}
}
