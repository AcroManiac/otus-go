#!/bin/bash

BINDIR=../bin
LOGDIR=../logs
CONFIGDIR=../configs

mkdir $LOGDIR
cd $LOGDIR

echo "Starting microservices..."
$BINDIR/calendar_api --config=$CONFIGDIR/calendar_api.yaml &> /dev/null &
$BINDIR/calendar_scheduler --config=$CONFIGDIR/calendar_scheduler.yaml &> /dev/null &
$BINDIR/calendar_sender --config=$CONFIGDIR/calendar_sender.yaml &> /dev/null &

sleep 2

echo "Creating and sending data for processing"
$BINDIR/calendar_api_client --config=$CONFIGDIR/calendar_api_client.yaml &> /dev/null &

echo "Data was sent. Watch logs for results"
