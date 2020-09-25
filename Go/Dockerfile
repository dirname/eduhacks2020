FROM golang:1.15 AS build-dist
ENV GOPROXY='https://mirrors.aliyun.com/goproxy'
WORKDIR /data/release
COPY . .
RUN go build

FROM centos:latest as prod
WORKDIR /data/Go
COPY --from=build-dist /data/release/Go ./
COPY --from=build-dist /data/release/conf /data/Go/conf

EXPOSE 555

CMD ["/data/Go/Go","-c","./conf/app.ini"]