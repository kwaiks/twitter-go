package app

import (
	"fmt"
	"log"
)

func Fatal(err error){
	if err != nil {
		log.Fatal(err)
	}
}

func PrintError(err *interface{}){
	if err != nil {
		fmt.Println("Error Log:",err)
	}
}