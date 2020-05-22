package app

import (
	"strings"
	"github.com/gorilla/context"
	"net/http"
	"../jwt"
)

func authMiddleware(handler http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		refToken := r.Header.Get("Refresh-Token")
		reqToken := r.Header.Get("Authorization")
		token := strings.Split(reqToken, "Bearer ")

		//refresh user token
		scs, user, tkn := jwt.ValidateToken(refToken, token[1])
		if scs == false{
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized"))
			return
		}
		if r.Method == "POST"{
			context.Set(r, "token", tkn)
			context.Set(r, "username", user)
		}
		if r.Method == "GET"{
			if r.URL.RawQuery == ""{
				r.URL.RawQuery = r.URL.RawQuery + "username=" + user + "&token=" + tkn
			}else{
				r.URL.RawQuery = r.URL.RawQuery + "&username=" + user + "&token=" + tkn
			}
			
		}
		handler.ServeHTTP(w,r)
	})
}