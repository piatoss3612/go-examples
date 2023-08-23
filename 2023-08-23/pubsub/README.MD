# pubsub example

## Run RabbitMQ

```sh
$ docker compose up -d
```

## Run subscriber

```sh
$ go run ./sub/
```

## Run publisher

```sh
$ go run ./pub/
```