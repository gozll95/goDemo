package crasher

import (
	"fmt"
	"log"
	"os"

	"github.com/astaxie/beego"
)

func Crasher() {
	log.Println("xxx")
	beego.Error("xxx")
	fmt.Println("Going down in flames!")
	os.Exit(1)
}
