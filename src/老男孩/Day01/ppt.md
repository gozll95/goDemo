# 开发环境搭建

## 1. 安装Go
a. 打开 址https://golang.org/dl/
b. 根据操作系统选择对应的安装包
c. 点击安装包进 安装(linux直接解压)
d. 设置环境变 (linux)
1. export GOROOT=$PATH:/path/to/go/
2. export PATH=$PATH:$GOROOT/bin/
3. export GOPATH=/home/user/project/go
e. 设置环境变 (window  设置)


## 2. IDE搭建(vscode)
a. 打开 址:https://code.visualstudio.com/ b. 根据操作系统选择对应的安装包
c. 点击安装包进 安装(linux直接解压)
d. 选择查看-》扩展-》搜索go，安装第 个


## 2. 新建项 
a. 新建 录/home/user/project/go/src/listen1
b.  vscode打开 录/home/user/project/go/src/listen1
c. 右键新建 件hello.go，保存
d. vscode会提示你安装 些go的 具，我们点击install all


## 3. 调试 具delve安装
a. 打开 址: https://github.com/derekparker/delve/tree/master/Documentation/installation
b. mac: brew install go-delve/delve/delve
c. linux&windows: go get github.com/derekparker/delve/cmd/dlv





# golang语 特性
## 1. 垃圾回收
- a. 内存 动回收，再也 需要开发 员管 内存
- b. 开发 员专注业务实现，降低  智负担 c. 只需要new分配内存， 需要释放
## 2. 天然并发
- a. 从语 层  持并发， 常简单
- b. goroute，轻 级线程，创建成千上万个goroute成为可能 
- c. 基于CSP(Communicating Sequential Process)模型实现
```
    func main() {
    go fmt.Println(“hello")
    }
```
## 3. channel
a. 管道，类似unix/linux中的pipe
b. 多个goroute之间通过channel进行通信 
c.  持任何类型
```
func main() {
pipe := make(chan int,3)
pipe <- 1
pipe <- 2 }
```


## 4. 多返回值
a.  个函数返回多个值
```
func calc(a int, b int)(int,int) { sum := a + b
    avg := (a+b)/2
    return sum, avg
}
```



#课后作业
1. 使 fmt分别打印字符 、 进制、 进制、 六进制、浮点数。