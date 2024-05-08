package model

import "github.com/dgrijalva/jwt-go"

const PathUserCreate = "/auth_v1.AuthV1/Create"

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int    `json:"role"`
	Email    string `json:"email"`
}
