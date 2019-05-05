package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.DialTimeout("tcp", "8.8.8.8:80", 2*time.Second)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
}

/*
如果网络延迟较大，TCP握手过程将更加艰难坎坷（各种丢包），时间消耗的自然也会更长。Dial这时会阻塞，如果长时间依旧无法建立连接，则Dial也会返回“ getsockopt: operation timed out”错误。
在连接建立阶段，多数情况下，Dial是可以满足需求的，即便阻塞一小会儿。但对于某些程序而言，需要有严格的连接时间限定，如果一定时间内没能成功建立连接，程序可能会需要执行一段“异常”处理逻辑，为此我们就需要DialTimeout了。下面的例子将Dial的最长阻塞时间限制在2s内，超出这个时长，Dial将返回timeout error：
*/
