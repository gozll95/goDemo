package main

import "log"
import "github.com/gin-gonic/gin"

type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	route := gin.Default()
	route.Any("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	if c.Bind(&person) == nil {
		log.Println("====== Only Bind Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}

/*
Bind只会bind query的,
对于GET:name=eason&address=xyz
对于POST:name=eason&address=xyz
*/
/*
# only bind query
$ curl -X GET "localhost:8085/testing?name=eason&address=xyz"
====== Only Bind Query String ======
eason
xyz


# only bind query string, ignore form data
$ curl -X POST "localhost:8085/testing?name=eason&address=xyz" --data 'name=ignore&address=ignore' -H "Content-Type:application/x-www-form-urlencoded"
====== Only Bind Query String ======
eason
xyz


*/

/*
package main

import "log"
import "github.com/gin-gonic/gin"

type Person struct {
	Name    string `form:"name" json:"name"`
	Address string `form:"address" json:"address"`
}

func main() {
	route := gin.Default()
	route.GET("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	if c.Bind(&person) == nil {
		log.Println("====== Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}

	if c.BindJSON(&person) == nil {
		log.Println("====== Bind By JSON ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}

	c.String(200, "Success")
}

*/
