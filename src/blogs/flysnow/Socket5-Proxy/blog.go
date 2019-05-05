前两天，使用Golang实现了一个简单的HTTP Proxy，具体实现参见 http://www.flysnow.org/2016/12/24/golang-http-proxy.html，这次使用Golang实现一个Socket5的简单代理。Socket5和HTTP并没有太大的不同，他们都可以完全给予TCP协议，只是请求的信息结构不同，所以这次我们不能像上次HTTP Proxy一样，解析请求和应答，要按照Socket的协议方式解析。

Socket协议版本
Socket协议分为Socket4和Socket5两个版本，他们最明显的区别是Socket5同时支持TCP和UDP两个协议，而SOcket4只支持TCP。目前大部分使用的是Socket5，我们这里只简单的介绍Socket5协议。

Socket5协议之授权认证
要想实现Socket5之间的连接会话，必须要懂SOcket5协议的实现细节和规范。这就好比我们都用普通话对话一样，彼此说的都明白，也可以给对方听得懂的回应。Socket5的客户端和服务端交流也一样，他们的语言就是Socket5协议。因为Socket5支持TCP和UDP两种，这里只介绍TCP这一种，UDP大同小异。

首先客户端会给服务端发送验证信息，这个是建立连接的前提。比如客户端：hi，哥们，借个火。服务端要认识它就说：好，给；如果不认识就说：你哪根葱啊！！客户端请求的暗号很简单：

VER	NMETHODS	METHODS
1	1	1 to 255
第一个字段VER代表Socket的版本，Soket5默认为0x05，其固定长度为1个字节
第二个字段NMETHODS表示第三个字段METHODS的长度，它的长度也是1个字节
第三个METHODS表示客户端支持的验证方式，可以有多种，他的尝试是1-255个字节。
目前支持的验证方式一共有：

X’00’ NO AUTHENTICATION REQUIRED（不需要验证）
X’01’ GSSAPI
X’02’ USERNAME/PASSWORD（用户名密码）
X’03’ to X’7F’ IANA ASSIGNED
X’80’ to X’FE’ RESERVED FOR PRIVATE METHODS
X’FF’ NO ACCEPTABLE METHODS（都不支持，没法连接了）
以上的都是十六进制常量，比如X’00’表示十六进制0x00。

服务端收到客户端的验证信息之后，就要回应客户端，服务端需要客户端提供哪种验证方式的信息。服务端的回应同样非常简洁。

VER	METHOD
1	1
第一个字段VER代表Socket的版本，Soket5默认为0x05，其值长度为1个字节
第二个字段METHOD代表需要服务端需要客户端按照此验证方式提供验证信息，其值长度为1个字节，选择为上面的六种验证方式。
举例说明，比如服务端不需要验证的话，可以这么回应客户端：

VER	METHOD
0x05	0x00
这就代表服务端说：哥们，我没啥要求，你来吧，我们使用Go实现代码如下：


var b [1024]byte
n, err := client.Read(b[:])
if err != nil {
	log.Println(err)
	return
}
if b[0] == 0x05 { //只处理Socket5协议
	//客户端回应：Socket服务端不需要验证方式
	client.Write([]byte{0x05, 0x00})
	n, err = client.Read(b[:])
}
我们这里以最简单的不需要验证的方式为例进行介绍，这种方式进行了上面的一问一答后就可以开始建立连接了。对于其他验证方式，还需要再进行一次一问一答，主要是客户端提供验证信息，服务端回应验证是否正确，这些细节可以参考http://www.ietf.org/rfc/rfc1928.txt以及http://www.ietf.org/rfc/rfc1929.txt的Socket5协议定义。

Socket5协议之建立连接。
Socket5的客户端和服务端进行双方授权验证通过之后，就开始建立连接了。连接由客户端发起，告诉Sokcet服务端客户端需要访问哪个远程服务器，其中包含，远程服务器的地址和端口，地址可以是IP4，IP6，也可以是域名。

VER	CMD	RSV	ATYP	DST.ADDR	DST.PORT
1	1	X’00’	1	Variable	2
VER代表Socket协议的版本，Soket5默认为0x05，其值长度为1个字节
CMD代表客户端请求的类型，值长度也是1个字节，有三种类型
CONNECT X’01’
BIND X’02’
UDP ASSOCIATE X’03’
RSV保留字，值长度为1个字节
ATYP代表请求的远程服务器地址类型，值长度1个字节，有三种类型
IP V4 address: X’01’
DOMAINNAME: X’03’
IP V6 address: X’04’
DST.ADDR代表远程服务器的地址，根据ATYP进行解析，值长度不定。
DST.PORT代表远程服务器的端口，要访问哪个端口的意思，值长度2个字节
从协议里解析我们需要的远程服务器信息，Go代码实现如下：

var host,port string
switch b[3] {
case 0x01://IP V4
	host = net.IPv4(b[4],b[5],b[6],b[7]).String()
case 0x03://域名
	host = string(b[5:n-2])//b[4]表示域名的长度
case 0x04://IP V6
	host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
}
port = strconv.Itoa(int(b[n-2])<<8|int(b[n-1]))
现在客户端把要请求的远程服务器的信息都告诉Socket5代理服务器了，那么Socket5代理服务器就可以和远程服务器建立连接了，不管连接是否成功等，都要给客户端回应，其回应格式为：

VER	REP	RSV	ATYP	BND.ADDR	BND.PORT
1	1	X’00’	1	Variable	2
VER代表Socket协议的版本，Soket5默认为0x05，其值长度为1个字节
REP代表响应状态码，值长度也是1个字节，有以下几种类型
X’00’ succeeded
X’01’ general SOCKS server failure
X’02’ connection not allowed by ruleset
X’03’ Network unreachable
X’04’ Host unreachable
X’05’ Connection refused
X’06’ TTL expired
X’07’ Command not supported
X’08’ Address type not supported
X’09’ to X’FF’ unassigned
RSV保留字，值长度为1个字节
ATYP代表请求的远程服务器地址类型，值长度1个字节，有三种类型
IP V4 address: X’01’
DOMAINNAME: X’03’
IP V6 address: X’04’
BND.ADDR表示绑定地址，值长度不定。
BND.PORT表示绑定端口，值长度2个字节
Go实现的连接和应答如下：


server, err := net.Dial("tcp", net.JoinHostPort(host, port))
if err != nil {
	log.Println(err)
	return
}
defer server.Close()
client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
绑定的地址和端口，这两个响应根据请求的CMD不同而不同，详细描述参考http://www.ietf.org/rfc/rfc1928.txt

数据转发
建立好连接之后，就是数据传递转发，TCP协议可以直接转发。UDP的话需要特殊处理，具体参考协议定义，其实就是一个特殊格式的回应，和上面的一问一答差不多，更多协议细节 http://www.ietf.org/rfc/rfc1928.txt。

TCP直接转发非常简单：

//进行转发
go io.Copy(server, client)
io.Copy(client, server)
完整代码实现
以下是完成的代码实现。


package main
import (
	"io"
	"log"
	"net"
	"strconv"
)
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
}
func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	if b[0] == 0x05 { //只处理Socket5协议
		//客户端回应：Socket服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return
		}
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		go io.Copy(server, client)
		io.Copy(client, server)
	}
}
这里目前是一个简易版本的Socket5 代理，还有很多没有实现，要实现完成的Socket5，可以自己根据他定义的协议试试。