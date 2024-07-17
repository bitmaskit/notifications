FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o slack-service .

FROM alpine:latest AS runner
COPY --from=builder /app/slack-service /app/slack-service

CMD ["/app/slack-service"]