Bind, BindJSON, BindQuery

ShouldBind, ShouldBindJSON, ShouldBindQuery


// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}


// Example for binding JSON ({"user": "manu", "password": "123"})
c.ShouldBindJSON(&json)



// Example for binding a HTML form (user=manu&password=123)
c.ShouldBind(&form)


type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

//与validation绑定
c.ShouldBindWith(&b, binding.Query)


#ShouldBindQuery
c.ShouldBindQuery(&person)




#Bind vs BindJson
only bind the query or post data

	if c.Bind(&person) == nil {
		log.Println("====== Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
only bind the json data

	if c.BindJSON(&person) == nil {
		log.Println("====== Bind By JSON ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}



#c.ShouldBind(&person)

// If `GET`, only `Form` binding engine (`query`) used.
// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48



curl -X POST "localhost:8085/testing?name=eason&address=xyz" --data 'name=ignore&address=ignore' -H "Content-Type:application/x-www-form-urlencoded"

ignore
ignore


 curl -X GET localhost:8085/testing --data '{"name":"JJ", "address":"xyz"}' -H "Content-Type:application/json"
 空


  curl -X GET "localhost:8085/testing?name=appleboy&address=xyz"
  appleboy
  xyz

