FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o frontend .

FROM alpine:latest AS runner
COPY --from=builder /app/frontend /app/frontend

CMD ["/app/frontend"]