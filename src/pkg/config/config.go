package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type DbParams struct {
	Uri    string
	DbName string
}

type JwtParams struct {
	SigningKey string
	Expiration int
}

func Init() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении текущей директории: %v", err)
	}

	envFilePath := filepath.Join(filepath.Dir(filepath.Dir(rootDir)), ".env")
	envErr := godotenv.Load(envFilePath)
	if envErr != nil {
		log.Fatalf("Ошибка при загрузке .env файла: %v", envErr)
	}
}

func GetDbParams() DbParams {
	dbParamsObject := DbParams{
		Uri:    os.Getenv("MONGO_URI"),
		DbName: os.Getenv("MONGO_DB_NAME"),
	}

	if dbParamsObject.DbName == "" || dbParamsObject.Uri == "" {
		log.Fatal("Ошибка: не все параметры базы данных были получены. Проверьте .env файл.")
	}

	return dbParamsObject
}

func GetJwtParams() JwtParams {
	signingKeyStr := os.Getenv("JWT_SIGNING_KEY")
	if signingKeyStr == "" {
		log.Fatal("Ошибка: параметр JWT_SIGNING_KEY не был получен. Проверьте .env файл.")
	}
	expirationStr := os.Getenv("JWT_EXPIRATION_HOUR")
	if expirationStr == "" {
		log.Fatal("Ошибка: параметр JWT_SIGNING_KEY не был получен. Проверьте .env файл.")
	}

	expirationInt, err := strconv.Atoi(expirationStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования JWT_EXPIRATION_HOUR: %v", err)
	}

	return JwtParams{SigningKey: signingKeyStr, Expiration: expirationInt}
}
