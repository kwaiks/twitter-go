package app

import (
	"database/sql"
	"net/http"
	"../handler"
)

const (
	post="POST"
	get="GET"
)

type ReqHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (app *App) requestHandler(handler ReqHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		handler(app.DB, w, r)
	}
}

// transform get Method
func (app *App) privateRoute(path string, method string, f func(w http.ResponseWriter, r *http.Request)){
	app.PrivateRoute.HandleFunc(path, f).Methods(method)
}

func (app *App) publicRoute(path string, method string, f func(w http.ResponseWriter, r *http.Request)){
	app.Router.HandleFunc(path, f).Methods("POST")
}

func (app *App) setRouters(){
	app.publicRoute("/login", post, app.requestHandler(handler.LoginHandler))
	app.publicRoute("/register", post, app.requestHandler(handler.RegisterHandler))
	app.publicRoute("/forgetPass", post, app.requestHandler(handler.ForgetPasswordHandler))
	app.publicRoute("/resetPass", post, app.requestHandler(handler.ResetPasswordHandler))

	app.privateRoute("/getProfile", get, app.requestHandler(handler.GetOwnProfile))
	app.privateRoute("/getUserProfile/{user}", get, app.requestHandler(handler.GetUserProfile))
}
