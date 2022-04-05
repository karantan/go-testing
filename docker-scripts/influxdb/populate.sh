#!/bin/sh

set -e

influx -execute "CREATE DATABASE $DOCKER_INFLUXDB_INIT_BUCKET"
influx -execute "CREATE RETENTION POLICY day ON $DOCKER_INFLUXDB_INIT_BUCKET DURATION 24h REPLICATION 1"
influx -execute "INSERT INTO "$DOCKER_INFLUXDB_INIT_BUCKET"."day" haproxy,proxy=apache,host=server1 rtime=200"
influx -execute "INSERT INTO "$DOCKER_INFLUXDB_INIT_BUCKET"."day" haproxy,proxy=apache,host=server2 rtime=5001"
influx -execute "INSERT INTO "$DOCKER_INFLUXDB_INIT_BUCKET"."day" haproxy,proxy=apache,host=server3 rtime=300"
