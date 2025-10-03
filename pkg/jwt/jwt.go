package jwt

import (
	"errors"
	"time"

	"go-my-blog/config" // 引入配置包，读取 JWT 密钥和过期时间

	"github.com/golang-jwt/jwt/v5"
)

// 自定义 Claims（包含用户 ID 和标准声明）
type CustomClaims struct {
	UserID               uint   `json:"user_id"`  // 存储用户 ID
	Username             string `json:"username"` // 存储用户名
	jwt.RegisteredClaims        // 标准声明（包含过期时间等）
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string) (string, int64, error) {
	// 1. 从配置中获取 JWT 密钥和过期时间（建议在 config/app.yaml 中配置）
	jwtConf := config.Conf.JWT
	secret := []byte(jwtConf.Secret) // 密钥（生产环境需复杂且保密）
	expireTime := time.Now().Add(time.Duration(jwtConf.ExpireHour) * time.Hour)
	expireAt := expireTime.Unix()

	// 2. 设置 Claims（包含用户 ID 和过期时间）
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间（立即生效）
		},
	}

	// 3. 使用 HS256 算法签名生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	return tokenString, expireAt, err
}

// VerifyToken 验证 Token 并返回用户 ID
func VerifyToken(tokenString string) (*CustomClaims, error) {
	// 1. 从配置中获取 JWT 密钥
	secret := []byte(config.Conf.JWT.Secret)
	if len(secret) == 0 {
		return nil, errors.New("JWT 密钥未配置")
	}

	// 2. 解析 Token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{}, // 传入自定义 Claims 指针，用于接收解析结果
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法是否为预期的 HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("不支持的签名算法")
			}
			return secret, nil
		},
	)
	if err != nil {
		return nil, errors.New("token 解析失败：" + err.Error())
	}

	// 3. 验证 Token 有效性（是否过期、签名是否正确）
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil // 验证通过，返回用户 ID
	}

	return nil, errors.New("token 无效或已过期")
}
