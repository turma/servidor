package main

import (
	"database/sql"
	"log"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	gorp.SqlExecutor
}

// The only one instance of db
var db DB

func OrmMiddleware(c martini.Context) {
	// Inject our db in the hanlders
	c.MapTo(db, (*DB)(nil))
}

func init() {
	var openSQL string
	if martini.Env == "production" {
		log.Println("Databse starting in production")
		openSQL = EnvProd.DBUser + ":" + EnvProd.DBPass + "@/" + EnvProd.DBName + "?charset=utf8&parseTime=true"
	} else {
		log.Println("Databse starting in development")
		openSQL = EnvDev.DBUser + ":" + EnvDev.DBPass + "@/" + EnvDev.DBName + "?charset=utf8&parseTime=true"
	}

	// Connect to db using standard Go database/sql API
	dbopen, err := sql.Open("mysql", openSQL)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	//defer db.Close() // I DUNNO IF IT WORKS HERE, LETS TEST

	dialect := gorp.MySQLDialect{"InnoDB", "UTF8"}

	// Construct a gorp DbMap using MySQL dialect
	dbmap := gorp.DbMap{Db: dbopen, Dialect: dialect}

	log.Println("Database connected!")

	// Adding schemes to my ORM
	dbmap.AddTableWithName(User{}, "user").SetKeys(false, "id")
	dbmap.AddTableWithName(Link{}, "link").SetKeys(false, "link")
	dbmap.AddTableWithName(Shared{}, "shared").SetKeys(false, "id")
	dbmap.AddTableWithName(Photo{}, "photo").SetKeys(false, "id")

	// Adding to local vairable
	db = &dbmap

	// Disabled for a while
	//dbmap.TraceOn("[SQL]", log.New(os.Stdout, "[DB]", log.Lmicroseconds))

}
