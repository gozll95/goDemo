# book-chatbot


1、"思路":
多种类型的聊天机器人
共性: 说话/关闭/else... -> "可以抽象为一个接口"
- Name()string 
- Begin() (string, error) //请输入你的名字: / please input your name:
- ReportError(err error) string
- End() error // 再见/bye

- 询问帮助
- 当获得指定词汇就做指定动作

以上两种方法,做成:可以传入/可以定制
我们叫method A/method B
如果可以传入,那么就要成为chatbot的一个传入参数C,C里有method A/method B
所以C应该也是一个interface

NewChatbot(xxx,C C-Interface{})

(n *NewChatbot) method A:
			- if n.C !=nil -> use n.C.A 
			- else: -> 自定义

(n *NewChatbot) method B:
			- if n.C !=nil -> use n.C.B 
			- else: -> 自定义	


多种类型的机器人-->注册进机器人池(it is a map)-->get机器人by name -> ... 


2. "流程"
- 使用flag问用户chatbotName?simple.cn/simple.en 
  // 向chatbot注册simple.cn/simple.en机器人
- chatbot.Register(chatbot.NewSimpleEN("simple.en", nil))
- chatbot.Register(chatbot.NewSimpleCN("simple.cn", nil))
  // 利用用户输入的chatbotName从chatbot里获取对应的chatbot
- myChatbot:= chatbot.Get(chatbotName)
- 用了一些方法... 
- begin, err := myChatbot.Begin()
- myChatbot.Talk(input)
- myChatbot.End()
- chatbot.ReportError(err)


3."主要数据结构和方法":
接口:
// Chatbot 定义了聊天机器人的接口类型。
type Chatbot interface {
	Name() string
	Begin() (string, error)
	Talk
	ReportError(err error) string
	End() error
}

// Talk 定义了聊天的接口类型。
type Talk interface {
	Hello(userName string) string
	Talk(heard string) (saying string, end bool, err error)
}


4."机器人池":
var chatbotMap = map[string]Chatbot{}

方法:
	- 注册机器人 
	- 根据机器人名字获取相应的机器人


5."机器人":
type simpleCN struct {
	name string
	talk Talk
}

new(name x,talk Talk)

6."错误":
func checkError(chatbot chatbot.Chatbot, err error, exit bool) bool {
	if err == nil {
		return false
	}
	if chatbot != nil {
		fmt.Println(chatbot.ReportError(err))
	} else {
		fmt.Println(err)
	}
	if exit {
		debug.PrintStack()
		os.Exit(1)
	}
	return true
}










