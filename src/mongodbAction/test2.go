package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		panic(err)
	}

	result := Person{}
	query:=make(bson.M)
	query["name"]="Ale"
//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	err = c.Find(query).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Phone:", result.Phone)
}

/*
如何使用

下面介绍如何使用mgo，主要介绍集合的操作。对数据库，用户等操作，请自行查看文档。

第一步当然是先导入mgo包

import ( "labix.org/v2/mgo" "labix.org/v2/mgo/bson" )

连接服务器

通过方法Dial()来和MongoDB服务器建立连接。Dial()定义如下：

func Dial(url string) (*Session, error)

具体使用：

session, err := mgo.Dial(url)

如果是本机，并且MongoDB是默认端口27017启动的话，下面几种方式都可以。

session, err := mgo.Dial("") session, err := mgo.Dial("localhost") session, err := mgo.Dial("127.0.0.1") session, err := mgo.Dial("localhost:27017") session, err := mgo.Dial("127.0.0.1:27017")

如果不在本机或端口不同，传入相应的地址即可。如：

mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb

切换数据库

通过Session.DB()来切换相应的数据库。

func (s Session) DB(name string) Database

如切换到test数据库。

db := session.DB("test")

切换集合

通过Database.C()方法切换集合（Collection），这样我们就可以通过对集合进行增删查改操作了。

func (db Database) C(name string) Collection

如切换到users集合。

c := db.C("users")

对集合进行操作

介绍插入、查询、修改、删除操作。

先提一下ObjectId，MongoDB每个集合都会一个名为_id的主键，这是一个24位的16进制字符串。对应到mgo中就是bson.ObjectId。

这里我们定义一个struct，用来和集合对应。

type User struct { Id_ bson.ObjectId bson:"_id" Name string bson:"name" Age int bson:"age" JoinedAt time.Time bson:"joined_at" Interests []string bson:"interests" }

注解

注意User的字段首字母大写，不然不可见。通过bson:”name”这种方式可以定义MongoDB中集合的字段名，如果不定义，mgo自动把struct的字段名首字母小写作为集合的字段名。如果不需要获得id，Id可以不定义，在插入的时候会自动生成。

插入

插入方法定义如下：

func (c *Collection) Insert(docs ...interface{}) error

下面代码插入两条集合数据。

err = c.Insert(&User{ Id_: bson.NewObjectId(), Name: "Jimmy Kuu", Age: 33, JoinedAt: time.Now(), Interests: []string{"Develop", "Movie"}, })

if err != nil { panic(err) }

err = c.Insert(&User{ Id_: bson.NewObjectId(), Name: "Tracy Yu", Age: 31, JoinedAt: time.Now(), Interests: []string{"Shoping", "TV"}, })

if err != nil { panic(err) }

这里通过bson.NewObjectId()来创建新的ObjectId，如果创建完需要用到的话，放在一个变量中即可，一般在Web开发中可以作为参数跳转到其他页面。

通过MongoDB客户端可以发现，两条即可已经插入。

{ "_id" : ObjectId( "5204af979955496907000001" ), "name" : "Jimmy Kuu", "age" : 33, "joned_at" : Date( 1376038807950 ), "interests" : [ "Develop", "Movie" ] }

{ "_id" : ObjectId( "5204af979955496907000002" ), "name" : "Tracy Yu", "age" : 31, "joned_at" : Date( 1376038807971 ), "interests" : [ "Shoping", "TV" ] }

查询

通过func (c Collection) Find(query interface{}) Query来进行查询，返回的Query struct可以有附加各种条件来进行过滤。

通过Query.All()可以获得所有结果，通过Query.One()可以获得一个结果，注意如果没有数据或者数量超过一个，One()会报错。

条件用bson.M{key: value}，注意key必须用MongoDB中的字段名，而不是struct的字段名。

无条件查询

查询所有

var users []User c.Find(nil).All(&users) fmt.Println(users)

上面代码可以把所有Users都查出来：

[{ObjectIdHex("5204af979955496907000001") Jimmy Kuu 33 2013-08-09 17:00:07.95 +0800 CST [Develop Movie]} {ObjectIdHex("5204af979955496907000002") Tracy Yu 31 2013-08-09 17:00:07.971 +0800 CST [Shoping TV]}]

根据ObjectId查询

id := "5204af979955496907000001" objectId := bson.ObjectIdHex(id)

user := new(User) c.Find(bson.M{"_id": objectId}).One(&user)

fmt.Println(user)

结果如下：

&{ObjectIdHex("5204af979955496907000001") Jimmy Kuu 33 2013-08-09 17:00:07.95 +0800 CST [Develop Movie]}

更简单的方式是直接用FindId()方法：

c.FindId(objectId).One(&user)

注解

注意这里没有处理err。当找不到的时候用One()方法会出错。

单条件查询

=($eq)

c.Find(bson.M{"name": "Jimmy Kuu"}).All(&users)

!=($ne)

c.Find(bson.M{"name": bson.M{"$ne": "Jimmy Kuu"}}).All(&users)

($gt)
c.Find(bson.M{"age": bson.M{"$gt": 32}}).All(&users)

<($lt)

c.Find(bson.M{"age": bson.M{"$lt": 32}}).All(&users)

=($gte)
c.Find(bson.M{"age": bson.M{"$gte": 33}}).All(&users)

<=($lte)

c.Find(bson.M{"age": bson.M{"$lte": 31}}).All(&users)

in($in)

c.Find(bson.M{"name": bson.M{"$in": []string{"Jimmy Kuu", "Tracy Yu"}}}).All(&users)

多条件查询

and($and)

c.Find(bson.M{"name": "Jimmy Kuu", "age": 33}).All(&users)

or($or)

c.Find(bson.M{"$or": []bson.M{bson.M{"name": "Jimmy Kuu"}, bson.M{"age": 31}}}).All(&users)

修改

通过func (*Collection) Update来进行修改操作。

func (c *Collection) Update(selector interface{}, change interface{}) error

注意修改单个或多个字段需要通过$set操作符号，否则集合会被替换。

修改字段的值($set)

c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$set": bson.M{ "name": "Jimmy Gu", "age": 34, }})

inc($inc)

字段增加值

c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$inc": bson.M{ "age": -1, }})

push($push)

从数组中增加一个元素

c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$push": bson.M{ "interests": "Golang", }})

pull($pull)

从数组中删除一个元素

c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$pull": bson.M{ "interests": "Golang", }})

删除

c.Remove(bson.M{"name": "Jimmy Kuu"})

注解

这里也支持多条件，参考多条件查询。
*/
