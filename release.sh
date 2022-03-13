#!/usr/bin/env bash

if [ "$#" -ne 1 ]; then
    echo "Illegal number of parameters"
    echo "Usage: ./release.sh X.X.X"
    exit 1
fi

version=$1
platforms=("linux/amd64" "darwin/amd64" "windows/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    if [ $GOOS = "windows" ]; then
        env GOOS=$GOOS GOARCH=$GOARCH go build -o papeer.exe
        zip "papeer-v$version-$GOOS-$GOARCH.exe.zip" papeer.exe
        rm papeer.exe
    else
        env GOOS=$GOOS GOARCH=$GOARCH go build -o papeer
        tar czvf "papeer-v$version-$GOOS-$GOARCH.tar.gz" papeer
        rm papeer
    fi
done
