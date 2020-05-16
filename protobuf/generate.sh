#!/bin/bash
# command to save file format correctly in vim
# :set fileformat=unix 
# command to set file exectuate permissions
# sudo chmod u+x run.sh

# generate go protobuf
protoc -I . --go_out=plugins=grpc:. ./caterpillar.proto

# generate python protobuf
# fix output file by changing import in caterpillar_pb2_grpc.py to from . import caterpillar_pb2 as caterpillar__pb2
python -m grpc_tools.protoc -I./ --python_out=../internal_py/. --grpc_python_out=../internal_py/. ./caterpillar.proto
