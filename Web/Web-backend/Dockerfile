FROM golang:1.17-alpine

#为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /project/docker_first/

#COPY go.mod , go.sum and download the dependencied
COPY go.* ./
RUN go mod download

#COPY All things inside the project and build
COPY . .
RUN go build -o /project/docker_first/build/myapp .

EXPOSE 8888
ENTRYPOINT ["/project/docker_first/build/myapp"]