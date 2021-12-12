FROM golang:alpine as buildstage
RUN mkdir generation
WORKDIR /generation
COPY . .
RUN go get ./... && go build -o generation cmd/main.go

FROM alpine:latest
WORKDIR /
COPY --from=buildstage generation /

CMD ["./generation"]