#!/bin/bash

set -euxo pipefail

WORKING_DIR="$(pwd)"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Helper packages
# TODO: Currently no tests

# Implementations

# Modules that don't require a service
array=( badgerdb bbolt bigcache file freecache gomap leveldb syncmap )
for MODULE_NAME in "${array[@]}"; do
    echo "testing $MODULE_NAME"
    (cd "$SCRIPT_DIR"/../"$MODULE_NAME" && go test -v -race -coverprofile=coverage.txt -covermode=atomic) || (cd "$WORKING_DIR" && echo " failed" && exit 1)
done

# Modules that require a service
# CockroachDB
docker run -d --rm --name cockroachdb -p 26257:26257 cockroachdb/cockroach start-single-node --insecure
sleep 10s
docker exec cockroachdb bash -c './cockroach sql --insecure --execute="create database gokv;"'
(cd "$SCRIPT_DIR"/../cockroachdb && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop cockroachdb) || (cd "$WORKING_DIR" && echo " failed" && docker stop cockroachdb && exit 1)
# Consul
docker run -d --rm --name consul -e CONSUL_LOCAL_CONFIG='{"limits":{"http_max_conns_per_client":1000}}' -p 8500:8500 bitnami/consul
sleep 10s
(cd "$SCRIPT_DIR"/../consul && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop consul) || (cd "$WORKING_DIR" && echo " failed" && docker stop consul && exit 1)
# Google Cloud Datastore via "Cloud Datastore Emulator"
# Using the ":slim" or ":alpine" tag would require the emulator to be installed manually.
# Both ways seem to be okay for setting the project: `-e CLOUDSDK_CORE_PROJECT=gokv` and CLI parameter `--project=gokv`
# `--host-port` is required because otherwise the server only listens on localhost IN the container.
docker run -d --rm --name datastore -p 8081:8081 google/cloud-sdk gcloud beta emulators datastore start --no-store-on-disk --project=gokv --host-port=0.0.0.0:8081
sleep 10s
(cd "$SCRIPT_DIR"/../datastore && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop datastore) || (cd "$WORKING_DIR" && echo " failed" && docker stop datastore && exit 1)
# DynamoDB via "DynamoDB local"
docker run -d --rm --name dynamodb-local -p 8000:8000 amazon/dynamodb-local
sleep 10s
(cd "$SCRIPT_DIR"/../dynamodb && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop dynamodb-local) || (cd "$WORKING_DIR" && echo " failed" && docker stop dynamodb-local && exit 1)
# etcd
docker run -d --rm --name etcd -p 2379:2379 --env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd
sleep 10s
(cd "$SCRIPT_DIR"/../etcd && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop etcd) || (cd "$WORKING_DIR" && echo " failed" && docker stop etcd && exit 1)
# Hazelcast
docker run -d --rm --name hazelcast -p 5701:5701 hazelcast/hazelcast
sleep 10s
(cd "$SCRIPT_DIR"/../hazelcast && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop hazelcast) || (cd "$WORKING_DIR" && echo " failed" && docker stop hazelcast && exit 1)
# Apache Ignite
docker run -d --rm --name ignite -p 10800:10800 apacheignite/ignite
sleep 10s
(cd "$SCRIPT_DIR"/../ignite && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop ignite) || (cd "$WORKING_DIR" && echo " failed" && docker stop ignite && exit 1)
# Memcached
docker run -d --rm --name memcached -p 11211:11211 memcached
sleep 10s
(cd "$SCRIPT_DIR"/../memcached && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop memcached) || (cd "$WORKING_DIR" && echo " failed" && docker stop memcached && exit 1)
# MongoDB
docker run -d --rm --name mongodb -p 27017:27017 mongo
sleep 10s
(cd "$SCRIPT_DIR"/../mongodb && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop mongodb) || (cd "$WORKING_DIR" && echo " failed" && docker stop mongodb && exit 1)
# MySQL
docker run -d --rm --name mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true -p 3306:3306 mysql
sleep 10s
(cd "$SCRIPT_DIR"/../mysql && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop mysql) || (cd "$WORKING_DIR" && echo " failed" && docker stop mysql && exit 1)
# PostgreSQL
docker run -d --rm --name postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=gokv -p 5432:5432 postgres:alpine
sleep 10s
(cd "$SCRIPT_DIR"/../postgresql && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop postgres) || (cd "$WORKING_DIR" && echo " failed" && docker stop postgres && exit 1)
# Redis
docker run -d --rm --name redis -p 6379:6379 redis
sleep 10s
(cd "$SCRIPT_DIR"/../redis && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop redis) || (cd "$WORKING_DIR" && echo " failed" && docker stop redis && exit 1)
# Amazon S3 via Minio
docker run -d --rm --name s3 -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE" -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" -p 9000:9000 minio/minio server /data
sleep 10s
(cd "$SCRIPT_DIR"/../s3 && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop s3) || (cd "$WORKING_DIR" && echo " failed" && docker stop s3 && exit 1)
# Tablestorage via Azurite
# In the past there was this problem: https://github.com/Azure/Azurite/issues/121
# With this Docker image:
#docker run -d --rm --name azurite -e executable=table -p 10002:10002 arafato/azurite
# Now with the official image it still doesn't work. // TODO: Investigate / create GitHub issue.
#docker run -d --rm --name azurite -p 10002:10002 mcr.microsoft.com/azure-storage/azurite azurite-table
#sleep 10s
#(cd "$SCRIPT_DIR"/../tablestorage && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop azurite) || (cd "$WORKING_DIR" && echo " failed" && docker stop azurite && exit 1)
#
# Alibaba Cloud Table Store
# TODO: Currently no emulator exists for Alibaba Cloud Table Store.
#
# Apache ZooKeeper
docker run -d --rm --name zookeeper -p 2181:2181 zookeeper
sleep 10s
(cd "$SCRIPT_DIR"/../zookeeper && go test -v -race -coverprofile=coverage.txt -covermode=atomic && docker stop zookeeper) || (cd "$WORKING_DIR" && echo " failed" && docker stop zookeeper && exit 1)

# Examples
# TODO: Currently no tests

cd "$WORKING_DIR"
