package jwt

import (
	"net/http"

	"github.com/klimenkokayot/avito-go/libs/jwt"
)

type TokenMiddleware struct {
	tokenManager *jwt.TokenManager
}

func NewTokenMiddleware(tokenManager *jwt.TokenManager) (*TokenMiddleware, error) {
	return &TokenMiddleware{
		tokenManager: tokenManager,
	}, nil
}

func (tm *TokenMiddleware) updateTokenPair(r *http.Request) error {
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		return ErrAccessTokenNotFound
	}
	accessToken := accessTokenCookie.Value

	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return ErrRefreshTokenNotFound
	}
	refreshToken := refreshTokenCookie.Value

	accessToken, refreshToken, err = tm.tokenManager.UpdateTokenPair(refreshToken)
	if err != nil {
		return err
	}

	r.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})
	r.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})
	return nil
}

func (a *TokenMiddleware) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "access_token cookie не найден", http.StatusUnauthorized)
			return
		}
		accessTokenString := accessTokenCookie.String()

		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "refresh_token cookie не найден", http.StatusUnauthorized)
			return
		}
		refreshTokenString := refreshTokenCookie.String()

		valid, err := a.tokenManager.ValidateTokenExpiration(accessTokenString)
		if valid {
			next.ServeHTTP(w, r)
		}

		valid, err = a.tokenManager.ValidateTokenExpiration(refreshTokenString)
		if valid {
			if err = a.updateTokenPair(r); err != nil {
				next.ServeHTTP(w, r)
			}
		}

		http.Error(w, "Не удалось обновить пару токенов", http.StatusUnauthorized)
		return
	})
}
