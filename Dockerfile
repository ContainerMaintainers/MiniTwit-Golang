# Base this docker container off the official golang docker image.
# Docker containers inherit everything from their base.
FROM alpine:edge AS build
RUN apk add --no-cache --update go gcc g++

ARG db_user
ARG db_host
ARG db_password
ARG db_name
ARG db_port
ARG port
ARG session_key
ARG gin_mode

# Create a directory inside the container to store all our application 
# and then make it the working directory.
WORKDIR /usr/src/app

# Copy everything
COPY . .

# Give permissions to run env_file.sh
RUN chmod +x env_file.sh

# Create .env if it doesn't exist
RUN ./env_file.sh ${db_user} ${db_password} ${db_name} ${db_port} ${db_host} ${port} ${session_key} ${gin_mode}

# Download Go modules
RUN go mod download && go mod tidy

# Build the application (this will build only minitwit.go)
RUN go build -o . minitwit.go

# Expose port 8080 to the host so we can access our application
EXPOSE 8080

# Now tell Docker what command to run when the container starts
# This will run the compiled minitwit file, when it is ready
CMD ["./minitwit"]
