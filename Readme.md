#### Context

A simple redis client (consumer) go application connected to redis HA setup.

The application is dockerized and the container is available at `docker.io/niteshsince1982/redis-poc:0.0.7`

Redis HA setup is deployed using Redis operator `https://github.com/spotahome/redis-operator`

#### APIs

* To store data to Redis

`http://<externalIP>/store?key=someKey&value=someValue`

* To read data from Redis

`http://<externalIP>/get?key=someKey`

#### Running redis server locally on docker

docker run -d --name some-redis -p 6379:6379 redis

docker exec -it some-redis bash

root@72c388dc2cb8:/data# redis-cli

Now from the redis application the host:port to store the content in redis will be 

`localhost:6379`

