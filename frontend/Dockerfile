FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o frontend .

FROM alpine:latest AS runner
COPY --from=builder /app/frontend /app/frontend
COPY --from=builder /app/views /views

# Verify that the views folder contains HTML files
RUN if [ -z "$(ls -A /views/*.html 2>/dev/null)" ]; then echo "No HTML files found in /app/views"; exit 1; fi

CMD ["/app/frontend"]