package main

import (
	"fmt"
	"npool/client"
	"strconv"
	"time"
)

var Job = []string{"aa", "bb", "cc", "dd", "ee", "ff"}

func main() {
	for i := 0; i < 1; i++ {
		c, err := client.NewClient("client" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}

		go func() {
			for {
				err = c.Write(Job)
				if err != nil {
					fmt.Println("err ====== is ", err)
				}

				time.Sleep(1 * time.Second)

				fmt.Println(c.GetConnCount())
				fmt.Println(c.GetHealthy())
			}
		}()
	}

	time.Sleep(100 * time.Second)

}
