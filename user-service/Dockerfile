FROM golang:1.24.5-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port
EXPOSE 8889

# Run the application
CMD ["./main"]