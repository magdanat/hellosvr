sh ./build.sh

docker login
docker push magdanat/thundahclient
# ssh ec2-user@ec2-54-183-61-235.us-west-1.compute.amazonaws.com 'docker pull magdanat/thundahclient && docker rm -

$SHELL