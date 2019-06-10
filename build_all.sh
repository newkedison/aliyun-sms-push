#!/bin/bash

echo "building sms_linux_x64"
go build
mv sms sms_linux_x64

echo "building sms_x86.exe"
GOOS=windows GOARCH=386 go build
mv sms.exe sms_x86.exe

echo "building sms_x64.exe"
GOOS=windows GOARCH=amd64 go build
mv sms.exe sms_x64.exe

echo "packing all"
rm -rf sms_push
mkdir -p sms_push
mv sms_linux_x64 *.exe sms_push
cp config.yaml.template sms_push/config.yaml
cp ip.txt phone.txt sms_push/
zip -r sms_push.zip sms_push
rm -rf sms_push
