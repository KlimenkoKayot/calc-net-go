package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenManager struct {
	jwtSecretKey           []byte
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

func NewTokenManager(jwtSecretKey string, accessTokenExpiration, refreshTokenExpiration time.Duration) (*TokenManager, error) {
	return &TokenManager{
		jwtSecretKey:           []byte(jwtSecretKey),
		accessTokenExpiration:  accessTokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

func (tm *TokenManager) NewAccessToken(values map[string]interface{}) (string, error) {
	payload := jwt.MapClaims{}
	payload = values
	payload["exp"] = time.Now().Add(tm.accessTokenExpiration).Unix()
	payload["ctd"] = time.Now().Unix()

	tokenData, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(tm.jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenData, nil
}

func (tm *TokenManager) NewRefreshToken(values map[string]interface{}) (string, error) {
	payload := jwt.MapClaims{}
	payload = values
	payload["exp"] = time.Now().Add(tm.refreshTokenExpiration).Unix()
	payload["ctd"] = time.Now().Unix()

	tokenData, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(tm.jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenData, nil
}

func (tm *TokenManager) ValidateTokenExpiration(token string) (bool, error) {
	valid, err := tm.ValidateToken(token)
	if !valid || err != nil {
		return false, err
	}
	claims, err := tm.ParseWithClaims(token)
	if err != nil {
		return false, err
	}
	expTime := (*claims)["exp"].(time.Time)
	expired := time.Now().After(expTime)
	// если истек, то невалидный
	return !expired, nil
}

func (tm *TokenManager) ValidateToken(tokenString string) (bool, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tm.jwtSecretKey, nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func (tm *TokenManager) ParseWithClaims(tokenString string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tm.jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

/*
Возвращает пару из access (1) и refresh (2) токенов, ошибку (3), если возникла.
*/
func (tm *TokenManager) UpdateTokenPair(refreshToken string) (string, string, error) {
	valid, err := tm.ValidateTokenExpiration(refreshToken)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", ErrNotValidToken
	}

	claims, err := tm.ParseWithClaims(refreshToken)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = tm.NewRefreshToken(*claims)
	if err != nil {
		return "", "", err
	}

	accessToken, err := tm.NewAccessToken(*claims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

/*
Возвращает пару из access (1) и refresh (2) токенов, ошибку (3), если возникла.
*/
func (tm *TokenManager) NewTokenPair(accessData, refreshData map[string]interface{}) (string, string, error) {
	refreshToken, err := tm.NewRefreshToken(refreshData)
	if err != nil {
		return "", "", err
	}

	accessToken, err := tm.NewAccessToken(accessData)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
