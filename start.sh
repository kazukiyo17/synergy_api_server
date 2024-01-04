docker stop synergy_api_server
docker rm synergy_api_server
docker rmi synergy_api_server
docker build -t synergy_api_server .
#docker pull registry.cn-shanghai.aliyuncs.com/game_xyunli/synergy_api_server:latest
docker run -d --name synergy_api_server -p 8081:8081 --network game-network synergy_api_server