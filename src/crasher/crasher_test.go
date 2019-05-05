package crasher

import (
	"os"
	"os/exec"
	"testing"

	"github.com/astaxie/beego"
)

func TestCrasher(t *testing.T) {
	var a int
	if os.Getenv("BE_CRASHER") == "1" {
		beego.Error("xxx")
		a++
		t.Log("crasher")
		Crasher()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCrasher")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
