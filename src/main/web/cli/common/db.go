package common

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

func Connect(dbName, cName, env string) (*mgo.Session, *mgo.Collection) {
	var session *mgo.Session
	var err error
	if env == "dev" {
		fmt.Println("environment:", "dev")
		session, err = mgo.Dial("")
	} else if env == "test" {
		fmt.Println("environment:", "test")
		session, err = mgo.Dial("")
	} else {
		fmt.Println("environment:", "local")
		session, err = mgo.Dial("192.168.136.130:27017")
		//session, err = mgo.Dial("192.168.244.128:27017")
	}
	ErrorHandler(err)
	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(1 * time.Hour)
	collection := session.DB(dbName).C(cName)
	return session, collection
}
