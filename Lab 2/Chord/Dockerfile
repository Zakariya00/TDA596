# Create build stage based on buster image
FROM golang:1.16-buster AS builder

# Create working directory under /app
WORKDIR /app

# Copy over all go config (go.mod, go.sum etc.)
COPY go.* ./

# Install any required modules
RUN go mod download

# Copy over Go source code
COPY *.go ./

# Run the Go build and output binary under chord
RUN go build -o /Chord

# Make sure to expose the port the Chord Node is using
EXPOSE 8080
EXPOSE 1234
EXPOSE 9090
EXPOSE 1235
EXPOSE 5050

# Run the app binary when we run the container
ENTRYPOINT ["/Chord"]

