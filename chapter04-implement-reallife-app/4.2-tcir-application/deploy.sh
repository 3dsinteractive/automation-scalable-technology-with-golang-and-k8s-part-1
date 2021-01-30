#!/bin/bash

# 1. Format of docker image is $DOCKER_REPOSITORY/$IMAGE:$DEPLOY_ENV-$APP_VERSION.$TIMESTAMP
#    for example 3dsinteractive/automation-technology:prd-1.0.20210118001006
APP_VERSION=1.0
TIMESTAMP=20210118001006
DEPLOY_ENV=prd
DOCKER_REPOSITORY=3dsinteractive

# 2. commit will push docker image to repository
function commit() {
    local IMAGE=$1
    echo "docker push image : $DOCKER_REPOSITORY/$IMAGE:$DEPLOY_ENV-$APP_VERSION.$TIMESTAMP"
    docker push $DOCKER_REPOSITORY/$IMAGE:$DEPLOY_ENV-$APP_VERSION.$TIMESTAMP
}

# 3. build_api is the main function to build Dockerfile
function build_api() {
    local IMAGE=automation-technology

    # If found go in default path, it will use go from default path
    GO=/usr/local/go/bin/go
    if [ -f "$GO" ]; then
        /usr/local/go/bin/go mod init automationworkshop/main
        /usr/local/go/bin/go get
        /usr/local/go/bin/go mod vendor
    else 
        go mod init automationworkshop/main
        go get
        go mod vendor
    fi

    # Build the Dockerfile
    docker build -f Dockerfile -t $DOCKER_REPOSITORY/$IMAGE:$DEPLOY_ENV-$APP_VERSION.$TIMESTAMP .
    commit $IMAGE
}

# 4. Validate APP_VERSION must not empty
if [ "$APP_VERSION" = "" ]; then
    echo -e "APP_VERSION cannot be blank"
    exit 1
fi

# 5. Validate TIMESTAMP must not empty
if [ "$TIMESTAMP" = "" ]; then
    echo -e "TIMESTAMP cannot be blank"
    exit 1
fi

# 6. Validate DEPLOY_ENV must not empty
if [ "$DEPLOY_ENV" == "" ]; then
    echo -e "DEPLOY_ENV cannot be blank"
    exit 1
fi

# 7. Validate DOCKER_REPOSITORY must not empty
if [ "$DOCKER_REPOSITORY" == "" ]; then
    echo -e "DOCKER_REPOSITORY cannot be blank"
    exit 1
fi

# 8. Run main build process
build_api