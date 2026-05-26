FROM golang:1.25-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o thecodeisme cmd/thecodeisme.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/thecodeisme .
COPY --from=builder /build/config /config

CMD ["./thecodeisme"]