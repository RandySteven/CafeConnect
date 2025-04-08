package jwt_client

import (
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

type JWTClaim struct {
	UserID   uint64
	RoleID   []uint64
	IsVerify bool
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User, roleUser []*models.RoleUser) string {

	claims := &JWTClaim{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Applications",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}
	for _, role := range roleUser {
		claims.RoleID = append(claims.RoleID, role.RoleID)
	}
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(JwtKey)
	if err != nil {
		return ""
	}
	return token
}
