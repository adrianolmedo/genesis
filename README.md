
# Practice RESTful API in Go

My first prototype of RESTful API written in Go whit persistence to MySQL or Postgres.

## TO-DO

[-] Add logger when creating invoice.

[-] Add logger when adding a product to an invoice.

## Architecture

![DDD style architecture](https://i.imgur.com/QPvGl0K.png)

## Content

* [Run with MySQL service](#run-with-mysql-service)
* [Run with Postgres service](#run-with-postgres-service)
* [Endpoints](#endpoints)
  * [Sign Up](#sign-up)
  * [Get user by ID](#get-user-by-id)
  * [Login](#login)
  * [Update user by ID](#update-user-by-id)
  * [Get all users](#get-all-users)
  * [Delete user by ID](#delete-user-by-id)

## Run with MySQL service:

```bash
$ git clone https://github.com/adrianolmedo/genesis.git
$ make genrsa
$ docker-compose up -d --build app mysql
```

## Run with Postgres service:

**1- Prepare database:**

```bash
$ git clone https://github.com/adrianolmedo/genesis.git
$ vim .env
```

Change the following values (remember that the postgres default internal port is 5432):

```
# Accepted values: mysql - postgres
DBENGINE=postgres

# Accepted values: mysql - postgres
DBHOST=postgres

# For only internal usage in the container. Accepted values: 3306 - 5432
DBPORT=5432
```

Then build the container:


```
$ docker-compose up -d --build postgres
```

**1.1- Join to `psql` and ingress the password `1234567a`:**

```bash
$ docker exec -it postgres /bin/sh
$ psql -U johndoe -d genesis
```

**1.2- Install tables:**

```bash
$ \i tables.sql
$ \q
$ exit
```

**2- Up application service:**

```bash
$ openssl genrsa -out app.sra 1024
$ openssl rsa -in app.sra -pubout > app.sra.pub
$ docker-compose up -d --build app
```

## Endpoints:

### **Sign Up**

**POST:** `/v1/signup`

Sing up users or create account. *First Name, Email and Password are fields required.*

Body (JSON):

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "jdoe@go.com",
  "password": "1234567b"
}
```

Reponse (201 Created):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": "user created"
  },
  "data": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "jdoe@go.com"
  }
}
```

---

### **Get user by ID**

**GET:** `/v1/users/:id`

Not required JWT Authorization.

Reponse (200 OK):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": ""
  },
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "jdoe@go.com"
  }
}
```

---

### **Login**

**POST:** `/v1/login`

Login users with data account.

Body (JSON):

```json
{
  "email": "jdoe@go.com",
  "password": "1234567b"
}
```

Reponse (201 Created):

```json
{
  "message_ok": {
    "code": "OK004",
    "content": "logged"
  },
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Impkb2VAZ28uY29tIiwiZXhwIjoxNjQ0NTc5NTA1LCJpc3MiOiJhZHJpYW5vbG1lZG8ifQ.qEYFi_ffDaI0aek01REQPS0L8dcTB6mteq09NK8PXf1fPCRp0H3EvIyjCRuJL6zddIxPsaUTi2-LERORc4-GsVwjA-qRPf0IpDwY75YroIC8LfZ_gd3icbxP1fTBy2ZQLy1cHLX11gBvxsXle-LX4dbIMmv81ulsbabkcVY_Vrw"
  }
}
```

---

### **Update user by ID**

**PUT:** `/v1/users/:id`

Required JWT Authorization.

Body (JSON):

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "j.doe@go.com",
  "password": "1234567a"
}
```

Response (200 OK):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": "user updated"
  },
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "j.doe@go.com"
  }
}
```

---

### **Get all users**

**GET:** `v1/users`

Required JWT Authorization.

Response (200 OK):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": ""
  },
  "data": [
    {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "j.doe@go.com"
    },
    {
      "id": 3,
      "first_name": "Jane",
      "last_name": "Doe",
      "email": "qwerty@example.com"
    }
  ]
}
```

---

### **Delete user by ID**

**DELETE:** `v1/users/:id`

Required JWT Authorization.

Response (200 OK):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": "user deleted"
  }
}
```
