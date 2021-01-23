
FROM ubuntu:16.04

MAINTAINER 3DS Interactive (contact@3dsinteractive.com)

ENV HOME /root

# install
# - redis-cli
# - net tools
# - mysql client
# - python3
# - aws-cli
# - telnet
# - git
# - wrk2
# - golang

# install requirements
RUN apt-get update
RUN apt-get -y install software-properties-common
RUN apt-get -y install build-essential
RUN apt-get -y install zlib1g-dev
RUN apt-get install -y git

# install python3.6
RUN add-apt-repository ppa:deadsnakes/ppa
RUN apt-get update
RUN apt-get -y --fix-missing install python3.6 \
    && update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.6 1

# install prerequisite
RUN apt-get update
RUN apt-get -y install traceroute 
RUN apt-get -y install curl 
RUN apt-get -y install dnsutils 
RUN apt-get -y install netcat-openbsd 
RUN apt-get -y install jq 
RUN apt-get -y install nmap 
RUN apt-get -y install net-tools 
RUN apt-get -y install iputils-ping
RUN apt-get -y install redis-tools
RUN apt-get -y install mysql-client-5.7
RUN apt-get -y install libsqlite3-dev
RUN apt-get -y install sqlite3
RUN apt-get -y install bzip2 
RUN apt-get -y install libbz2-dev

# install kafkacat
RUN apt-get -y install ca-certificates
RUN apt-get -y install kafkacat

# install pip
RUN curl -sS https://bootstrap.pypa.io/get-pip.py | python3
# install aws-cli (required python3 & pip installed above)
RUN pip3 install awscli --upgrade --user
# PATH for python
ENV PATH="/root/.local/bin:${PATH}"

# install wrk2
RUN apt-get -y install libssl-dev 
RUN git clone https://github.com/giltene/wrk2.git /home/root/wrk2 &&\
    cd /home/root/wrk2 &&\
    make &&\
    cp wrk /usr/local/bin &&\
    rm -Rf /home/root/wrk2

# Install Lua and Luarocks - a lua package manager
RUN apt-get -y install lua5.1
RUN apt-get -y install liblua5.1-dev
RUN cd /home/root &&\
    curl http://keplerproject.github.io/luarocks/releases/luarocks-2.2.2.tar.gz -L --output luarocks-2.2.2.tar.gz &&\
    tar -xzvf luarocks-2.2.2.tar.gz &&\
    cd luarocks-2.2.2 &&\
    ./configure &&\
    make build &&\
    make install &&\
    rm -Rf /home/root/luarocks-2.2.2 &&\
    rm /home/root/luarocks-2.2.2.tar.gz

# Install the cjson package (unzip is needed for luarocks)
RUN apt-get -y install unzip
RUN cd /home/root &&\
    luarocks install lua-cjson
# Raise the limits for wrk to open many file while load testing
RUN ulimit -c -m -s -t unlimited

# Install go 1.14.4
RUN curl https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz -L --output go1.14.4.linux-amd64.tar.gz 
RUN tar -xvzf go1.14.4.linux-amd64.tar.gz 
RUN rm go1.14.4.linux-amd64.tar.gz
ENV PATH="/go/bin:${PATH}"

# add user 1001
RUN useradd -u 1001 --no-create-home 1001