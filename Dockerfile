# Base Image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /GPT-Bot

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download

RUN go get

# Build the Go app
RUN go build

# Command to run the executable
CMD ["./GPT-Bot"]