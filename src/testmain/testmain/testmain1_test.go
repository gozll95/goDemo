package testmain

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	arg := os.Getenv("PKGDIR")
	log.Println("xxxxxxxxxx:", os.Getenv("PKGDIR"))

	if arg != "" {
		exitVal := m.Run()
		os.Exit(exitVal)
	}
	os.Exit(-1)

}

func Test1(t *testing.T) {
	log.Println("[Test1] running ")
}
