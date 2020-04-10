docker build -t docker-go .

docker run -it docker-go

docker run -it -p 8080:8080 docker-go

docker-compose -f docker-compose.yml up -f

docker-compose -f docker-compose.yml up
docker-compose -f docker-compose.yml down