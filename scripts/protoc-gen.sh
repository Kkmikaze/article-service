#!/bin/bash
protoc -I ./proto \
  --proto_path ./third_party/proto \
  --go_out ./stubs --go_opt paths=source_relative \
  --go-grpc_out ./stubs --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./stubs --grpc-gateway_opt paths=source_relative \
  ./proto/*/v1/*.proto


# Generate Open API V2
protoc -I ./proto \
  --proto_path ./third_party/proto \
  --openapiv2_out ./api \
  --openapiv2_opt=disable_default_responses=true \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt use_go_templates=true \
  ./proto/*/v1/*_service.proto