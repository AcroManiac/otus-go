#!/bin/bash

cd ../api

# Generate event interface
protoc event.proto --go_out=plugins=grpc:../pkg/api

# Generate Calendar API server
protoc api.proto --go_out=plugins=grpc:../pkg/api
