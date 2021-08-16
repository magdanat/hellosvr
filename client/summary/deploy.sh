sh ./build.sh

docker login
docker push magdanat/thundahclient
ssh ec2-user@ec2-54-183-61-235.us-west-1.compute.amazonaws.com 'docker pull magdanat/thundahclient && docker rm -f client && export TLSCERT=/etc/letsencrypt/live/thundah.com/fullchain.pem && export TLSKEY=/etc/letsencrypt/live/thundah.com/privkey.pem && docker run --name client -d -p 80:80 -p 443:443 -v ~/etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY magdanat/thundahclient'
$SHELL