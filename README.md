# nyanbot
linebot

# installation

dependency: docker

From project root
```
$ cp ./userdata ~/nyanbot
```

Edit config
```
$ cp ~/nyanbot/config/config.yml.sample ~/nyanbot/config/config.yml
$ vi ~/nyanbot/config/config.yml

$ cp ~/nyanbot/csv/push_message_sample.csv ~/nyanbot/csv/push_message.csv
$ vi ~/nyanbot/csv/push_message.csv
```

Docker build
```
$ docker build -t nyanbot -f ./docker/Dockerfile
```

Docker run
```
$ docker run -v ~/nyanbot:/root/nyanbot nyanbot
```
