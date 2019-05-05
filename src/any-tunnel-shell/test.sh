#! /bin/bash

#  命令:
#  bash xx.sh \
#   -hostname "beijing-any-tunnel-2" \
#   -up443 "xxx.com" -up443 "xxx.com" \
#   -up80 "xxx.com" -up80 "xxx.com" \
#   -io443 "qvm-z1.zhu.me"  \
#   -io80 "~.*" \
#   -backend 10.0.1.131:80 -backend 10.0.1.130:80 \
#   -nginx_dir "/root/nginx-template"


# 全局变量
HOSTNAME=""
IP=""
UP_HTTPS_SERVER_NAME=("")
UP_HTTP_SERVER_NAME=("")
IO_HTTPS_SERVER_NAME=("")
IO_HTTP_SERVER_NAME=("")
UPSTREAMS=("")
NGINX_PATH=""

UP80=""
UP443=""
IO80=""
IO443=""
BACKENDS=""


# 使用
Usage(){
    echo -e "usage:
             -hostname hostname
             -io443 io443.domain -io443 io443.domain
             -io80 io80.domain -io80 io80.domain
             -up443 up443.domain -up443 up443.domain
             -up80 up80.domain -up80 up80.domain
             -backend 1.1.1.1:80 -backend 2.2.2.2:80
             -nginx_dir /tmp/nginx
             "
}

