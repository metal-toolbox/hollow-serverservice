# Server Service

> This repository is [Experimental](https://github.com/packethost/standards/blob/master/experimental-statement.md) meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
> This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

The server service is a microservice within the Hollow eco-system. Server service is responsible for providing a store for physical server information. Support to storing the device components that make up the server is available. You are also able to create attributes and versioned-attributes for both servers and the server components.

## Running locally

To run the api server locally you can bring it up with docker-compose.

```
docker compose up
```

If you have never ran `serverservice` before then the server will fail to start the first time. After the DB service is running you need to create the dev database. Run:

```
make dev-database
```
