echo "Running deployment script..."

# Calls Build Script
sh ./build.sh

echo "Finished running build script without errors."

echo "Logging into Docker..."
# Push First
docker login

echo "Pushing current build..."

docker push magdanat/apithundahsvr

echo "SSH and running upgrading server..."

# Run Upgrade-Server script
ssh ec2-user@ec2-54-219-192-75.us-west-1.compute.amazonaws.com 'docker pull magdanat/apithundahsvr && docker rm -f 443gateway && export TLSCERT=/etc/letsencrypt/live/api.thundah.com/fullchain.pem && export TLSKEY=/etc/letsencrypt/live/api.thundah.com/privkey.pem && docker run -d --name 443gateway -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY magdanat/apithundahsvr'
# 'bash -s' < upgrade-server.sh

$SHELL