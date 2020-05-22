package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"database/sql"
	"../models"
)

func GetOwnProfile(db *sql.DB, w http.ResponseWriter, r *http.Request){
	token := r.URL.Query().Get("token")
	username := r.URL.Query().Get("username")

	var user models.User
	err := db.QueryRow("SELECT username,name,email,photo from users WHERE username = $1",username).Scan(&user.Username,&user.Name,&user.Email,&user.Photo)
	if err != nil {
		sendResponseNoPayload(w, 404, "Data not found")
		return
	}
	sendResponse(w, 200 ,user,token,"")
}

func GetUserProfile(db *sql.DB, w http.ResponseWriter, r *http.Request){
	token := r.URL.Query().Get("token")
	currentUser := r.URL.Query().Get("username")
	vars := mux.Vars(r)
	username := vars["user"]

	var data struct{
		User models.User `json:"user"`
		Followed bool `json:"followed"`
		Following bool `json:"following"`
	}

	db.QueryRow("SELECT id,username,name,email,photo from users WHERE username = $1",username).Scan(&data.User.ID,&data.User.Username,&data.User.Name,&data.User.Photo)
	data.Followed = checkFollow(db, username, currentUser)
	data.Following = checkFollow(db, currentUser, username)

	sendResponse(w, 200, data, token, "")
}

func checkFollow(db *sql.DB, current string, target string) bool {
	var id int
	q := db.QueryRow("SELECT a.id FROM followers f INNER JOIN users a ON a.id = f.userid INNER JOIN users b ON b.id = f.followed_userid WHERE a.username=$1 AND b.username=$2",current,target).Scan(id)
	if q == sql.ErrNoRows{
		return false
	}
	return true
}