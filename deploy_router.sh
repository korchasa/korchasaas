#!/usr/bin/env bash
set -ex
docker-compose -f ./docker-compose.yml -f ./docker-compose.dev.yml build router
docker tag korchasaas_router korchasa/korchasaas:latest
docker push korchasa/korchasaas:latest
IP=$(hyper fip detach korchasaas)
! hyper rm -f korchasaas
hyper pull korchasa/korchasaas:latest
hyper run -d -p 80 --size=s1 --name korchasaas korchasa/korchasaas
hyper fip attach $IP korchasaas
! hyper rmi $(hyper images | grep "^<none>" | awk '{print $3}')
hyper logs -f korchasaas
