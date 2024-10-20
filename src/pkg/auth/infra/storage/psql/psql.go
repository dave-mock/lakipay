package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func New(log *log.Logger) (*sql.DB, error) {
	envFilePath := "C:\\Users\\Rinatoh computer\\Desktop\\lakipay\\.env" // Use double backslashes for Windows paths
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Println("Error loading .env file:", err)
		return nil, err
	}


	const (
		ErrFailedToConnect = "FAILED_TO_CONNECT"
		ErrFailedToPing    = "FAILED_TO_PING"
		ErrFailedToCreate  = "FAILED_TO_CREATE"
	)

	log.SetPrefix(log.Prefix() + " [AUTH] [INFRA] [STORAGE] [PSQL] ")
	log.Println("Initiating pg db")

	var (
		host, _     = os.LookupEnv("DB_HOST")
		user, _     = os.LookupEnv("DB_USER")
		password, _ = os.LookupEnv("DB_PASS")
		dbName, _   = os.LookupEnv("DB_NAME")
		sslMode, _  = os.LookupEnv("SSL_MODE")
		port, _     = os.LookupEnv("DB_PORT")
	)

	// pg connection
	fmt.Println("|||||||||", host, user, password)

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode))
	if err != nil {
		log.Println("Failed to connect to db server: " + err.Error())
		return nil, err
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(0)

	if err = db.Ping(); err != nil {
		log.Println("Ping is not responding")
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "invalid_catalog_name":
				log.Println("The specified database does not exist, trying to create one")
				_, err := db.Exec("create database " + dbName)
				if err != nil {
					log.Println(err.Error())
				}
				if err = db.Ping(); err != nil {
					log.Println("Pinging the new database failed: " + err.Error())
				}
				return db, nil
			}
			return nil, err
		}
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}
