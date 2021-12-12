FROM golang:alpine as buildstage
RUN mkdir persistence
WORKDIR /persistence
COPY . .
RUN go get ./... && go build -o persistence cmd/main.go

FROM alpine:latest
WORKDIR /
COPY --from=buildstage persistence /

CMD ["./persistence"]