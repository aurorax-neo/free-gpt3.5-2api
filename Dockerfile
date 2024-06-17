# Start with a base image containing Go runtime
FROM golang:1.22.0 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the Go app ensuring that the binary is statically linked
RUN CGO_ENABLED=0 go build -o /app/f99871322b15ff5007f0 .

# Now use a smaller image to run the app
FROM alpine:latest

# Set the working directory in the new container
WORKDIR /app

# Copy the statically-linked binary into the new container
COPY --from=builder /app/f99871322b15ff5007f0 /app/f99871322b15ff5007f0

# This container exposes port 3040 to the outside world
EXPOSE 3040

# Run the binary.
CMD [ "/app/f99871322b15ff5007f0" ]
