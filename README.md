# Hollow

> This repository is [Experimental](https://github.com/packethost/standards/blob/master/experimental-statement.md) meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
> This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

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
