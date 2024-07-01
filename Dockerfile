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

# RUN go build -o main cmd/server/main.go
RUN go build -o main main.go

# Expose port 8080 via TCP
EXPOSE 8080

# Run app using local postgres container
CMD ["/app/main", "-env=dev"]

# Run app using prod DB
# CMD ["/app/main", "-env=prod"]
