// mongodb
package main

import (
    "fmt"
    "log"
    "time"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Task struct {
    Description string
    Due         time.Time
}

type Category struct {
    Id          bson.ObjectId `bson:"_id,omitempty"`
    Name        string
    Description string
    //Tasks       []Task
}

func main() {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    //获取一个集合
    c := session.DB("taskdb").C("categories")
    c.RemoveAll(nil)

    //index
    index := mgo.Index{
        Key:        []string{"name"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }

    //create Index
    err = c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }

    //插入三个值
    err = c.Insert(&Category{bson.NewObjectId(), "R & D", "R & D Tasks"},
        &Category{bson.NewObjectId(), "Project", "Project Tasks"},
        &Category{bson.NewObjectId(), "Open Source", "Tasks for open-source projects"})

    if err != nil {
        panic(err)
    }

    result := Category{}
    err = c.Find(bson.M{"name": "Open-Source"}).One(&result)
    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println("Description:", result.Description)

    }
