### docker build -t dialogflow-boobs:latest .
### docker run --rm -it dialogflow-boobs:latest

FROM golang:latest as builder

COPY . /go/src/dialogflow-boobs/
WORKDIR /go/src/dialogflow-boobs/
RUN go get
RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s"

FROM alpine:3.5
RUN apk update && apk add ca-certificates
COPY --from=builder /go/src/dialogflow-boobs/dialogflow-boobs /dialogflow-boobs
WORKDIR /
ENTRYPOINT ["/dialogflow-boobs"]
CMD ["--help"]