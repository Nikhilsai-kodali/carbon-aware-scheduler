FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY scripts/fetch_carbon.go ./
RUN go build -o carbon-checker fetch_carbon.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/carbon-checker .
CMD ["./carbon-checker"]