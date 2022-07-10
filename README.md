# q-server
Test server for upload bulk data

## Configuration

To configure server was used environment approach.

- *APP_SQLITE_DB_PATH* defines path to database. **REQUIRED**
- *APP_SERVER_PORT* defines port on which server will listen. Default is **8080**.
- *APP_LOG_LEVEL* is used to define log level. Possible values are `debug`, `info`, `warn`, `error`, `fatal`.Default is **info**.
- *APP_SQLITE_DB_POOL_SIZE* defines pool size of connections to database.
- *APP_REQUESTS_LIMIT* defines request limit to server. Default is **100**.

## Start application

There are some ways to launch server

### 1. Run standalone application

```shell
APP_SQLITE_DB_PATH=assets/schema.sqlite go run main.go
```

### 2. Run docker

```shell
docker build -t q-server .
docker run -it --rm -p 8080:8080 -v $(pwd)/assets:/app/data --env  APP_SQLITE_DB_PATH=/app/data/schema.sqlite q-server
```

### 3. Run docker-compose

```shell
docker-compose up
```
## Usage

There are two endpoints to upload credit transfers:
- */transfer*
```shell
curl --request POST 'http://127.0.0.1:8080/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "organization_name": "ACME Corp",
    "organization_bic": "OIVUSCLQXXX",
    "organization_iban": "FR10474608000002006107XXXXX",
    "credit_transfers": [
        {
            "amount": "23.17",
            "currency": "EUR",
            "counterparty_name": "Bip Bip",
            "counterparty_bic": "CRLYFRPPTOU",
            "counterparty_iban": "EE383680981021245685",
            "description": "Neverland/6318"
        },
        {
            "amount": "200",
            "currency": "EUR",
            "counterparty_name": "Daffy Duck",
            "counterparty_bic": "DDFCNLAM",
            "counterparty_iban": "NL24ABNA5055036109",
            "description": "2020/RabbitSeason/"
        }
    ]
}'
```
- */upload*
```shell
curl --request POST 'http://127.0.0.1:8080/upload' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/PATH/TO/transfers.json"'
```

## TODO

This approach will be work fine with count transfers less than 1000. If there are expecting more transfers in one request, should split `INSERT INTO` queries to bulks.

Store amounts in cents is not good approach when we need to get some statistic (average, percentage etc). It's better to store cents which multiplied to 1000000 (common practice).
