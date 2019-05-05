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
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SuiteGroup struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Name    string        `bson:"name,omitempty" json:"name,omitempty"`
	SuiteID int           `bson:"_suiteId" json:"_suiteId"`
	Kind    string        `bson:"kind,omitempty" json:"kind,omitempty"`

	// pipeline get, Suites []Suite
	Suite []Suite `bson:"suites" json:"suites"`
}

type Suite struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	//Created *time.Time     `bson:"created,omitempty" json:"created,omitempty"`
	IDD int `bson:"_idd" json:"_idd"`
	Age int `bson:"age" json:"age"`
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
	c1.RemoveAll(nil)
	c2 := session.DB("test").C("suite")
	c2.RemoveAll(nil)

	//8 插入数据
	newSuiteGroup := &SuiteGroup{
		Name:    "Admin",
		SuiteID: 1,
		Kind:    "zhu1-kind",
	}

	err = c1.Insert(newSuiteGroup)
	if err != nil {
		fmt.Println(err)
	}

	a := &Suite{
		IDD: 1,
		Age: 2,
	}
	err = c2.Insert(a)
	if err != nil {
		fmt.Println(err)
	}

	err = c2.Insert(&Suite{
		IDD: 1,
		Age: 3,
	})
	if err != nil {
		fmt.Println(err)
	}

	err = c2.Insert(&Suite{
		IDD: 2,
		Age: 2,
	})
	if err != nil {
		fmt.Println(err)
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"name": "Admin"}},
		bson.M{"$lookup": bson.M{"from": "suite", "localField": "_suiteId", "foreignField": "_idd", "as": "suites"}},
	}

	var resp SuiteGroup
	//var resp interface{}
	c1.Pipe(pipeline).One(&resp)

	fmt.Println(resp)

	//fmt.Printf("SuiteID [%v] SuiteObj %#v\n", resp.SuiteID, resp.Suite[0])
}


//在测试中