package main

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

func main() {

    err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }

	a := App{}
	a.Initialize(os.Getenv("DB_URI"),
                 os.Getenv("NOTIFICATION_URI"),
                 os.Getenv("MAIL_URI"),
                 os.Getenv("DOCUMENT_LIMIT"))

	a.Run(os.Getenv("PORT"))

}

