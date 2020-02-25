#!/bin/bash

cd api

# Generate event interface
protoc event.proto --go_out=plugins=grpc:../internal/domain/grpc

# Generate Calendar API server
protoc api.proto --go_out=plugins=grpc:../internal/domain/grpc
