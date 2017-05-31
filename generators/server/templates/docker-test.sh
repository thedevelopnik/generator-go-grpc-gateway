#! /usr/bin/env bash
export DISABLE_CONSOLE_LOG=true

IMAGE_NAME="$CIRCLE_PROJECT_REPONAME"
 
IMAGE_NAME="$(tr [A-Z] [a-z] <<< "$IMAGE_NAME")"

CONTAINER_IP=127.0.0.1
if [[ "$1" != "" ]]; then
    CONTAINER_IP=$1
fi

if [[ "$CIRCLE_BUILD_NUM" == "" ]]; then
    CIRCLE_BUILD_NUM="dev"
fi

set -ex

echo "Setting version file"
v=$(<version)
echo "$v.$CIRCLE_BUILD_NUM" > version
echo $(<version)

echo "building docker container"
docker build --no-cache -t repo.fanaticslabs.com/$IMAGE_NAME:$CIRCLE_BUILD_NUM .

echo $IMAGE_NAME
echo $CIRCLE_BUILD_NUM

docker run -d \
--name $IMAGE_NAME \
-p 443:443 \
-e "VERSION=$(<version)"
repo.fanaticslabs.com/$IMAGE_NAME:$CIRCLE_BUILD_NUM

set +ex
# wait for container to start
echo 'test route ''http://'$CONTAINER_IP':80/'$HEALTH_CHECK_LOCATION
echo 'awaiting api'
MAX_SECONDS=30
SECONDS=0
until curl --fail https://$CONTAINER_IP/$HEALTH_CHECK_LOCATION; do
    if [[ $SECONDS -gt $MAX_SECONDS ]]; then
        echo "api failed to start in $MAX_SECONDS seconds"
        docker logs $IMAGE_NAME
        exit 1
    fi
    printf '.'
    sleep 1
done
echo ''
set -e

# TODO: swap to actually testing docker, rather than just checking it starts
# TEST_USING_DOCKER=http://$CONTAINER_IP:8088 ./node_modules/.bin/gulp test
