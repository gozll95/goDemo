并发其主要思想是使多个任务可以在同一时间执行以便能够更快的得到结果。并发编程的思想来自于多任务操作系统。

多任务操作系统允许同时运行多个程序。与之相对的是单用户计算机系统的操作系统，任务是被一个接一个的读取，寻找资源并运行的，各任务的运行完全是串行的。

并发程序内部会被划分为多个部分，每个部分都可以被看作是一个串行程序，在这些串行程序之间可能会存在交互的需求，这就需要操作系统去协调。在这之前，我们先来看下进程。

进程
我们通常把一个程序的执行称为一个进程，同时进程也被用来描述程序的执行过程。

一个进程可以使用系统调用fork创建若干新的进程。前者被称为后者的父进程，每一个进程都有父进程。所有的进程共同组成了一个树状结构，内核启动进程作为进程树的根并负责系统的初始化操作。它的父进程就是它自己。

为了管理进程，内核必须对每个进程的属性，行为进行详细的记录，包括进程的优先级，状态，虚拟地址范围以及各种访问权限等等。这些信息都会被记录在每个进程的进程描述符中，而被保存在进程描述符中的进程ID(常叫做PID)是进程在操作系统中的唯一标识，同时进程描述符中还会包含当前进程的父进程的ID(常被称为PPID)。

进程的状态共有6个，分别是可运行状态，可中断的睡眠状态，不可中断的睡眠状态，暂停状态或跟踪状态，僵尸状态和退出状态。

linux操作系统可以凭借cpu快速在多个进程之间切换，以产生多个进程在同时运行的假象。但切换正在运行的进程是需要付出代价的。

内核对进程的合理切换和调度使多个进程可以有条不紊的并发执行，在很多时候，多个进程之间需要相互配合并合作完成一个任务，这就需要进程间通讯机制(IPC)的支持。下面就讲一下go语言支持的IPC方法。它们是管道，信号和Socket。

管道
管道(pipe)是一种是单向的通讯方式。它只能被用于父进程与子进程以及同祖先的子进程之间的通讯。例如，我们在使用shell命令的时候常常会用到管道：

1
bogon:~ yang$ ps aux | grep go
shell命令为每个命令都创建一个进程，然后把左边的命令的标准输出用管道与右边的命令的标准输入连接起来。

管道的优点在于它的简单，而缺点则是只能单向通讯以及对通讯双方关系上的严格限制。

对于管道，go语言是支持的。通过标准库代码包os/exec中的API，我们可以执行操作系统命令并在此基础上建立管道。

package main
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)
func main() {
	demo1()
	fmt.Println()
	demo2()
}
func demo2() {
	fmt.Println("Run command `ps aux | grep apipe`: ")
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "apipe")
	stdout1, err := cmd1.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command: %s", err)
		return
	}
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error: The command can not running: %s\n", err)
		return
	}
	//使用带缓冲的读取器可以非常方便和灵活的读取需要的内容，而不是只能先把所有内容读出来再作处理
	outputBuf1 := bufio.NewReader(stdout1)
        //StdinPipe方法返回一个在命令Start后与命令标准输入关联的管道
	stdin2, err := cmd2.StdinPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdin pipe for command: %s\n", err)
		return
	}
	//WriteTo把所属值中缓存的数据全部写入到参数值代表的写入器中
	outputBuf1.WriteTo(stdin2)
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: The command can not be startup: %s\n", err)
		return
	}
	err = stdin2.Close()
	if err != nil {
		fmt.Printf("Error: Can not close the stdio pipe: %s\n", err)
		return
	}
	//wait会一直阻塞到所属命令完全运行结束为止
	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: Can not wait for the command: %s\n", err)
		return
	}
	fmt.Printf("%s\n", outputBuf2.Bytes())
}
func demo1() {
	useBufferIo := false
	fmt.Println("Run command `echo -n \"My first command from golang.\"`: ")
	cmd0 := exec.Command("echo", "-n", "My first command from golang.")
	//创建一个能获取此命令输出的管道
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command No.0: %s\n", err)
		return
	}
	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}
	if !useBufferIo {
		var outputBuf0 bytes.Buffer
		for {
			tempOutput := make([]byte, 5)
			//Read把读出的数据存入调用方传递给他的字节切片
			n, err := stdout0.Read(tempOutput)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Printf("Error: Can not read data from the pipe: %s\n", err)
					return
				}
			}
			if n > 0 {
				outputBuf0.Write(tempOutput[:n])
			}
		}
		fmt.Printf("%s\n", outputBuf0.String())
	} else {
		outputBuf0 := bufio.NewReader(stdout0)
		output0, _, err := outputBuf0.ReadLine()
		if err != nil {
			fmt.Printf("Error: Can not read data from the pipe: %s\n", err)
			return
		}
		fmt.Printf("%s\n", string(output0))
	}
}
上面我们讲的是匿名管道，与之相对的是命名管道。与匿名管道不同，任何进程都可以通过命名管道交换数据。实际上，命名管道以文件的形式存在于文件系统中，使用它的方法与使用文件类似，linux支持使用shell命令创建和使用命名管道，例如：


