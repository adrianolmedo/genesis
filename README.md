# Genesis

REST API of CRM, architecture based on the main objective of Hexagonal Architecture: to isolate business logic from external dependencies.

## Architecture

![Golang structure based on Hexagonal Architecture purpose](https://i.imgur.com/eZxD6S1.png)

## Run with Postgres service

```bash
$ git clone https://github.com/adrianolmedo/genesis.git
$ cp .env.example .env
$ docker-compose up -d --build postgres
```

**Run migrations**

```bash
./migrate.sh
```

**Up application service:**

```bash
$ make genrsa
$ make swagger
$ docker-compose up -d --build app
```
