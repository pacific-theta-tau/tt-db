# Import official golang image
FROM golang

# Create /app directory in container
RUN mkdir /app

# Move to /app directory
WORKDIR /app

# Copy project into container
COPY . .

# Install dependencies
RUN go mod download

# Build server binary
RUN go build -o main ./cmd/server/main.go

# Expose port 8080 via TCP
EXPOSE 8080

# Recursively test all packages
# CMD ["go", "test", "-v", "./..."]
CMD ["/app/main"]