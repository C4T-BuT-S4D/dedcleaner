FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git upx

WORKDIR /app
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download
COPY . /app
RUN go build -o /app/dedcleaner
RUN upx --best --ultra-brute /app/dedcleaner

FROM scratch

COPY --from=builder /app/dedcleaner /dedcleaner
CMD ["/dedcleaner"]
