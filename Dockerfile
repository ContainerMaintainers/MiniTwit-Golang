# Base this docker container off the official golang docker image.
# Docker containers inherit everything from their base.
FROM golang:1.19-alpine

# Create a directory inside the container to store all our application 
# and then make it the working directory.
WORKDIR /app

# Copy and Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything from the repo (where the Dockerfile lives) into the container.
COPY . .

# Set the PORT environment variable inside the container
ENV PORT=8080

# Build the application (this will build only minitwit.go)
RUN go build -o . ./src/minitwit.go

# Expose port 8080 to the host so we can access our application
EXPOSE $PORT

# Now tell Docker what command to run when the container starts
CMD ["run"]