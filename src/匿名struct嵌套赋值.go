package main 

import(
	"fmt"
	"encoding/json"
)

// type Student struct {
   
// }

 type Account struct {
    Id uint32
    Name string
    Nested struct {
          Age uint8
     }
}

func main(){
	account := &Account { 
		Id : 10,
		Name : "jim",
		Nested : struct {Age uint8}{Age: 20},
	}
	

	fmt.Println(account)
	b, _ := json.Marshal(account)
	fmt.Println(string(b))
}

