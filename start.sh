docker compose -f "compose.yaml" up -d --build
sleep 45
./sdcs-test.sh 3