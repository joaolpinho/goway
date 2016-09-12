package handlers
import (
	"github.com/andrepinto/goway/router"
	"fmt"
	"net/http"
)


func Jwt(route *router.Route, r *http.Request)(bool){
	fmt.Println("auth")
	r.Header.Set("USERID","123456")
	return true
}


