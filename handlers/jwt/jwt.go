package handlers

import (
	"fmt"
	"regexp"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/andrepinto/goway/router"
	"github.com/pkg/errors"
)

var regexpBearer = regexp.MustCompile("^Bearer\\s([A-Za-z0-9\\-\\._~\\+\\/]+=*)$")


func opCallback( route *router.Route, r *http.Request, claim *jwt.MapClaims ) bool { return true }
func errCallback(  route *router.Route, r *http.Request, err error ) bool { return false }


func NewJWTHandler( secret string, algorithm *jwt.SigningMethodHMAC ) (*JWTHandler) {
	return &JWTHandler{
		Secret: secret,
		Algorithm: algorithm,
		OnSuccess: opCallback,
		OnError: errCallback,
	}
}
func MakeJWTHandler( secret string, algorithm *jwt.SigningMethodHMAC ) (JWTHandler) {
	return *NewJWTHandler(secret, algorithm)
}


func ( handler *JWTHandler ) validateSignature(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
		return nil, errors.New("Invalid claim.")
	}
	return  &claims, nil
}

func (handler *JWTHandler) Handle(route *router.Route, r *http.Request)(bool){

	matches := regexpBearer.FindStringSubmatch(r.Header.Get("Authorization"))


	if ( len(matches) != 2 ) {
		return false
	}

	claim, err := handler.decode(matches[1])
	if (err != nil) {
		fmt.Errorf("%s", err.Error())
		return handler.OnError( route, r, err  )
	}

	if err := claim.Valid(); err != nil {
		fmt.Errorf("%s", err.Error())
		return handler.OnError( route, r, err  )
	}


	return handler.OnSuccess( route, r, claim )
}


