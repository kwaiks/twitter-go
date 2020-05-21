package app

import (
	"github.com/rs/cors"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"database/sql"

	"../config"
)

type App struct {
	Router		*mux.Router
	DB			*sql.DB
	BasePath	string
}

func (app *App) Initialize(){
	dbConfig := config.GetDBConfig() //Retrieve DB Config
	app.BasePath = config.GetAPIPath()
	fmt.Println(app.BasePath)
	app.DB = app.setDBConnection(dbConfig) // Set the connection
	app.Router = mux.NewRouter()
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
	err := http.ListenAndServe(port, cors.Default().Handler(app.Router))
	Fatal(err)
}