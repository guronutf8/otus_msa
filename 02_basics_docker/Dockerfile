FROM golang:alpine3.19 as builder
LABEL authors="Pavel.Ok"
WORKDIR .
COPY ./app app
RUN go build -o build/webserver app/main.go

FROM alpine:3.19.1
WORKDIR /home/
COPY --from=builder /go/build/webserver webserver
EXPOSE 8000
CMD ./webserver