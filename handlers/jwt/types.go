package handlers

import (
	"net/http"
	"github.com/andrepinto/goway/router"
	"github.com/dgrijalva/jwt-go"
)

type OpCallback func ( *router.Route, *http.Request, *jwt.MapClaims ) bool
type ErrCallback func (  *router.Route, *http.Request, error ) bool

type JWTHandler struct {
	Secret string
	Algorithm 	*jwt.SigningMethodHMAC
	OnSuccess	OpCallback
	OnError		ErrCallback
}