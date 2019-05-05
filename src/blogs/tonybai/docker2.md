dockerapp3利用fork创建了一个子进程,这样dockerapp3实际上是两个进程在运行,各自有自己的signal监听goroutine,goroutine的处理逻辑是相同的。

利用dockerapp3，我们创建image "test:v3":
***$ sudo docker build -t="test:v3" ./ ***

启动基于test:v3 image的容器：

***$ sudo docker run -d "test:v3" ***
781cecb4b3628cb33e1b104ea57e506ad5cb4a44243256ebd1192af86834bae6
$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND                CREATED             STATUS              PORTS               NAMES
781cecb4b362        test:v3             "/bin/sh -c /bin/doc   5 seconds ago       Up 4 seconds                            insane_bohr      


通过docker logs查看dockerapp3的输出:
***$ sudo docker logs 781cecb4b362 ***
i am parent process, pid = 1
fork ok, childpid = 13
i am in child process, pid = 13

1: handle signal: terminated
1: signal termiate received, app exit normally

我们可以看到主进程收到了stop发来的SIGTERM并退出,主进程的退出导致容器退出,导致子进程13也无法生存,并且没有优雅退出。而在非容器状态下,子进程是可以被init进程接管的。

因此对于docker容器内运行的多进程程序,stop命令只会将SIGTERM发送给容器主进程,要想让其他进程也能优雅退出,需要在主进程与其他进程之间建立一种通信机制。在主进程退出前,等待其他子进程退出。待所有其他进程退出后,主进程再退出,容器停止。这样才能保证服务程序的优雅退出。

# 容器内启动多个服务程序
虽说docker best practice建议一个container内只放置一个服务程序，但对已有的一些遗留系统，在架构没有做出重构之前，很可能会有在一个 container中部署两个以上服务程序的情况和需求。而docker Dockerfile只允许执行一个CMD，这种情况下，我们就需要借助类似supervisor这样的进程监控管理程序来启动和管理container 内的多个程序了。

下面我们来自制作一个基于centos:centos6的安装了supervisord以及两个服务程序的image。我们将dockerapp1拷贝一份，并将拷贝命名为dockerapp1-brother。下面是我们的Dockerfile：

FROM centos:centos6
MAINTAINER Tony Bai <bigwhite.cn@gmail.com>
RUN yum install python-setuptools -y
RUN easy_install supervisor
RUN mkdir -p /var/log/supervisor
COPY ./supervisord.conf /etc/supervisord.conf
COPY ./dockerapp1 /bin/
COPY ./dockerapp1-brother /bin/
CMD ["/usr/bin/supervisord"]

通过容器的log，我们看出supervisord是等待两个程序退出后才退出的，不过我们还是要看看两个程序的输出日志以最终确认。重新启动容器，通过nsenter进入到容器中。


BTW，在物理机上测试supervisord以daemon形式运行，当kill掉supervisord时，supervisord是不会通知其监控 和管理的程序退出的。只有在以non-daemon形式运行时，supervisord才会在退出前先通知下面的程序退出。如果在一段时间内下面程序没有 退出，supervisord在退出前会kill -9强制杀死这些程序的进程。