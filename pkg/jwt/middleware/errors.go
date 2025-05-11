package jwt

import "fmt"

var (
	ErrAccessTokenNotFound  = fmt.Errorf("токен доступа не найден")
	ErrRefreshTokenNotFound = fmt.Errorf("токен обновления не найден")
)
