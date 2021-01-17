#!/bin/bash

CMD=build_api
APP_VERSION=1.0
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
NAMESPACE=prd
COMPANY=3dsinteractive

function print_howto() {
    echo "How to use : deploy.sh CMD APP_VERSION TIMESTAMP NAMESPACE OTHER_PARAM1 OTHER_PARAM2 (eg. deploy.sh build_api 1.0 201906011030 prd)"
}

function commit() {
    local IMAGE=$1
    echo "docker push image : $COMPANY/$IMAGE:$NAMESPACE-$APP_VERSION.$TIMESTAMP"
    # docker push $COMPANY/$IMAGE:$NAMESPACE-$APP_VERSION.$TIMESTAMP
}

function build_api() {
    local IMAGE=automation-technology

    GO=/usr/local/go/bin/go
    if [ -f "$GO" ]; then
        /usr/local/go/bin/go get
        /usr/local/go/bin/go mod vendor
    else 
        go get
        go mod vendor
    fi

    docker build -f Dockerfile -t $COMPANY/$IMAGE:$NAMESPACE-$APP_VERSION.$TIMESTAMP .
    commit $IMAGE
}

if [ "$CMD" = "" ]; then
    echo -e "CMD cannot be blank"
    print_howto
    exit 1
fi

if [ "$APP_VERSION" = "" ]; then
    echo -e "APP_VERSION cannot be blank"
    print_howto
    exit 1
fi

if [ "$TIMESTAMP" = "" ]; then
    echo -e "TIMESTAMP cannot be blank"
    print_howto
    exit 1
fi

if [ "$NAMESPACE" == "" ]; then
    echo -e "NAMESPACE cannot be blank"
    print_howto
    exit 1
fi

eval ${CMD}