package database

import (
	"gopkg.in/mgo.v2"
	"os"
)

type Mongo struct {
	Posts *mgo.Collection
}

func CreateConn() *Mongo {
	session, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}

	db := session.DB("remotejopbapp")

	return &Mongo{
		Posts: db.C("posts"),
	}
}
