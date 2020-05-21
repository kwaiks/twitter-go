package app

import (
	"database/sql"
	"net/http"
	"fmt"
	"../handler"
)

type ReqHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (app *App) requestHandler(handler ReqHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		handler(app.DB, w, r)
	}
}

// transform get Method
func (app *App) get(path string, f func(w http.ResponseWriter, r *http.Request)){
	endPoint := fmt.Sprintf("%s%s",app.BasePath,path)
	app.Router.HandleFunc(endPoint, f).Methods("GET")
}

func (app *App) post(path string, f func(w http.ResponseWriter, r *http.Request)){
	endPoint := fmt.Sprintf("%s%s",app.BasePath,path)
	app.Router.HandleFunc(endPoint, f).Methods("POST")
}

func (app *App) put(path string, f func(w http.ResponseWriter, r *http.Request)){
	endPoint := fmt.Sprintf("%s%s",app.BasePath,path)
	app.Router.HandleFunc(endPoint, f).Methods("PUT")
}

func (app *App) setRouters(){
	app.post("/login", app.requestHandler(handler.LoginHandler))
	app.post("/register", app.requestHandler(handler.RegisterHandler))
	app.post("/forgetPass", app.requestHandler(handler.ForgetPasswordHandler))
	app.post("/resetPass", app.requestHandler(handler.ResetPasswordHandler))
}