package main

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"log"
)

var url string

func init() {
	url = "192.168.20.22:27017"
}

func Conn(url string) (*mgo.Session, error) {
	session, err := mgo.Dial(url)
	return session, err
}

func main() {
	//1. 创建session
	session, err := Conn(url)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(session.BuildInfo())

	//2. 获取database
	db := session.DB("mydb")

	c := db.C("c1")

	indexes, err := c.Indexes()
	if err != nil {
		panic(err)
	}
	log.Println(indexes)

	// query := c.Find(bson.M{"age": 20})

	// count, _ := query.Count()
	// log.Println(count)

}
