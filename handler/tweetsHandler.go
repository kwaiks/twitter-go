package handler

import (
	"net/http"
	_"log"
	_"fmt"

	_"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"database/sql"
)

func postTweet(db *sql.DB, w http.ResponseWriter, r *http.Request){
	
}

func GetTweetById(db *sql.DB, w http.ResponseWriter, r *http.Request){
	//vars := mux.Vars(r)
	//userId := vars["uid"]
	
}