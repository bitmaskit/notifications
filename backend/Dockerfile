FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o backend .

FROM alpine:latest AS runner
COPY --from=builder /app/backend /app/backend

CMD ["/app/backend"]