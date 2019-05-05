// {
//    $lookup:
//      {
//        from: <collection to join>,
//        localField: <field from the input documents>,
//        foreignField: <field from the documents of the "from" collection>,
//        as: <output array field>
//      }
// }
// From: $lookup (aggregation)

// 其中的lookup功能可以实现类似于mysql中的join操作，方便于关联查询。

// 2. mgo中的实现外键关联

// See also
// https://docs.mongodb.com/manual/aggregation/
// https://docs.mongodb.com/manual/core/aggregation-pipeline/
// https://github.com/go-mgo/mgo/issues/248

package main

import (
    "time"

    "fmt"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type SuiteGroup struct {
    ID     bson.ObjectId `bson:"_id,omitempty" json:"_id"`
    Name    string        `bson:"name,omitempty" json:"name,omitempty"`
    SuiteID int `bson:"_suiteId" json:"_suiteId"`
    Kind    string        `bson:"kind,omitempty" json:"kind,omitempty"`

    // pipeline get, Suites []Suite
    Suite []Suite `bson:"suite" json:"suites"`
}

type Suite struct {
	ID *bson.ObjectId `bson:"_id" json:"_id,omitempty"`
	//Created *time.Time     `bson:"created,omitempty" json:"created,omitempty"`
	IDD string `bson:"idd" json:"idd"`
	Age int    `bson:"age" json:"age"`
}


func main() {
    session, err := mgo.Dial("127.0.0.1")
    if err != nil {
        panic(err)
    }

    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c1 := session.DB("test").C("suitegroups")
    c2:= session.DB("test").C("suite")

     //8 插入数据
	newSuiteGroup := &SuiteGroup{
        Name: "zhu1",
        SuiteID: 1,
        Kind: "zhu1-kind",
    }

	err=c1.Insert(newSuiteGroup)
    if err!=nil{
        fmt.Println(err)
    }

         //8 插入数据


	err=c1.Insert(newSuiteGroup)
    if err!=nil{
        fmt.Println(err)
    }


    pipeline := []bson.M{
        bson.M{"$match": bson.M{"name": "Admin"}},
        bson.M{"$lookup": bson.M{"from": "suites", "localField": "_suiteId", "foreignField": "idd", "as": "suite"}},
    }

    var resp SuiteGroup
    c1.Pipe(pipeline).One(&resp)

    fmt.Printf("SuiteID [%s] SuiteObj %#v\n", resp.SuiteID.String(), resp.Suite[0])
}