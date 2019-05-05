#! /bin/bash
Usage(){
        echo "-i ip -h virtname"
}

I=1;
flagI=0
flagH=0
if [ $# -gt 0 ];then
        while [ $I -le $# ];do
                case $1 in
                -i)
                        if [ $flagI -eq 0 ];then
                            flagI=1
                            ip=$2
                            shift 2
                        else
                            Usage
                            exit 0
                        fi
                        ;;
                -h)
                        if [ $flagH -eq 0 ];then
                            flagH=1
                            virtname=$2
                            shift 2
                        else
                            Usage
                            exit 0
                        fi
                        ;;
                *)
                        Usage
                        exit 0
                        ;;
                esac
        done
fi

if [ -z $ip ] || [ -z $virtname ];then
    Usage
    exit 0
 fi

diskChoose(){
    disk=`df -TH | grep srv | awk '{print $NF"\t" $(NF-1)}'|sort -nk 2|head -n 1|awk '{print $1}'`
    echo -e "\033[41;36m $disk \033[0m"
}
ask(){
        read -p "do u want to rm? please input yes or no": option
        if [ -z $option ];then
                ask
        elif [ $option != "yes" -a  $option != "no" ];then
                ask
        else
                return 0
        fi
}
checkExist(){
    sudo virsh list --all | grep $virtname  && echo "have already!!" && exit 0
    if [ $option = "yes" ];then
       echo "nonononononononononono"
       exit 1
    fi
}
checkDns(){
    host $ip|grep $virtname || exit 1
}

changeFile(){
    sudo cp -rpf /tmp/example.qcow2 $disk/$virtname.qcow2
    sudo cp -rpf /tmp/virt.xml $disk/$virtname.xml
    sudo sed -i 's/virtname/'$virtname'/g' $disk/$virtname.xml
    sed -i -e '/source\ file/s#/srv/kvm/point#'$disk'#g' -e '/source\ file/s#image#'$virtname'.qcow2#g' $disk/$virtname.xml
}

changeEth0(){
    sudo virsh define $disk/$virtname.xml
    sudo virt-copy-out -d $virtname /etc/sysconfig/network-scripts/ifcfg-eth0 /tmp
    cat /tmp/ifcfg-eth0
    cat <<EOF > /tmp/ifcfg-eth0
DEVICE="eth0"
BOOTPROTO="static"
ONBOOT="yes"
TYPE="Ethernet"
IPADDR=10.120.82.138
NETMASK="255.255.254.0"
GATEWAY="10.120.82.1"
EOF
    sed -ri '/IPADDR/s/(.*)/IPADDR='$ip'/g' /tmp/ifcfg-eth0 #可以用cat eof方式生成随机文件
    sudo virt-copy-in /tmp/ifcfg-eth0 -d $virtname /etc/sysconfig/network-scripts/
    sudo virt-cat -d $virtname /etc/sysconfig/network-scripts/ifcfg-eth0
}

addOtherShell(){
    cat <<EOFFF > /tmp/rc.local
#! /bin/bash
mkdir /root/.ssh
  cat <<EOF > /root/.ssh/authorized_keys
ssh-rsa XXXX
EOF
chmod 600 /root/.ssh/authorized_keys
  cat <<EOFF > /etc/ssh/sshd_config
Port 22
Port 1046
Port 18211
Protocol 2
ChallengeResponseAuthentication no
PasswordAuthentication no
X11Forwarding yes
PrintMotd no
AcceptEnv LANG LC_*
Subsystem sftp /usr/libexec/openssh/sftp-server
UsePAM yes
LogLevel VERBOSE
UseDNS no
EOFF
/etc/init.d/sshd restart
EOFFF
    sudo virt-copy-in /tmp/rc.local -d $virtname /etc/rc.d/
    #sudo virt-cat -d $virtname /etc/rc.local
}

sudo dpkg -l | grep guestfish || sudo apt-get install -y guestfish
diskChoose
checkExist.
checkDns
changeFile
changeEth0
addOtherShell

#virsh start $virtname && virsh console $virtname
sudo virsh start $virtname

echo "Info: $virtname init done!"
exit 0