#!/bin/bash
VERSION="0.02"
echo "Building Uptime Monitor Version $VERSION"
echo "Building check as uptime monitor"
echo "Building for Mac - Intel"
GOOS=darwin GOARCH=amd64 go build -o  ./releases/check_intel_mac_$VERSION main.go

echo "Building for Mac - ARM"
GOOS=darwin GOARCH=arm64 go build -o  ./releases/check_arm_mac_$VERSION main.go

echo "Building for Linux - Intel"
GOOS=linux GOARCH=amd64 go build -o ./releases/check_intel_linux_$VERSION main.go

echo "Building for Linux - ARM"
GOOS=linux GOARCH=arm64 go build -o ./releases/check_arm_linux_$VERSION main.go

echo "Building for Windows - Intel"
GOOS=windows GOARCH=amd64 go build -o ./releases/check_intel_win_$VERSION.exe main.go

echo "Building for Windows - ARM"
GOOS=windows GOARCH=arm64 go build -o ./releases/check_arm_win_$VERSION.exe main.go