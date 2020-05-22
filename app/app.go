package app

import (
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"database/sql"
	

	"../config"
)

type App struct {
	Router			*mux.Router
	PrivateRoute	*mux.Router
	DB				*sql.DB
}

func (app *App) Initialize(){
	err := godotenv.Load() //load .env file
	Fatal(err)

	
	dbConfig := config.GetDBConfig() //Retrieve DB Config
	app.DB = app.setDBConnection(dbConfig) // Set the connection
	app.Router = mux.NewRouter().PathPrefix("/api").Subrouter()
	app.PrivateRoute = app.Router.PathPrefix("/service").Subrouter()
	app.PrivateRoute.Use(authMiddleware)
	app.setRouters()
}

func (app *App) setDBConnection(config *config.DBConfig) *sql.DB{
	connString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.Type,
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.DBName,
	)

	db, err := sql.Open(config.DB.Type, connString)
	Fatal(err)

	return db
}

func (app *App) Run(port string){
	c := cors.New(cors.Options{
		AllowedOrigins:[]string{"http://localhost:3000","localhost"},
		Debug:true,
		AllowedHeaders: []string{"Content-Type","Bearer","Bearer ","content-type","Origin","Accept","Authorization","Refresh-Token"},
	})
	err := http.ListenAndServe(port, c.Handler(app.Router)) //Added CORS to integrate with React
	Fatal(err)
}