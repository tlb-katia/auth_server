FROM golang:1.20.3-alpine AS builder

COPY . /github.com/tlb_katia/auth/source/
WORKDIR /github.com/tlb_katia/auth/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/tlb_katia/auth/source/bin/crud_server .

CMD ["./crud_server"]