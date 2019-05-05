近期在试验如何将我们的产品部署到docker容器中去,这其中涉及到一个技术环节,那就是如何让docker容器退出时其内部运行的服务程序也可以优雅的退出。所谓优雅退去,指的就是程序在退出前有清理资源(比如关闭文件描述符、关闭socket),保存必要中间状态,持久化内存数据(比如将内存中的数据flush到文件中)的机会。docker作为目前最火的轻量级虚拟化技术,其在后台服务领域的应用是极其广泛的,其设计者在程序优雅退出方面是有考虑的。下面我们由简单到复杂逐一考量一下。

# 优雅退出的原理
对于服务程序而言，一般都是以daemon形式运行在后台的。通知这些服务程序退出需要使用到系统的signal机制。一般服务程序都会监听某个 特定的退出signal，比如SIGINT、SIGTERM等（通过kill -l命令你可以查看到几十种signal）。当我们使用kill + 进程号时，系统会默认发送一个SIGTERM给相应的进程。该进程通过signal handler响应这一信号，并在这个handler中完成相应的“优雅退出”操作。

与“优雅退出”对立的是“暴力退出”，也就是我们常说的使用kill -9，也就是kill -s SIGKILL + 进程号，这个行为不会给目标进程任何时间空隙，而是直接将进程杀死，无论进程当前在做何种操作。这种操作常常导致“不一致”状态的出现。***SIGKILL这 个信号比较特殊***，进程无法有效监听该信号，无法有效针对该信号设置handler，无法改变其信号的默认处理行为。


# 测试用"服务程序"

sudo docker pull centos:centos6
sudo docker run -t -i centos:centos6 /bin/bash

制作我们的第一个测试用image。安装官方推荐的Best Practice,我们使用Dockerfile来build一个测试用image。
步骤如下:

- 建立~/ImagesFactory目录
- 将构建好的dockerapp1拷贝到~/ImagesFactory目录下
- 进入~/ImagesFactory目录,创建DockerFile文件,

***Dockerfile***内容如下:

FROM centos:centos6
MAINTAINER Tony Bai <bigwhite.cn@gmail.com>
COPY ./dockerapp1 /bin/
CMD /bin/dockerapp1

执行docker build,结果如下:
sudo docker build -t="test:v1" ./

$ sudo docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
test                v1                  21a720a808a7        59 seconds ago      214.6 MB

# 第一个测试容器
我们基于image "test:v1"启动第一个测试容器:
$ sudo docker run -d "test:v1"
daf3ae88fec23a31cde9f6b9a3f40057953c87b56cca982143616f738a84dcba

$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND                CREATED             STATUS              PORTS               NAMES
daf3ae88fec2        test:v1             "/bin/sh -c /bin/doc   17 seconds ago      Up 16 seconds                           condescending_sammet  

通过docker run命令，我们基于image"test:v1"启动了一个容器。通过docker ps命令可以看到容器成功启动，容器id：daf3ae88fec2，别名为：condescending_sammet。

根据Dockerfile我们知道，容器启动后将执行"/bin/dockerapp1"这个程序，dockerapp1退出，容器即退出。 run命令的"-d"选项表示容器将以daemon的形式运行，我们在前台无法看到容器的输出。那么我们怎么查看容器的输出呢？我们可以通过 docker logs + 容器id的方式查看容器内应用的标准输出或标准错误。我们也可以进入容器来查看。


进入容器有多种方法，比如用***sudo docker attach daf3ae88fec2***。attach后，就好比将daemon方式运行的容器 拿到了前台，你可以Ctrl + C一下，可以看到如下dockerapp1的输出:

^Chandle signal: interrupt

另外一种方式是利用nsenter工具进入我们容器的namespace空间。
另外一种方式是利用nsenter工具进入我们容器的namespace空间。ubuntu 14.04下可以通过如下方式安装该工具：

$ wget https://www.kernel.org/pub/linux/utils/util-linux/v2.24/util-linux-2.24.tar.gz; tar xzvf util-linux-2.24.tar.gz
$ cd util-linux-2.24
$ ./configure –without-ncurses && make nsenter
$ sudo cp nsenter /usr/local/bin


echo $(sudo docker inspect -format "{{.State.Pid}}") daf3ae88fec2)
5494

sudo nsenter -target 5494 -mount -uts -ipc -net -pid
-bash-4.1# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 09:20 ?        00:00:00 /bin/dockerapp1
root        16     0  0 09:32 ?        00:00:00 -bash
root        27    16  0 09:32 ?        00:00:00 ps -ef
-bash-4.1#

