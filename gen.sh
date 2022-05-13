#!/bin/bash
set -e

protoc --go_out=paths=source_relative:./gin_gateway gin_gateway_option.proto
