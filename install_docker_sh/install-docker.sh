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

# yum list docker-ce --showduplicates | sort -r

install_docker(){
    echo "install docker..."
    sudo yum install -y yum-utils
    sudo yum-config-manager \
        --add-repo \
        https://download.docker.com/linux/centos/docker-ce.repo
    sudo yum install docker-ce-19.03.13 docker-ce-cli-19.03.13 containerd.io-1.3.7 # install specific version
}

install_nvidia_docker(){
    echo "install nvidia docker..."
    distribution=$(. /etc/os-release;echo $ID$VERSION_ID) \
    && curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | sudo tee /etc/yum.repos.d/nvidia-docker.repo
    sudo yum install -y nvidia-docker2-2.6.0
}

set_docker_log_rotate(){
    echo "set docker log driver and rotate..."
    mkdir -p /etc/docker
    echo '{
    "log-driver": "local",
    "log-opts": {
        "max-size": "100m",
        "max-file": "3"
    },
    "runtimes": {
        "nvidia": {
            "path": "nvidia-container-runtime",
            "runtimeArgs": []
        }
    }
  }' > /etc/docker/daemon.json
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

# install docker-compose
install_docker_compose(){
    echo "install docker-compose..."
    curl --connect-timeout 5 -O -L https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64
    sudo chmod 770 docker-compose-Linux-x86_64
    cp docker-compose-Linux-x86_64 /usr/local/bin/docker-compose
}

main(){
    # remove_docker
    install_docker
    install_nvidia_docker
    set_docker_log_rotate
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
set_docker_log_rotate)
    set_docker_log_rotate
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
