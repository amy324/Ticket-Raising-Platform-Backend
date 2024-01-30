package main

import (
	"database/sql"
	"log"
	"backend-project/data"
	"sync"

	"github.com/alexedwards/scs/v2"
)

type Configuration struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
	Mailer   Mail
}
