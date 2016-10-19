package handlers

import (
	"regexp"
	"net/http"
	"github.com/pkg/errors"
	"github.com/dgrijalva/jwt-go"

	. "github.com/andrepinto/goway/handlers"
	"github.com/andrepinto/goway/router"
)

var regexpBearer = regexp.MustCompile("^Bearer\\s([A-Za-z0-9\\-\\._~\\+\\/]+=*)$")

type OpCallback func ( *router.Route, *http.Request, *jwt.MapClaims ) *HandlerError
type ErrCallback func (  *router.Route, *http.Request, *HandlerError) *HandlerError

//noinspection GoUnusedParameter
func opCallback( route *router.Route, r *http.Request, claim *jwt.MapClaims ) *HandlerError { return nil }
//noinspection GoUnusedParameter
func errCallback(  route *router.Route, r *http.Request, err *HandlerError) *HandlerError {
	return err
}



type JWTHandler struct {
	Secret string
	Algorithm 	*jwt.SigningMethodHMAC
	OnSuccess	OpCallback
	OnError		ErrCallback
}
//noinspection GoUnusedExportedFunction
func NewJWTHandler( secret string, algorithm *jwt.SigningMethodHMAC ) (*JWTHandler) {
	return &JWTHandler{
		Secret: secret,
		Algorithm: algorithm,
		OnSuccess: opCallback,
		OnError: errCallback,
	}
}


func ( handler *JWTHandler ) validateSignature(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("Invalid token.")
	}

	return []byte(handler.Secret), nil
}

func (handler *JWTHandler) decode(tk string) (*jwt.MapClaims, error) {

	var claims jwt.MapClaims
	var err error
	var ok bool

	token, err := jwt.Parse(tk, handler.validateSignature)
	if ( err != nil ) {
		return nil, err
	}

	claims, ok = token.Claims.(jwt.MapClaims);
	if (!ok) {
		return nil, errors.New("Token is invalid")
	}
	return  &claims, nil
}

func (handler *JWTHandler) Handle(route *router.Route, req *http.Request) (*HandlerError){

	matches := regexpBearer.FindStringSubmatch(req.Header.Get("Authorization"))


	if ( len(matches) != 2 ) {
		return handler.OnError( route, req, NewHttpError(  http.StatusUnauthorized, "Unauthorized"))
	}

	claim, err := handler.decode(matches[1])
	if (err != nil) {
		return handler.OnError( route, req, NewHttpError( http.StatusUnauthorized, err.Error()))
	}

	if err := claim.Valid(); err != nil {
		return handler.OnError( route, req, NewHttpError( http.StatusUnauthorized, err.Error()))
	}


	return handler.OnSuccess( route, req, claim )
}


