
# Run in Linux for Go Web Servers
# GOOS=linux go build

# Container image that we're building upon
FROM alpine

# Causes Docker daemon to run command that follows
# inside the new container as it builds the iamge
RUN apk add --no-cache ca-certificates

# Copy current host directory 
# into Docker's 
COPY . /gateway

# Expose the port 443 
EXPOSE 443

# Tells Docker  what program to execute 
ENTRYPOINT ["/gateway"]