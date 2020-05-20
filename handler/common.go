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
}

func sendResponse(w http.ResponseWriter, status int, message string, payload interface{}){
	res := Response{}
	res.Data = &payload
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

func decodeUser(w http.ResponseWriter, body io.ReadCloser) models.User{
	var user models.User
	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return user
}