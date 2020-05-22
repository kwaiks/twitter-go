package handler

import (
	_"github.com/dgrijalva/jwt-go"
	"context"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	_ "github.com/lib/pq"
	"database/sql"
	"log"
	"fmt"
	gCon "github.com/gorilla/context"
	"../models"
	myJwt "../jwt"
	"time"
)

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request){
	res := decodeRequest(w, r.Body)
	user := decodeUserFromMap(res)
	tes := gCon.Get(r,"username")

	fmt.Println(tes)


	if user.Username == "" || user.Password == "" {
		sendResponseNoPayload(w, 400, "Found empty field")
		return
	}

	var err error 

	getUserQuery := "SELECT password,username FROM users WHERE username = $1"
	var password string
	var userRes models.User
	
	err = db.QueryRow(getUserQuery, user.Username).Scan(&password,&userRes.Username)
	if err == sql.ErrNoRows {
		sendResponseNoPayload(w, 404, "User not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password),[]byte(user.Password))
	if err != nil {
		sendResponseNoPayload(w, 401, "Password doesn't match")
		return
	}
	//jwt := jwt.New()
	token := myJwt.GenerateJWT(userRes.Username)
	refToken := myJwt.GenerateRefreshToken()
	refCookie := http.Cookie{Name:"rtk",Value:refToken,Expires:time.Now().AddDate(1,0,0)}
	tknCookie := http.Cookie{Name:"qtk",Value:token}
	http.SetCookie(w, &tknCookie)
	http.SetCookie(w, &refCookie)
	sendResponse(w, 200, nil, token, refToken)
}

func RegisterHandler(db *sql.DB, w http.ResponseWriter, r *http.Request){
	res := decodeRequest(w, r.Body)
	user := decodeUserFromMap(res)

	checkUserNameQuery := "SELECT username, email from users WHERE username=$1 OR email=$2"

	// check if username already registered
	var username string
	var email string
	errDb := db.QueryRow(checkUserNameQuery,user.Username,user.Email).Scan(&username,&email)
	if errDb == nil { 
		sendResponseNoPayload(w, 202, "Username or Email has been used")
		return
	}

	if user.Password == "" {
		sendResponseNoPayload(w, 200, "Username and email could be used")
		return
	}

	// register user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password),10) //default cost
	queryString := "INSERT INTO users (username,password,email,name) VALUES ($1,$2,$3,$4)"
	_, err := db.Exec(queryString, user.Username, hashedPassword, user.Email, user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error Ocurred"))
		return
	}

	sendResponseNoPayload(w, 201, "Successfully registered")
}

func ForgetPasswordHandler(db *sql.DB, w http.ResponseWriter, r *http.Request){
	res := decodeRequest(w, r.Body)
	user := decodeUserFromMap(res)

	var err error

	queryString := "SELECT username,id FROM users WHERE username=$1 AND email=$2"
	var userMdl models.User
	err = db.QueryRow(queryString, user.Username, user.Email).Scan(&userMdl.Username,&userMdl.ID)
	if err == sql.ErrNoRows {
		sendResponseNoPayload(w, 404, "User not found")
		return
	}

	generatedKey := randString(30)
	
	ctx := context.Background()
	tx,err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, execErr := tx.ExecContext(ctx,"INSERT INTO reset_pass_keys (userid,key) VALUES ($1,$2)",userMdl.ID, generatedKey)
	if execErr != nil {
		tx.Rollback()
		fmt.Println("insert") // just need to know where its error
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}
	_, execErr = tx.ExecContext(ctx,"UPDATE users SET reset_password=1")
	if execErr != nil {
		tx.Rollback()
		fmt.Println("	")
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		http.Error(w, commitErr.Error(), http.StatusInternalServerError)
		return
	}
	sendResponse(w, http.StatusOK, generatedKey, "", "")
}

func ResetPasswordHandler(db *sql.DB, w http.ResponseWriter, r *http.Request){
	res := decodeRequest(w, r.Body)
	key := res["key"]
	user := decodeUserFromMap(res)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password),10)

	_, err := db.Exec("UPDATE users SET password=$1, reset_password=0 FROM reset_pass_keys WHERE users.id = reset_pass_keys.userid AND reset_pass_keys.key = $2 AND users.reset_password = 1",hashedPassword,key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResponseNoPayload(w, http.StatusOK, "Password has been reseted")
}

func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}