进入容器后通过ps命令可以看到正在运行的dockerapp1程序。在容器内，我们可以通过kill来测试dockerapp1的运行情况：

-bash-4.1# kill -s SIGINT 1

通过前面的attach窗口，我们可以看到dockerapp1输出:

handle signal: interrupt

如果你发送SIGTERM信号，那么dockerapp1将终止运行，容器也就停止了。

***-bash-4.1# kill 1***

attach窗口显示：

***signal termiate received, app exit normally***

我们可以看到容器启动后默认执行的时Dockerfile中的CMD命令，如果Dockerfile中有多行CMD命令，Docker在启动容器 时只会执行最后一条CMD命令。如果在docker run中指定了命令，docker则会执行命令行中的命令而不会执行dockerapp1，比如：

$ sudo docker run -t -i "test:v1" /bin/bash
bash-4.1#

这里我们看到直接执行的时bash，dockerapp1并未执行。


# docker stop的行为
我们先来看看docker stop的manual:
sudo docker stop --help
Usage: docker stop [OPTIONS] CONTAINER [CONTAINER...]
Stop a running container by sending SIGTERM and then SIGKILL after a grace period
  -t, –time=10      Number of seconds to wait for the container to stop before killing it. Default is 10 seconds.

可以看出当我们执行docker stop时,docker会首先向容器内的主程序发送一个SIGTERM信号,用于容器内程序的退出。如果容器在收到SIGTEM后没有马上退出,那么stop命令会在等待一段时间(默认是10s)后,再向容器发送SIGKILL信号,将容器杀死,变为退出状态。



我们来验证一下docker stop的行为。启动刚才那个容器：

$ sudo docker start daf3ae88fec2
daf3ae88fec2

attach到容器daf3ae88fec2
$ sudo docker attach daf3ae88fec2

新打开一个窗口，执行docker stop命令：
$ sudo docker stop daf3ae88fec2
daf3ae88fec2

可以看到attach窗口输出：
handle signal: terminated
signal termiate received, app exit normally

通过docker ps查看，发现容器已经退出。

也许通过上面的例子还不能直观的展示stop命令的两阶段行为，因为dockerapp1收到SIGTERM后直接就退出 了，stop命令无需等待容器慢慢退出，也无需发送SIGKILL。我们改造一下dockerapp1这个程序。

我们复制一下dockerapp1.go为dockerapp2.go，编辑dockerapp2.go，将handler中对SIGTERM的 处理注释掉，其他不变：

handler := func(s os.Signal, arg interface{}) {
                fmt.Printf("handle signal: %v\n", s)
                /*
                if s == syscall.SIGTERM {
                        fmt.Printf("signal termiate received, app exit normally\n")
                        os.Exit(0)
                }
                */
        }

我们使用dockerapp2来构建一个新image：test:v2，将Dockerfile中得dockerapp1换成 dockerapp2即可。

$ sudo docker build -t="test:v2" ./
Sending build context to Docker daemon 9.369 MB
Sending build context to Docker daemon
Step 0 : FROM centos:centos6
 —> 68edf809afe7
Step 1 : MAINTAINER Tony Bai <bigwhite.cn@gmail.com>
 —> Using cache
 —> c617b456934a
Step 2 : COPY ./dockerapp2 /bin/
 —> 27cd613a9bd7
Removing intermediate container 07c760b6223b
Step 3 : CMD /bin/dockerapp2
 —> Running in 1aac086452a7
 —> 82eb876fefd2
Removing intermediate container 1aac086452a7
Successfully built 82eb876fefd2

利用image "test:v2"创建一个容器来测试stop。

$ sudo docker run -d "test:v2"
29f3ec1af3c355458cbbd802a5e8a53da28e9f51a56ce822c7bba2a772edceac

$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND                CREATED             STATUS              PORTS               NAMES
29f3ec1af3c3        test:v2             "/bin/sh -c /bin/doc   7 seconds ago       Up 6 seconds                            romantic_feynman   

Attach到这个容器并观察，在另外一个窗口stop该container。我们在attach窗口只看到如下输出：

handle signal: terminated

stop命令的执行没有立即返回，而是等待容器退出。等待10s后，容器退出，stop命令执行结束。从这个例子我们可以明显看出stop的两阶 段行为。

如果我们以sudo docker run -i -t "test:v1" /bin/bash形式启动容器，那stop命令会将SIGTERM发送给bash这个程序，即使你通过nsenter进入容 器，启动了dockerapp1，dockerapp1也不会收到SIGTERM，dockerapp1会随着容器的退出而被强行终止，就像被 kill -9了一样。

