package handler

import (
	"net/http"
	"encoding/json"
	"log"
	"io"

	"../models"
)

type Response struct{
	StatusCode int `json:"statusCode"`
	Message string `json:"message,omitempty"`
	Data *interface{} `json:"data,omitempty"`
	Token TokenStr `json:"token,omitempty"`
}

type TokenStr struct{
	Auth string `json:"reqToken"`
	Refresh string `json:"refToken,omitempty"`
}

func sendResponseNoPayload(w http.ResponseWriter, status int, message string){
	res := Response{}
	res.StatusCode = status
	res.Message = message

	response, err :=  json.Marshal(res)
	if err != nil {
		log.Fatalf("Error when sending Response : %s", err)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func sendResponse(w http.ResponseWriter, status int, payload interface{}, token string, refToken string){
	res := Response{}
	res.Data = &payload
	res.StatusCode = status
	res.Token.Auth = token
	res.Token.Refresh = refToken

	response, err :=  json.Marshal(res)
	if err != nil {
		log.Fatalf("Error when sending Response : %s", err)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func decodeUserFromMap(a map[string]interface{}) *models.User{
	var user models.User
	ab := a["user"].(map[string]interface{})
	if ab["username"] != nil {
		user.Username = ab["username"].(string)
	}
	if ab["name"] != nil {
		user.Name = ab["name"].(string)
	}
	if ab["password"] != nil {
		user.Password = ab["password"].(string) 
	}
	if ab["email"] != nil {
		user.Email = ab["email"].(string)
	}
	return &user
}

func decodeRequest(w http.ResponseWriter, body io.ReadCloser) map[string]interface{}{
	var m map[string]interface{}
	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return m
}