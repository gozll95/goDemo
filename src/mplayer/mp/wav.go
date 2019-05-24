package mp

import (
	"fmt"
	"time"
)

type WAVPlayer struct {
	stat    int
	process int
	signal  chan int
}

func (p *WAVPlayer) Play(source string) {
	fmt.Println("Playing wav music", source)

	p.process = 0

	for p.process < 100 {
		time.Sleep(100 * time.Millisecond)
		fmt.Print(".")
		p.process += 10
	}
	fmt.Println("\nFinished playing", source)
}
