## Content

* [Run with Postgres service](#run-with-postgres-service)

## Run with Postgres service:

```bash
$ git clone https://github.com/adrianolmedo/genesis.git
$ cp .env.example .env
$ docker-compose up -d --build postgres
```

**Join to `psql` and ingress the password `1234567a`:**

```bash
$ docker exec -it postgres /bin/sh
$ psql -U johndoe -d genesis
```

**Install tables:**

```bash
$ \i tables.sql
$ \q
$ exit
```

**Up application service:**

```bash
$ make genrsa
$ docker-compose up -d --build app
```
