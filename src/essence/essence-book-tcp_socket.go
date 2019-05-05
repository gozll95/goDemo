# book-tcp_socket


// 现象:
/*
Server[0]: Got listener for the server. (local address: 127.0.0.1:8085)
Client[1]: Connected to server. (remote address: 127.0.0.1:8085, local address: 127.0.0.1:54276)
Server[0]: Established a connection with a client application. (remote address: 127.0.0.1:54276)
Client[1]: Sent request (written 11 bytes): 1298498081.
Client[1]: Sent request (written 11 bytes): 2019727887.
Client[1]: Sent request (written 11 bytes): 1427131847.
Client[1]: Sent request (written 10 bytes): 939984059.
Client[1]: Sent request (written 10 bytes): 911902081.
Server[0]: Received request: 1298498081.
Server[0]: Sent response (written 44 bytes): The cube root of 1298498081 is 1090.972418..
Server[0]: Received request: 2019727887.
Server[0]: Sent response (written 44 bytes): The cube root of 2019727887 is 1264.050100..
Client[1]: Received response: The cube root of 1298498081 is 1090.972418..
Server[0]: Received request: 1427131847.
Server[0]: Sent response (written 44 bytes): The cube root of 1427131847 is 1125.869444..
Server[0]: Received request: 939984059.
Server[0]: Sent response (written 42 bytes): The cube root of 939984059 is 979.580571..
Server[0]: Received request: 911902081.
Client[1]: Received response: The cube root of 2019727887 is 1264.050100..
Server[0]: Sent response (written 42 bytes): The cube root of 911902081 is 969.726809..
Client[1]: Received response: The cube root of 1427131847 is 1125.869444..
Client[1]: Received response: The cube root of 939984059 is 979.580571..
Client[1]: Received response: The cube root of 911902081 is 969.726809..
Server[0]: The connection is closed by another side.
*/

// 开启 服务端 goroutine
for accept connection and go handle connection
go serverGo():
	listener,err:=net.Listen(xx,xx)
	for{
		conn,err:=listener.Accept() //阻塞直至连接到来
		go handleConn(conn)
			- defer func(){
				- conn.Close()
				- // else
			}()
			- for {
				- conn.SetReadDeadline(time.Now().Add(10 * time.Second))
				// 这里的read里面包括for,会将数据读完
				- read(conn)
				- // 处理数据
				- // 这里的write就write
				- write(conn,xxx)
			}
	}


// 开启客户端goroutine
go clientGo(1):
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	defer conn.Close()
	// 发送请求
	n, err := write(conn, fmt.Sprintf("%d", req))
	// 接收请求
	strResp, err := read(conn)


1."注意点":
read(conn)的时候判断边界:
strResp, err := read(conn)
if err != nil {
	if err == io.EOF {
		printClientLog(id, "The connection is closed by another side.")
	} else {
		printClientLog(id, "Read Error: %s", err)
	}
	break
}


2."write":

buffer.Write
conn.Write(buffer.Bytes())


func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}


3."read":
// 千万不要使用这个版本的read函数！
//func read(conn net.Conn) (string, error) {
//	reader := bufio.NewReader(conn)
//	readBytes, err := reader.ReadBytes(DELIMITER)
//	if err != nil {
//		return "", err
//	}
//	return string(readBytes[:len(readBytes)-1]), nil
//}


// 因为TCP read是分段的,并不能保证一次性能读到多少
func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}
// for read a bytes to buffer then return buffer string

4."log":
func printLog(role string, sn int, format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}

func printServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

func printClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}



5."思路":
server------connection-------client
server探测到一个connection就会开启一个handle function()

server handle function():
	for 一直从 connection 里去 read(这里read是一条完整的数据,以DELIEM为结束符,所以这里read里也有一个for去一直读,直到读到DELIEM)->处理->write to connection


client connection:
	- 发送请求 write to connection
	- read from connection



