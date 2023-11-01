# Metrics ingestion into Kafka and consumption with Golang 

The project is inspired by The Primeagen's article review - [From $erverless to Elixir](https://www.youtube.com/watch?v=UGG2HMonQ1c). The reviewed article is [here](https://medium.com/coryodaniel/from-erverless-to-elixir-48752db4d7bc). Instead of Elixir, I decide to use Golang, and build a simplified version.

## Resources

The list of resource and tools used to build the project:

- [kcat](https://github.com/edenhill/kcat), CLI to query Kafka;
- [wrk](https://github.com/wg/wrk), for load testing;
- [Tutorial used to setup Kafka locally via Docker Compose](https://hackernoon.com/setting-up-kafka-on-docker-for-local-development).

## Development and testing

### CURL with JSON

For convenience, here is a CURL command to send a JSON payload to the server:

For API gateway:

```
curl -X POST -H "Content-Type: application/json" \
    -d '{"type": "red", "client_id": 1}' \
    http://localhost:8080/metrics
```

For consumer:

```
curl -X GET http://localhost:8080/statistics
```
