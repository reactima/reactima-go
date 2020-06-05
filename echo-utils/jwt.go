package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func GetJWTUserID (c echo.Context) (int64, error) {
	ctxUser := c.Get("user").(*jwt.Token)

	// TODO rewrite, add errors
	id := ctxUser.Claims.(jwt.MapClaims)["id"].(float64)
	iid:=int64(id)

	return iid, nil
}