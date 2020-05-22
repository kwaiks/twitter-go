package config

import (
	"strconv"
	"os"
)

type DBConfig struct{
	DB *DB
}

type DB struct{
	Type		string
	Host		string
	Username	string
	Password	string
	Port		int64
	DBName		string
}

func GetDBConfig() *DBConfig{
	port,_ := strconv.ParseInt(os.Getenv("db_port"),0,64)
	return &DBConfig{
		DB: &DB{
			Type: os.Getenv("db_type"),
			Host: os.Getenv("db_host"),
			Username: os.Getenv("db_username"),
			Password: os.Getenv("db_password"),
			Port: port,
			DBName: os.Getenv("db_name"),
		},
	}
}