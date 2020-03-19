#!/bin/bash

echo "Stopping microservices..."

# Calendar API gRPC server
unset arr
arr=($(ps aux | grep calendar_api | awk 'BEGIN {FS=" "}{print $2}'))
kill -2 "${arr[@]}"

# Calendar scheduler
unset arr
arr=($(ps aux | grep calendar_scheduler | awk 'BEGIN {FS=" "}{print $2}'))
kill -2 "${arr[@]}"

# Calendar sender
unset arr
arr=($(ps aux | grep calendar_sender | awk 'BEGIN {FS=" "}{print $2}'))
kill -2 "${arr[@]}"

echo "Microservices ware stopped successfully. Bye!"
