# Genesis

REST API of CRM, its architecture is inspired (not based) on Hexagonal Architecture.

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
$ make swagger
$ docker-compose up -d --build app
```
