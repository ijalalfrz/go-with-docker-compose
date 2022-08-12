package model

import "github.com/golang-jwt/jwt/v4"

// UserBeaerClaims is a model.
type UserBeaerClaims struct {
	jwt.StandardClaims
	UserName string `json:"user_name"`
}

// UserSuccessSignIn is a model.
type UserSuccessSignIn struct {
	Name  string  `json:"name"`
	Token *string `json:"token,omitempty"`
}

// UserLoginParams is a model.
type UserLoginParams struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
