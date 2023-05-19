#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: $0 [32|64] [linux|windows|darwin]"
    exit 1
fi

if [ -z "$2" ]; then
    echo "Usage: $0 [32|64] [linux|windows|darwin]"
    exit 1
fi

if [ "$1" == "32" ]; then
    export GOARCH=386
elif [ "$1" == "64" ]; then
    export GOARCH=amd64
else
    echo "Usage: $0 [32|64] [linux|windows|darwin]"
    exit 1
fi

if [ "$2" == "linux" ]; then
    export GOOS=linux
elif [ "$2" == "windows" ]; then
    export GOOS=windows
elif [ "$2" == "darwin" ]; then
    export GOOS=darwin
else
    echo "Usage: $0 [32|64] [linux|windows|darwin]"
    exit 1
fi

go build -ldflags "-s -w -linkmode internal -extldflags -static" -o bin/$(basename $PWD)_$GOOS-$GOARCH
upx --brute -9 bin/$(basename $PWD)_$GOOS-$GOARCH
