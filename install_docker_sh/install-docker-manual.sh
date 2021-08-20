#! /bin/sh
## use for centos7
set -e
remove_docker(){
    set +e
    echo "remove docker..."
    sudo yum remove docker \
        docker-client \
        docker-client-latest \
        docker-common \
        docker-latest \
        docker-latest-logrotate \
        docker-logrotate \
        docker-engine

    sudo yum remove docker-ce docker-ce-cli containerd.io
    sudo rm -rf /var/lib/docker
    sudo rm -rf /var/lib/containerd
    return 0
}

# https://download.docker.com/linux/centos/7/x86_64/stable/Packages/

# docker-ce 
# curl --connect-timeout 5 -O -L https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-19.03.13-3.el7.x86_64.rpm

# docker-ce-cli 
# curl --connect-timeout 5 -O -L https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-cli-19.03.13-3.el7.x86_64.rpm

# containerd.io 
# curl --connect-timeout 5 -O -L https://download.docker.com/linux/centos/7/x86_64/stable/Packages/containerd.io-1.3.7-3.1.el7.x86_64.rpm

install_docker(){
    echo "docker ce..."
    sudo yum install docker-ce-19.03.13-3.el7.x86_64.rpm
    echo "docker ce cli..."
    sudo yum install docker-ce-cli-19.03.13-3.el7.x86_64.rpm
    echo "install containerd.io..."
    sudo yum install containerd.io-1.3.7-3.1.el7.x86_64.rpm
}

# nvidia-docker need add repo
# sudo yum install --downloadonly -y nvidia-docker2-2.6.0 --downloaddir=nvidia-docker2

install_nvidia_docker(){
    echo "install nvidia docker..."
    nvidia_path="nvidia-docker2-2.6.0"
    for file in $(ls $nvidia_path)
    do  
    echo "install" $nvidia_path/$file;
    sudo yum install $nvidia_path/$file
    done  
}
start_docker(){
    echo "start docker..."
    sudo systemctl start docker
}

set_auto_start_docker(){
    echo "set auto start docker..."
    sudo systemctl enable docker.service
    sudo systemctl enable containerd.service
}

# docker-compose 
# curl --connect-timeout 5 -O -L https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64
install_docker_compose(){
    echo "install docker-compose..."
    sudo chmod 770 docker-compose-Linux-x86_64
    cp docker-compose-Linux-x86_64 /usr/local/bin/docker-compose
}

main(){
    # remove_docker
    install_docker
    install_nvidia_docker
    start_docker
    set_auto_start_docker
    install_docker_compose
    echo "all succeed."
    return 0
}

case $1 in
all)
    main
    ;;
remove_docker)
    remove_docker
    ;;
install_docker)
    install_docker
    ;;
install_nvidia_docker)
    install_nvidia_docker
    ;;
start_docker)
    start_docker
    ;;
set_auto_start_docker)    
    set_auto_start_docker
    ;;
install_docker_compose)
    install_docker_compose
    ;;
*)
    echo "unknow command"
    exit 1
    ;;
esac
