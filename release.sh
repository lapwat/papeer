#!/usr/bin/env bash

version=$1
platforms=("linux/amd64" "darwin/amd64" "windows/amd64")

if [ "$#" -ne 1 ]; then
    echo "Illegal number of parameters"
    echo "Usage: ./release.sh X.X.X"
    exit 1
fi

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='papeer-v'$version'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
