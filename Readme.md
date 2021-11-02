#### Context

A simple redis client (consumer) go application connected to redis HA setup.

The application is dockerized and the container is available at `docker.io/niteshsince1982/redis-poc:0.0.7`

Redis HA setup is deployed using Redis operator `https://github.com/spotahome/redis-operator`

#### APIs

* To store data to Redis

`http://<externalIP>:8000/store?key=someKey&value=someValue`

* To read data from Redis

`http://<externalIP>:8000/get?key=someKey`

