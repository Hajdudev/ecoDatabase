FROM golang:1.24 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Use a minimal image for the final container
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
