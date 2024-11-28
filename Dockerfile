# Step 1: Use official Go image as a base image
FROM golang:1.22.3-alpine as builder

# Step 2: Set the Current Working Directory inside the container
WORKDIR /usr/src/app

# Step 3: Copy go.mod and go.sum and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the source code into the container
COPY . .

# Step 5: Build the Go app
RUN go build -o main .

# Step 6: Start a new stage for the final image (this reduces image size)
FROM alpine:latest

# Step 7: Install PostgreSQL client (optional, for troubleshooting and testing)
RUN apk --no-cache add postgresql-client

# Step 8: Set the Current Working Directory inside the container
WORKDIR /root/

# Step 9: Copy the binary file from the build stage
COPY --from=builder /usr/src/app/main .

# Step 10: Expose the port the app will run on
EXPOSE 8080

# Step 11: Command to run the executable
CMD ["./main"]
