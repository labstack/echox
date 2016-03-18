package main

import (
	"fmt"
	"net/http"
	
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User represent a document in the collection users
type User struct {
	ID    bson.ObjectId `json:"_id" bson:"_id"`
	Fname string        `json:"fname" bson:"fname"`
}

// getUser returns one document from the collection users
func getUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := getMgo(c)
		user := User{}
		id := bson.ObjectIdHex(c.Param("id"))
		
		err := db.C("users").FindId(id).One(&user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, user)
	}
}

// getMgo retreives mgo key
func getMgo(c echo.Context) *mgo.Database {
	return c.Get("mgok").(*mgo.Database)
}

// mgo middleware
func mgoMid(db *mgo.Database) echo.MiddlewareFunc {
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			s := db.Session.Clone()
			defer s.Close()
			c.Set("mgok", s.DB(db.Name))
			return next.Handle(c)
		})
	}
}

// newMgo dial mongo database url and name
func newMgo(url string, name string) *mgo.Database {
	s, err := mgo.Dial(url)
	if err != nil {
		panic(fmt.Sprintf("%v - %s", err, url))
	}
	return &mgo.Database{s, name}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	logg := e.Logger()

	// gopkg.in/mgo.v2
	db := newMgo("localhost:27017", "dbname")
	defer db.Session.Close()
	e.Use(mgoMid(db))
	
	// Route for user
	e.Get("/user/:id", getUser())

	// mgo database info
	mgoInf, _ := db.Session.BuildInfo()
	logg.Printf("MongoDB %s %d bits", mgoInf.Version, mgoInf.Bits)

	e.Run(standard.New(":1323"))
}
