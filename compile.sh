#!/bin/bash

# Sendlix Go SDK Proto Compilation Script
# This script generates Go code from protocol buffer definitions
# and places them in the internal package to exclude from public documentation

# Create output directory
mkdir -p internal/proto

# Clean old generated files
rm -f internal/proto/*.pb.go

# Generate Go protobuf and gRPC code
protoc \
  --proto_path=proto_files \
  --go_out=internal/proto \
  --go_opt=paths=source_relative \
  --go-grpc_out=internal/proto \
  --go-grpc_opt=paths=source_relative \
  proto_files/*.proto

echo "Proto files compiled successfully and placed in internal package!"