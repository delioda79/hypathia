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
RUN npm config set unsafe-perm true && \
npm i api2html -g && \
npm config set unsafe-perm false

RUN apk add --no-cache git make musl-dev go

COPY --from=builder hypatia .

EXPOSE 50000

ENTRYPOINT ["./hypatia"]