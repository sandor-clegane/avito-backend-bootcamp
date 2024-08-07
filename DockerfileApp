# Specifies a parent image
FROM golang:1.21.6-bullseye
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN go mod download
 
# Tells Docker which network port your container listens on
EXPOSE 8082
 
# Specifies the executable command that runs when the container starts
CMD ["go", "run", "cmd/house-service/main.go", "-config-path=config/local.yaml"] 
