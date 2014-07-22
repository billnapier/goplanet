package commands

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

var mongodbSession *mgo.Session

func init() {
	RootCmd.PersistentFlags().String("mongodb_url", "localhost", "host where mongodb is")
	viper.BindPFlag("mongodb_url", RootCmd.PersistentFlags().Lookup("mongodb_url"))

	CreateUniqueIndexes()
}

func DBSession() *mgo.Session {
	if mongodbSession == nil {
		uri := viper.GetString("mongodb_url")

		var err error
		mongodbSession, err = mgo.Dial(uri)
		if mongodbSession == nil || err != nil {
			log.Fatalf("Can't connect to database, %v\n", err)
		}

		mongodbSession.SetSafe(&mgo.Safe{})
	}
	return mongodbSession
}

func DB() *mgo.Database {
	return DBSession().DB(viper.GetString("database_name"))
}

func Items() *mgo.Collection {
	return DB().C("items")
}

func Channels() *mgo.Collection {
	return DB().C("chanels")
}

func CreateUniqueIndexes() {
	idx := mgo.Index{
		Key:        []string{"key"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := Items().EnsureIndex(idx); err != nil {
		fmt.Println(err)
	}

	if err := Channels().EnsureIndex(idx); err != nil {
		fmt.Println(err)
	}
}
