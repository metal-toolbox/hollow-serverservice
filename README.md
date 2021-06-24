# hollow


## Running locally

To run the api server locally you can bring it up with docker-compose.

```
docker compose up
```

If you have never ran hollow before then the server will fail to start the first time. After the DB service is running you need to create the dev database. Run:

```
make local-dev-databases
```
