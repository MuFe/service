FROM registry.cn-guangzhou.aliyuncs.com/mufe/golang:1.17.8 AS builder

WORKDIR $GOPATH/src/mufe_service

ENV  GO111MODULE=off

COPY ./jsonRpc ./jsonRpc
COPY ./vendor ./vendor
COPY ./model ./model
COPY ./camp ./camp
COPY ./api ./api
COPY ./manager ./manager
COPY ./service ./service

COPY ./adminapifactory/good/src ./adminapifactory/good/src
RUN CGO_ENABLED=0 GOOS=linux go build -o /adminapifactory/good/server -work -v ./adminapifactory/good/src
FROM registry.cn-guangzhou.aliyuncs.com/mufe/alpine:ca as prod
WORKDIR /root/
COPY --from=0 go/src/mufe_service/camp/cert .
COPY --from=0 /adminapifactory/good/ .
CMD ["./server"]