bogon:test yang$ mkfifo -m 644 myfifo2
bogon:test yang$ tee of_login < myfifo2 &
[1] 10028
bogon:test yang$ vi tepipi.txt
bogon:test yang$ cat tepipi.txt >myfifo2
[1]+  Done                    tee of_login < myfifo2
在上面的实例中，我们先使用命令mkfifo在当前目录创建了一个命名管道mififo2，然后又使用这个命名管道和命名tee把tepipe.txt文件中的内容写到了of_login文件中。

这里只是使用了命名管道搬运了数据，我们也可以在此基础上实现诸如数据的过滤或转换，以及管道的多路复用等功能。注意，命名管道默认是阻塞式的，更具体的说，只有在对这个命令管道的读操作和写操作都已准备就绪后，数据才会流转。它相对于匿名管道的优势就是通讯双方可以毫不相关。但命名管道也是单向的。

package main
import (
	"fmt"
	"io"
	"os"
	"time"
)
func main() {
	fileBasedPipe()
	inMemorySyncPipe()
}
func fileBasedPipe() {
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error: Can not create the named pipe: %s\n", err)
	}
	//命名管道默认会在其中一端还未就绪时阻塞另一端的进程
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Can not read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [file-based pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Can not write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)
	time.Sleep(200 * time.Millisecond)
}
func inMemorySyncPipe() {
	//与上面的两个命名管道不同，这两个是被存于内存中的，有原子性操作保证的管道
	reader, writer := io.Pipe()
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Can not read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [in-memory pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Can not write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [in-memory pipe]\n", n)
	time.Sleep(200 * time.Millisecond)
}
信号
它是IPC中唯一一种异步的通讯方法。它的本质是利用软件来模拟硬件的中断机制。信号被用来通知某个进程有某个事件发生了。使用kill命令查看当前系统支持的信号。

linux支持的信号有62种，分别分为两大类，1到31号为标准信号，也叫不可靠信号，34到64为实时信号，也叫可靠信号。

对同一进程来说，每种标准信号只会被记录并处理一次，并且如果某一进程的标准信号种类有好多，其处理顺序也是完全不确定的。而实时信号正好相反，即同种类的多个信号都可以被记录，并且可以按照发送的顺序被处理。

进程响应信号的方式有3种：忽略，捕捉和执行默认操作 .

linux对每个标准信号都有默认的操作方式。对大多数标准信号，我们可以自定义当进程接收到他们之后进行怎样的处理。在程序中，这些作为信号响应的自定义操作往往是由函数来代表的。

go命令会对其中的一些以键盘输入为来源的标准信号作出相应。这是由于go命令使用了在标准库代码包os/signal 中的被用于处理信号的API。

下面我们看下os.Signal接口类型的声明：

type Signal interface {
  String() string
  Signal() // to distinguish from other Stringers
}
从声明可知，其中的Signal()方法的声明并没有实际意义。它只是作为os.Signal接口类型的一个标识。

在标准库代码包syscall中，已经为不同的操作系统的所支持的每一个标准信号都声明了一个同名常量，其类型都为syscall.Signal——os.Signal接口类型的一个实现类型，同时也是一个int类型的别名类型。每个信号常量的整数值与他所代表的信号在操作系统中的编号一致。

代码包os/signal 中的Notify函数用来把操作系统发给当前进程的指定信号通知给该函数的调用方。声明如下：

1
func Notify(c chan<- os.Signal, sig ...os.Signal)
signal处理程序在向接受通道发送值的时候，并不会因为通道已满而产生阻塞。

前面说过，大部分的标准信号我们都可以自定义其处理方法，不过有两种信号除外。SIGKILL和SIGSTOP。对他们的响应只执行系统默认操作。

对于其他信号，我们可以自行处理也可以恢复对他们的系统默认操作，这需要使用到os/signal包中的Stop函数。声明如下：

1
func Stop(c chan<- os.Signal)
只需将Notify中的输入通道作为参数传入即可取消对这些信号的自行处理。

当然，我们也可以只对部分信号取消自定义处理，这时可以重新调用Notify函数，只需要第一个参数相同即可。

下面通过一个程序来实现以下功能：

1.执行一系列操作系统命令并获取演示进程的进程ID；

2.使用该进程值之上的API相对应的进程发送一个SIGINT信号，并输出演示进程已受到信号的凭证。

package main
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)
func main() {
	go func() {
		time.Sleep(5 * time.Second)
		sigSendingDemo()
	}()
	sigHandleDemo()
}
func sigHandleDemo() {
	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv1]\n", sigs1)
	signal.Notify(sigRecv1, sigs1...)
	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv2]\n", sigs2)
	signal.Notify(sigRecv2, sigs2...)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		//在sigRecv通道被关闭后，for语句会立即被退出执行
		for sig := range sigRecv1 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv1]\n")
		wg.Done()
	}()
	go func() {
		for sig := range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv2]\n")
		wg.Done()
	}()
	fmt.Println("Wait for 3 seconds... ")
	time.Sleep(3 * time.Second)
	fmt.Printf("Stop notification...")
	signal.Stop(sigRecv1)
	close(sigRecv1)
	fmt.Printf("done. [sigRecv1]\n")
	wg.Wait()
}
func sigSendingDemo() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal Error: %s\n", err)
			debug.PrintStack()
		}
	}()
	// ps aux | grep "mysignal" | grep -v "grep" | grep -v "go run" | awk '{print $2}'
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "mysignal"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{print $2}"),
	}
	output, err := runCmds(cmds)
	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}
	pids, err := getPids(output)
	if err != nil {
		fmt.Printf("PID Parsing Error: %s\n", err)
		return
	}
	fmt.Printf("Target PID(s):\n%v\n", pids)
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Process Finding Error: %s\n", err)
			return
		}
		sig := syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process (pid=%d)...\n", sig, pid)
		err = proc.Signal(sig)
		if err != nil {
			fmt.Printf("Signal Sending Error: %s\n", err)
			return
		}
	}
}
func getPids(strs []string) ([]int, error) {
	pids := make([]int, 0)
	for _, str := range strs {
		pid, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}
func runCmds(cmds []*exec.Cmd) ([]string, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("The cmd slice is invalid!")
	}
	first := true
	var output []byte
	var err error
	for _, cmd := range cmds {
		fmt.Printf("Run command: %v ...\n", getCmdPlaintext(cmd))
		if !first {
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}
		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}
		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}
		output = stdoutBuf.Bytes()
		//fmt.Printf("Output:\n%s\n", string(output))
		if first {
			first = false
		}
	}
	lines := make([]string, 0)
	var outputBuf bytes.Buffer
	outputBuf.Write(output)
	for {
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, getError(err, nil)
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}
func getCmdPlaintext(cmd *exec.Cmd) string {
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)
	for _, arg := range cmd.Args[1:] {
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}
	return buf.String()
}
func getError(err error, cmd *exec.Cmd, extraInfo ...string) error {
	var errMsg string
	if cmd != nil {
		errMsg = fmt.Sprintf("%s  [%s %v]", err, (*cmd).Path, (*cmd).Args)
	} else {
		errMsg = fmt.Sprintf("%s", err)
	}
	if len(extraInfo) > 0 {
		errMsg = fmt.Sprintf("%s (%v)", errMsg, extraInfo)
	}
	return errors.New(errMsg)
}