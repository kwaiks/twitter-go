package config

type DBConfig struct{
	DB *DB
}

type DB struct{
	Type		string
	Host		string
	Username	string
	Password	string
	Port		int16
	DBName		string
}

func GetDBConfig() *DBConfig{
	return &DBConfig{
		DB: &DB{
			Type: "postgres",
			Host: "localhost",
			Username: "postgres",
			Password: "alex02",
			Port: 5432,
			DBName: "twitterGo",
		},
	}
}

func GetAPIPath() string{
	return "/api/v1"
}