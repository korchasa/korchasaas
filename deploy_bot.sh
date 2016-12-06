#!/usr/bin/env bash
set -ex

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <telegram_token> <api_host>"
    exit 1
fi

docker-compose  -f ./docker-compose.yml -f ./docker-compose.dev.yml build bot
docker tag korchasaas_bot korchasa/korchasaas-bot:latest
docker push korchasa/korchasaas-bot:latest

! hyper rm -f korchasaas-bot
hyper pull korchasa/korchasaas-bot:latest
hyper run -e "TELEGRAM_TOKEN=$1" -e "API_HOST=$2" -d -p 80 --size=s1 --name korchasaas-bot korchasa/korchasaas-bot
! hyper rmi $(hyper images | grep "^<none>" | awk '{print $3}')

hyper logs -f korchasaas-bot
