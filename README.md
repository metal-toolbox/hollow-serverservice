[![slack](https://img.shields.io/badge/slack-delivery--eng-brightgreen.svg?style=for-the-badge&logo=slack)](https://app.slack.com/client/T02BWUJ78/C013J3XQHM4) [![pagerduty](https://img.shields.io/badge/pagerduty-oncall-red.svg?style=for-the-badge&logo=pagerduty)](https://packet.pagerduty.com/schedules#P09YKZW)

# Hollow

North of Neverland, `Hollow` is the magical kingdom where fairies are born and is their home.

More pragmatically, Hollow is the combination of a golang-based API built on top of cockroachdb primarily responsible for the handling of physical asset information (e.g. metadata about our physical assets)

## Running locally

To run the api server locally you can bring it up with docker-compose.

```
docker compose up
```

If you have never ran hollow before then the server will fail to start the first time. After the DB service is running you need to create the dev database. Run:

```
make local-dev-databases
```
