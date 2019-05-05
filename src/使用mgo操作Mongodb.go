package main

/*
usage:mgo的使用
author:Mr-Luo
date:2015年7月16日 09:54:14
version:0.1
*/
import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var url string

//注意struct 的字段要大写
type Blog struct {
	Sex    string
	Author string
	Date   string
}

func Get() *Blog {
	return nil
}

func init() {
	url = "localhost:27017"
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
	log.Println(session.BuildInfo())

	//2. 获取database
	db := session.DB("")

	c := db.C("blog")
	//3 获取query
	query := c.Find(bson.M{"author": "milo"})

	//4 获取数量
	count, _ := query.Count()
	log.Println(count)

	//5 获取一个结果
	blog := Blog{}
	query.One(&blog)
	log.Println(blog) //2015/07/16 10:45:21 {0 milo 2015-07-08 16:45:14}

	//6 获取所有结果
	//blogs := make([]Blog, 1)
	blogs := []Blog{}
	query.All(&blogs)
	log.Println(blogs)

	//7 更新数据
	selector := bson.M{"_id": bson.ObjectIdHex("559dcd43178eed0768dd21ad")}
	update := bson.M{"$set": bson.M{"address": "kp"}}
	e := c.Update(selector, update)
	log.Println(e)

	//8 插入数据
	newblog := &Blog{Author: "m", Sex: "female"}
	e1 := c.Insert(newblog)
	log.Println(e1)

	//9 删除数据
	e2 := c.RemoveId(bson.ObjectIdHex("559ce34d178eed10e0b6dab4")) //这里的id也要hex
	log.Println(e2)
}
