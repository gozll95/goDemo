# book-game

思路:

框架:
分为client/server两部分
client-->发送请求--->server 
	  <--接受相应<---

请求类型:
1.AddPlayer
2.ListPlayer
3.RemovePlayer

所以:
client需要有这些方法
client
	methods:
		- AddPlayer
		- ListPlayer
		- RemovePlayer 

所以client发出的请求格式应该为:
type Request struct{
	Method string 
	Params string 
}

相应的
server也需要有这些方法
server
	methods:
		- AddPlayer
		- ListPlayer
		- RemovePlayer

server根据client发送的Request解析出对应的func再去执行

client还要等待server的response,这里用channel来传输

client:
send Request to channel 
wait Response from channel


server:
wait Resquest from channel
send Response to channel

将server抽象成一个接口
type server interface{
	Name() string
	Handler(method,params string)*Response
}
那么这个传输通道就由new server的时候实现
new server 
	- make channel 
	- go func (for)
	- return channel

这样也不行,这样没办法将channel暴露出来作为client的一个属性

所以将channel加在server的一个方法上,那么server就不是接口了,而是包着interface的struct

type IpcServer struct{
	Server
}

IpcServer有个方法来暴露channel
method:
	- xxx() channel
		- make channel 
		- go func() for 
		- return channel

	
这个channel用来作client的属性,那么干脆:
type client struct{
	conn chan string
}
new client的时候将channel传进去

//总结下
IpcServer
(
	Server[
		Name()
		Handler(...) //处理不同的方法
	]
)[
	Connect() channel --------------
]									|
									|
IpcClient(				<--------------
	channel
)[
	AddPlayer---send AddPlayer+xxx to channel------->Server Handler()
	RemovePlayer---send RemovePlayer+xxx to channel------->Server Handler()
	ListPlayer---send ListPlayer+xxx to channel------->Server Handler()
]


//精华2:
交互式命令行:

//将命令和处理函数对应
func GetCommandHandlers() map[string]func(args []string) int {
	return map[string]func([]string) int{
		"help":       Help,
		"h":          Help,
		"quit":       Quit,
		"q":          Quit,
		"login":      Login,
		"logout":     Logout,
		"listplayer": ListPlayer,
		"send":       Send,
	}
}


func main() {
	fmt.Println("Casual Game Server Solution")

	startCenterService()

	Help(nil)

	r := bufio.NewReader(os.Stdin)

	handlers := GetCommandHandlers()

	for {
		//循环读取用户输入
		fmt.Println("Command> ")
		b, _, _ := r.ReadLine()
		line := string(b)

		tokens := strings.Split(line, " ")

		if handler, ok := handlers[tokens[0]]; ok {
			ret := handler(tokens)
			if ret != 0 {
				break
			} else {
				fmt.Println("Unknown command:", tokens[0])
			}
		}
	}
}

