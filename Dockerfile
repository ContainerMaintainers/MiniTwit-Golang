# Base this docker container off the official golang docker image.
# Docker containers inherit everything from their base.
FROM golang:1.16-alpine

# Create a directory inside the container to store all our application and then make it the working directory.
#RUN mkdir -p /go/src/app
#WORKDIR /go/src/app
WORKDIR /app

# Copy and Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything from the repo (where the Dockerfile lives) into the container.
ADD . /app

# Download and install any required third party dependencies into the container.
RUN go-wrapper download
RUN go-wrapper install

# Set the PORT environment variable inside the container
ENV PORT=8080

# Expose port 8080 to the host so we can access our application
EXPOSE $PORT

# Build the application
RUN go build -o main .

# Now tell Docker what command to run when the container starts
CMD ["go-wrapper", "run"]