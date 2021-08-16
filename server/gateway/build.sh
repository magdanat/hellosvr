echo "Running build script..."

# Build go executable for linux
# GOOS=linux go build
# go mod init hellosvr

echo "Building Docker Container..."

# Build Docker Container
DOCKER_BUILDKIT=0 docker build -t magdanat/apithundahsvr . 

echo "Cleaning executable..."

# Cleans built Go executable
go clean