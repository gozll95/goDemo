package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/zhu/qvm/server/tools/migrate_image/lib"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("fetch workspace dir failed %v", err)
		return
	}

	defer func() {
		os.Remove(path.Join(dir, "user_config.json"))
		// remove exe right
		os.Chmod(path.Join(dir, ".data"), 0000)

	}()

	// deal config files
	lib.DealConfig()

	os.Chmod(path.Join(dir, ".data"), 0777)

	cmd := exec.Command(path.Join(dir, ".data"))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// after 1 second remove temple config
	time.Sleep(time.Second)

	// remove temp config file
	os.Remove(path.Join(dir, "user_config.json"))

	// remove exe right
	os.Chmod(path.Join(dir, ".data"), 0000)

	lib.ReadPipeline(stdout)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}