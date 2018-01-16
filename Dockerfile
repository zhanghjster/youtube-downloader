FROM python:3-alpine

LABEL maintaier "benx <zhanghjster@gmail.com>"

# install you-get
RUN apk add --no-cache ffmpeg
RUN pip3 install you-get

RUN mkdir -p /downloader

WORKDIR /downloader
VOLUME /downloader

ADD youtube-downloader /downloader/youtube-downloader

ENTRYPOINT ["/downloader/youtube-downloader"]