FROM golang:alpine as buildstage
RUN mkdir rest-api
WORKDIR /rest-api
COPY . .
RUN go get ./... && go build -o rest-api cmd/main.go

FROM alpine:latest
WORKDIR /
COPY --from=buildstage rest-api /

CMD ["./rest-api"]