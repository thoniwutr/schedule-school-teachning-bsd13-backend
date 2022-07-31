FROM golang:latest AS builder
WORKDIR /go/src/github.com/Beam-Data-Company/merchant-config-svc/
COPY . .
RUN mkdir output
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o output/merchant-config-svc ./cmd/merchant-config-svc

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/Beam-Data-Company/merchant-config-svc/output/merchant-config-svc .
