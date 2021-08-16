# Build go executable for linux
GOOS=linux go build
# go mod init hellosvr
# Build Docker Container
docker build -t magdanat/apithundahsvr . 

echo "Cleaning executable..."

# Cleans built Go executable
go clean