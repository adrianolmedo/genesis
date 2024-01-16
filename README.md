# Genesis

REST API of CRM, its architecture is inspired (not based) on Hexagonal Architecture.

## Architecture

![Golang structure inspired on Hexagonal Architecture](https://i.imgur.com/eZxD6S1.png)

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
