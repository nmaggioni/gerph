#!/bin/sh

echo -e "\n* Cleaning up"
rm -rf ./dist/
echo "* Creating directory"
mkdir dist &>/dev/null
echo "* Preparing asset bundle"
$GOPATH/bin/rice embed-go
echo -e "\n* Compiling for Windows x86_64"
GOOS=windows GOARCH=amd64 go build -o dist/gerph_windows_amd64.exe
echo "* Compiling for Windows x86"
GOOS=windows GOARCH=386 go build -o dist/gerph_windows_i386.exe
echo "* Compiling for Linux x86_64"
GOOS=linux GOARCH=amd64 go build -o dist/gerph_linux_amd64.bin
echo "* Compiling for Linux x86"
GOOS=linux GOARCH=386 go build -o dist/gerph_linux_i386.bin
echo "* Compiling for Mac OS x86_64"
GOOS=darwin GOARCH=amd64 go build -o dist/gerph_darwin_amd64
echo "* Compiling for ARMv6"
GOARCH=arm GOARM=6 go build -o dist/gerph_armv6.bin
echo "* Compiling for ARMv5"
GOARCH=arm GOARM=5 go build -o dist/gerph_armv5.bin
echo -e "\n* Cleaning up"
rm rice-box.go

echo -ne "\n"
file dist/*
