package config

import "os"

type MongoDb struct {
	Password string
	Name     string
}

func MongoDbConfig() MongoDb {
	dbConfig := MongoDb{
		Password: os.Getenv("MONGO_DB_PASSWORD"),
		Name:     os.Getenv("MONGO_DB_NAME"),
	}

	return dbConfig
}

type AppEnvironment struct {
	AppEnv string
}

func AppConfig() AppEnvironment {
	appConfig := AppEnvironment{
		AppEnv: os.Getenv("APP_ENV"),
	}

	return appConfig
}
