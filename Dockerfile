FROM golang:1.12.5 as builder
WORKDIR /hypatia/
COPY . /hypatia/
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o hypatia ./cmd/main.go

FROM alpine:latest
USER root
WORKDIR /home/app

RUN apk add npm
RUN npm install
RUN npm i api2html -g

RUN apk add --no-cache git make musl-dev go

COPY --from=builder hypatia .

EXPOSE 50000

ENTRYPOINT ["./hypatia"]