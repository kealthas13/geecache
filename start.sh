docker network create --driver bridge --subnet 188.168.0.0/24 --gateway 188.168.0.1 mynetwork
docker compose -f "compose.yaml" up -d --build
sleep 30
docker logs geecache1
./sdcs-test.sh 3