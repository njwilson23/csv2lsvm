#! /bin/sh
# Compile for multiple architectures

env GOOS=windows GOARCH=386 go build -v -o csv2lsvm_386.exe
env GOOS=windows GOARCH=amd64 go build -v -o csv2lsvm.exe

env GOOS=linux GOARCH=amd64 go build -v -o csv2lsvm
