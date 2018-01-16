#!/usr/bin/env bash

for GOOS in darwin linux; do
    for GOARCH in 386 amd64; do
        docker run -it -v $PWD:/usr/src/youtube-downloader  \
            -w /usr/src/youtube-downloader \
            -e GOPATH=/gopath \
            -v ${GOPATH}:/gopath \
            -e GOOS=${GOOS} \
            -e GOARCH=${GOARCH} \
            golang \
            go build -o youtube-downloader

        tar zcvf youtube-downloader-${GOOS}-${GOARCH}.tar.gz youtube-downloader

        rm -f youtube-downloader

    done
done