## Content

* [Run with Postgres service](#run-with-postgres-service)

## Run with Postgres service:

```bash
$ git clone https://github.com/adrianolmedo/genesis.git
$ cp .env.example .env
$ docker-compose up -d --build postgres
```

**Run migrations**

```bash
./migrations.sh
```

**Up application service:**

```bash
$ make genrsa
$ docker-compose up -d --build app
```
