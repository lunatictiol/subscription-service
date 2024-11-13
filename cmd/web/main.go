package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lunatictiol/subscription-service/data"
)

const PORT = "80"

func main() {
	db := initDB()
	session := initSession()
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	wg := sync.WaitGroup{}

	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLogger,
		ErrorLog: errorLogger,
		Wait:     &wg,
		Models:   data.New(db),
	}

	//app.listenForShutDown()
	app.startServer()

}

func (app *Config) startServer() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting server")
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

//database

func initDB() *sql.DB {
	conn := connectDB()
	if conn == nil {
		log.Panic("cant connnect to db")
	}
	return conn
}

func connectDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Panicln("postgres not ready yet")
		} else {
			log.Println("connection successful")
			return conn
		}

		if counts > 10 {
			return nil
		}

		log.Println("backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

//session

func initSession() *scs.SessionManager {
	gob.Register(data.User{})
	session := scs.New()
	session.Store = redisstore.New(initReddis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return session

}

func initReddis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}

// func (app *Config) listenForShutDown() {
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	app.ShutDown()
// 	os.Exit(0)
// }

// func (app *Config) ShutDown() {
// 	//perform clean up

// 	app.InfoLog.Println("Performing clean up")

// 	app.Wait.Wait()

// 	app.InfoLog.Println("closing all channels and terminating the application")
// }
