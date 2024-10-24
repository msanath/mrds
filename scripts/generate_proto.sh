#!/bin/bash

PROTO_DIR="./api"
GEN_GO_DIR="./gen"

# Generate Go code for all .proto files
protoc -I=$PROTO_DIR --go_out=$GEN_GO_DIR --go-grpc_out=$GEN_GO_DIR $PROTO_DIR/*.proto