PrintArgs(){
    echo -e "args:
            hostname: $HOSTNAME
            up http server_name: $UP80
            up https server_name: $UP443
            io http server_name: $IO80
            io https server_name: $IO443
            backends: $BACKENDS
            nginx_template: $NGINX_PATH"
}
CheckArgs(){
    # 判断参数有效性
    echo "HOSTNAME is $HOSTNAME"
    if [ "$HOSTNAME"x = ""x ];then
        echo "hostname is invalid"
        Usage
        exit 1
    fi

    if [ ${#UP_HTTP_SERVER_NAME[*]} -eq 1 ] && [ "${UP_HTTP_SERVER_NAME[0]}"x = ""x ];then 
        echo "nginx up http server_name is invalid"
        Usage
        exit 1
    fi 
    if [ ${#UP_HTTPS_SERVER_NAME[*]} -eq 1 ] && [ "${UP_HTTPS_SERVER_NAME[0]}"x = ""x ];then 
        echo "nginx up https server_name is invalid"
        Usage
        exit 1
    fi 
    if [ ${#IO_HTTPS_SERVER_NAME[*]} -eq 1 ] && [ "${IO_HTTPS_SERVER_NAME[0]}"x = ""x ];then 
        echo "nginx io https server_name is invalid"
        Usage
        exit 1
    fi 
    if [ ${#IO_HTTP_SERVER_NAME[*]} -eq 1 ] && [ "${IO_HTTP_SERVER_NAME[0]}"x = ""x ];then 
        echo "nginx io http server_name is invalid"
        Usage
        exit 1
    fi
    if [ ${#UPSTREAMS[*]} -eq 1 ];then 
        echo "nginx upstream is invalid"
        Usage
        exit 1
    fi 
    if [ "$NGINX_PATH"x = ""x ] || [ ! -d $NGINX_PATH ];then
        echo "nginx template path is invalid"
        Usage
        exit 1
    fi 

    # 变量初始化
    for v in ${IO_HTTP_SERVER_NAME[@]}
    do
        IO80="$IO80  $v"
    done

    for v in ${IO_HTTPS_SERVER_NAME[@]}
    do
        IO443="$IO443  $v"
    done

    for v in ${UP_HTTP_SERVER_NAME[@]}
    do
        UP80="$UP80 $v"
    done

    for v in ${UP_HTTPS_SERVER_NAME[@]}
    do
        UP443="$UP443 $v"
    done

    for v in ${UPSTREAMS[@]}
    do
        BACKENDS="$BACKENDS $v"
    done

    echo "args is valid"
    PrintArgs
}


# 1.参数输入
Input(){
    I=1;
    flagHostname=0
    flagIo443=0
    flagIo80=0
    flagUp443=0
    flagUp80=0
    flagNginxPath=0

    if [ $# -gt 0 ];then
            while [ $I -le $# ];do
                    case $1 in
                    -hostname)
                            if [ $flagHostname -eq 0 ];then
                                flagHostname=1
                                HOSTNAME=$2
                                shift 2
                            else
                                Usage
                                exit 0
                            fi
                            ;;
                    -io443)
                            io443=( $2 )
                            IO_HTTPS_SERVER_NAME+=("${io443[@]}")
                            shift 2
                            ;;
                    -io80)
                            io80=( $2 )
                            IO_HTTP_SERVER_NAME+=("${io80[@]}")
                            shift 2
                            ;;
                    -up443)
                            up443=( $2 )
                            UP_HTTPS_SERVER_NAME+=("${up443[@]}")
                            shift 2
                            ;;
                    -up80)
                            up80=( $2 )
                            UP_HTTP_SERVER_NAME+=("${up80[@]}")
                            shift 2
                            ;;
                    -nginx_dir)
                            if [ $flagNginxPath -eq 0 ];then
                                flagNginxPath=1
                                NGINX_PATH=$2
                                shift 2
                            else
                                Usage
                                exit 0
                            fi
                            ;;
                    -backend)  
                            backend=( $2 )
                            UPSTREAMS+=("${backend[@]}")
                            shift 2
                            ;;
                    *)
                            Usage
                            exit 0
                            ;;
                    esac
            done
    fi
}



# 2.修改主机名
ModifyHostName(){             
    getIp
    echo $HOSTNAME > /etc/hostname

    cat <<EOF > /etc/hosts
127.0.0.1	localhost
$IP $HOSTNAME
EOF

    hostname -F /etc/hostname 
    echo name is `hostname -f`
}


# 3.修改sshd_config
ModifySSHPort(){
    cat /etc/ssh/sshd_config| grep -q 'Port 18922' && return 
    echo "begin set ssh port to 18922"
    sed -i 's/Port 22/Port 18922/g' /etc/ssh/sshd_config
    /etc/init.d/ssh restart
    echo "restarted ssh daemon"
}


# 4.修改apt源
ModifyApt(){
    cat <<EOF > /etc/apt/sources.list
deb http://mirrors.163.com/ubuntu/ trusty main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ trusty-security main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ trusty-updates main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ trusty-proposed main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ trusty-backports main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ trusty main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ trusty-security main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ trusty-updates main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ trusty-proposed main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ trusty-backports main restricted universe multiverse

deb http://nginx.org/packages/mainline/ubuntu/ trusty nginx
deb-src http://nginx.org/packages/mainline/ubuntu/ trusty nginx
EOF
    apt-get update
}

# 5.修改内核参数
AddSysctl(){
    cat <<EOF >> /etc/sysctl.conf
    # use kodo
    kernel.unknown_nmi_panic=1
    kernel.sysrq=1
    fs.file-max=655360
    net.ipv4.ip_local_port_range=2048 65000
    net.ipv4.ip_local_reserved_ports=1024-10000
    net.ipv4.tcp_tw_recycle=1
    net.ipv4.tcp_tw_reuse=1
    net.ipv4.tcp_max_orphans=262144
    net.ipv4.tcp_syn_retries=1
    net.ipv4.tcp_fin_timeout=30
    net.ipv4.tcp_keepalive_time=600
    net.ipv4.tcp_timestamps=0
    net.core.somaxconn=655350
    net.core.netdev_max_backlog=262144
    vm.swappiness=5
    vm.dirty_background_bytes=104857600
EOF
    sysctl -p
}


# 6.nginx
Nginx(){
    # 判断nginx是否开启
    echo `ps aux | grep nginx`| grep -q worker 
    if [ $? -eq 0 ];then 
         echo "nginx already running" 
         return 
    fi

    # 判断是否安装了nginx
    dpkg -l | grep nginx| grep -q '^ii'
    if [ $? -eq 0 ];then 
        echo "nginx already installed"
        echo `nginx -v`
        # 如果不是1.13以上版本则删除nginx
        dpkg -l | grep nginx| grep -o '1.1[3-9]' || removeNginx
    else
        echo "begin install nginx" 
        wget https://nginx.org/keys/nginx_signing.key
        apt-key add nginx_signing.key
        apt-get remove nginx-common -y
        apt-get update
        apt-get install nginx -y
        dpkg -l | grep nginx| grep '^ii' && echo "install nginx successfully" || (echo "install nginx failed" && exit 1)
        echo `nginx -v`
    fi

    # 设置nginx config
    setNginxConfig

    # 开启nginx
    nginx -t
    /etc/init.d/nginx start 
}




# 删除现有nginx
removeNginx(){
    echo "begin remove nginx"
    dpkg -l | grep nginx 
    apt-get autoremove -y nginx
    apt-get purge -y nginx
    echo "remove nginx successfully"
}

setNginxConfig(){
    echo "start set nginx config from template" 

   # 准备工作
   test -e /disk1/nginx || mkdir -p /disk1/nginx
   test -e ~/nginx_init_bak || mkdir ~/nginx_init_bak
   cp -rpf /etc/nginx/* ~/nginx_init_bak  
   rm -rf /etc/nginx/*
   cp -rpf $NGINX_PATH/* /etc/nginx 

   sed -i "s/UP_HTTPS_SERVER_NAME/$UP443/g" /etc/nginx/sites-enabled/up_proxy.conf
   sed -i "s/UP_HTTP_SERVER_NAME/$UP80/g" /etc/nginx/sites-enabled/up_proxy.conf
   sed -i "s/IO_HTTPS_SERVER_NAME/$IO443/g" /etc/nginx/sites-enabled/io_proxy.conf
   sed -i "s/IO_HTTP_SERVER_NAME/$IO80/g" /etc/nginx/sites-enabled/io_proxy.conf

   rm -rf /etc/nginx/sites-enabled/upstream.conf && touch /etc/nginx/sites-enabled/upstream.conf
   echo 'upstream backend_hosts {' >>  /etc/nginx/sites-enabled/upstream.conf

    for upstream in ${UPSTREAMS[@]}  
    do  
        echo "server ${upstream}  max_fails=0 ;" >> /etc/nginx/sites-enabled/upstream.conf
    done 

    echo '}' >> /etc/nginx/sites-enabled/upstream.conf 
}


getIp(){
    IP=`ip a show eth0| grep inet| awk '{print $2}'|awk -F '/' '{print $1}'`
}


# 主函数
main(){
    Input $@
    CheckArgs
    ModifyApt
    ModifyHostName
    ModifySSHPort
    AddSysctl
    Nginx
}

main $@


