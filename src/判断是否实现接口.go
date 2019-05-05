package main 

import(
	"fmt"
)
type Stringer interface {
    String() string
}

type A struct{
	AA string 
}

func(t *A)String()string{
	return t.AA
}
func main(){
	v:=A{"A"}

	if sv, ok := interface{}(&v).(Stringer); ok {
		fmt.Printf("v implements String(): %s\n", sv.String()) // note: sv, not v
	}

	var i Stringer
	i=&v
	fmt.Println(i.String())
}

