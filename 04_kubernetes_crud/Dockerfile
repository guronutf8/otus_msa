FROM golang:alpine3.19 as builder
LABEL authors="Pavel.Ok"
WORKDIR .
COPY ./app app
WORKDIR ./app
RUN go build -o build/webserver main.go

FROM alpine:3.19.1
WORKDIR /home/
COPY --from=builder /go/app/build/webserver webserver
EXPOSE 80
CMD ./webserver