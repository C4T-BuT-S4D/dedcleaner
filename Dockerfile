FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/dedcleaner/
COPY main.go .
RUN go build -o /go/bin/dedcleaner

FROM scratch

COPY --from=builder /go/bin/dedcleaner /go/bin/dedcleaner
CMD ["/go/bin/dedcleaner"]
