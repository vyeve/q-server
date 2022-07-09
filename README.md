# q-server
Test server for bulk upload

## Configuration

To configure server was used environment approach.

- *APP_SQLITE_DB_PATH* defines path to database. **REQUIRED**
- *APP_SERVER_PORT* defines port on which server will listen (default is **8080**)
- *APP_LOG_LEVEL* is used to define log level. Possible values are `debug`, `info`, `warn`, `error`, `fatal`. Default is **info**.
- *APP_SQLITE_DB_POOL_SIZE* defines pool size of connections to database.

## Usage

There are some ways to launch server

### 1. Run standalone application

```shell
APP_SQLITE_DB_PATH=assets/schema.sqlite go run main.go
```

### 2. Run docker

```shell
docker build -t q-server .
docker run -it --rm -p 8080:8080 -v $(pwd)/assets:/app/data --env-file config.env q-server
```

### 3. Run docker-compose

```shell
docker-compose up
```

## TODO

This approach will be work with count transfers less than 1000. If there are expecting more transfers in one receipt, should split `INSERT INTO` queries to bulks.
