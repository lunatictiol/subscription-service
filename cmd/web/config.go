package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/lunatictiol/subscription-service/data"
)

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
}
