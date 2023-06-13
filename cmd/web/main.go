package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/omeg21/project-repo/internal/config"
	"github.com/omeg21/project-repo/internal/controller"
	"github.com/omeg21/project-repo/internal/driver"
	"github.com/omeg21/project-repo/internal/helpers"
	"github.com/omeg21/project-repo/internal/models"
	"github.com/omeg21/project-repo/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db,err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB,error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	
	//
	infoLog = log.New(os.Stdout,"info\t",log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	//
	errorLog =log.New(os.Stdout,"error\t",log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("connecting to database")
	db,err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=yep password=1234")
	if err != nil{
		log.Fatal("can not connect to database Deadge")
	}

	log.Println("connected to database")

	
	
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil,err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := controller.NewRepo(&app,db)
	controller.NewController(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db,nil
}