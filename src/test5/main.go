package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	sender := bufio.NewScanner(os.Stdin)
	for sender.Scan() {
		fmt.Println(sender.Text())
	}
}


rd := bufio.NewReader(r)
if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
name, err := rd.ReadString('\n')
