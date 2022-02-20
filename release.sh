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
    output_name=papeer

    if [ $GOOS = "windows" ]; then
        env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_name.exe"
        zip "$output_name-v$version-$GOOS-$GOARCH.exe.zip" "$output_name.exe"
        rm "$output_name.exe"
    else
        env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_name"
        tar czvf "$output_name-v$version-$GOOS-$GOARCH.tar.gz" "$output_name"
        rm "$output_name"
    fi
done
