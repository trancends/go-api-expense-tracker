# use official Golang image
FROM golang:alpine as build

# set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o api .

FROM alpine:latest

# Copy the binary from the build stage
COPY --from=build /api /api

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]
