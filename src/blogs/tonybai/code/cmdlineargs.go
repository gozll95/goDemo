//cmdlineargs.go
package main
import (
    //      "fmt"
    "os"
    "path/filepath"
)
func main() {
    println("I am ", os.Args[0])
    baseName := filepath.Base(os.Args[0])
    println("The base name is ", baseName)
    // The length of array a can be discovered using the built-in function len
    println("Argument # is ", len(os.Args))
    // the first command line arguments
    if len(os.Args) > 1 {
        println("The first command line argument: ", os.Args[1])
    }
}

/*
$go build cmdlineargs.go
$cmdlineargs test one
I am  cmdlineargs
The base name is  cmdlineargs
Argument # is  3
The first command line argument:  test
*/