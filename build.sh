#!/bin/bash
VERSION="0.01"
echo "Building Uptime Monitor Version $VERSION"
echo "Building check as uptime monitor"
echo "Building for Mac - Intel"
export GOOS=darwin
export GOARCH=amd64
go build -o ./releases/ ./... 
mv ./releases/check ./releases/check_intel_mac_$VERSION

echo "Building for Mac - ARM"
export GOOS=darwin
export GOARCH=arm64
go build -o ./releases/ ./... 
mv ./releases/check ./releases/check_arm_mac_$VERSION

echo "Building for Linux - Intel"
export GOOS=linux
export GOARCH=amd64
go build -o ./releases/ ./...
mv ./releases/check ./releases/check_intel_linux_$VERSION

echo "Building for Linux - ARM"
export GOOS=linux
export GOARCH=arm64
go build -o ./releases/ ./...
mv ./releases/check ./releases/check_arm_linux_$VERSION

echo "Building for Windows - Intel"
export GOOS=windows
export GOARCH=amd64
go build -o ./releases/ ./...
mv ./releases/check.exe ./releases/check_intel_win_$VERSION.exe

echo "Building for Windows - ARM"
export GOOS=windows
export GOARCH=arm64
go build -o ./releases/ ./...
mv ./releases/check.exe ./releases/check_arm_win_$VERSION.exe