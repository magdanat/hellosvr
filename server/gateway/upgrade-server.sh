# echo 
# docker pull magdanat/apithundahsvr
# docker rm -f 443gateway
# export TLSCERT=/etc/letsencrypt/live/api.thundah.com/fullchain.pem
# export TLSKEY=/etc/letsencrypt/live/api.thundah.com/privkey.pem
# docker run -d --name 443gateway -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY magdanat/apithundahsvr

# echo We have reached the end