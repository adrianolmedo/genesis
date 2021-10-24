
# Practice RESTful API in Go

My first prototype of RESTful API written in Go based in [Dao Pattern](https://github.com/adrianolmedo/design-patterns/tree/main/go/dao-pattern). Whit persistence to MySQL and Postgres.

## Content

* [Run](#run)
* [Endpoints](#endpoints)
  * [Sign Up](#sign-up)
  * [Get user by ID](#get-user-by-id)
  * [Login](#login)
  * [Update user by ID](#update-user-by-id)
  * [Get all users](#get-all-users)
  * [Delete user by ID](#delete-user-by-id)

## TO-DO:

- [ ] Connect login to storage.

## Run:

1- Prepare certificates for JWT:

```bash
$ git clone https://github.com/adrianolmedo/go-restapi-practice.git
$ cd auth/
$ openssl genrsa -out app.sra 1024
$ openssl rsa -in app.sra -pubout > app.sra.pub
```

2- Create a database in Postgres and execute `$ \i postsgres.sql` for install tables.

3- Write configuration of server and database in `config.json` file.

4- Compile and run:

```bash
$ ./compile.sh
$ ./rest -c config.json
```

## Endpoints:

### Sign Up

**POST:** `/v1/signup`

Sing up users or create account. *First Name, Email and Password are fields required.*

Body (JSON):

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "exmaple@gmail.com",
  "password": "1234567@"
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
    "created_at": "2021-07-02T01:20:19.493927615-04:00",
    "updated_at": "0001-01-01T00:00:00Z",
    "deleted_at": "0001-01-01T00:00:00Z",
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "exmaple@gmail.com"
  }
}
```

---

### Get user by ID

**GET:** `/v1/users/:id`

For example to get user with ID 1 make GET request to `/v1/users/1` route. Not required JWT Authorization.

Reponse (200 OK):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": ""
  },
  "data": {
    "created_at": "2021-07-02T01:20:19.493928Z",
    "updated_at": "0001-01-01T00:00:00Z",
    "deleted_at": "0001-01-01T00:00:00Z",
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "exmaple@gmail.com"
  }
}
```

---

### Login

**POST:** `/v1/login`

Login users with data account.

Body (JSON):

```json
{
  "email": "a@g.com",
  "password": "1234567a"
}
```

Reponse (201 Created):

```json
{
  "message_ok": {
    "code": "OK004",
    "content": ""
  },
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkcmlhbm9sbWVkb0BnbWFpbC5jb20iLCJleHAiOjE2MjUxOTUxNTQsImlzcyI6ImFkcmlhbm9sbWVkbyJ9.n-t_X3pQVHa1sz10QqNjQMH6VCtmx9RmBy6J9sjVvbl74cCtFIxFDN9r6M9j4ZjOC_HAvWNdC_mOzOhk0Idrui_18Rqp_D6BqWmthaXAZPIi8qpQ6nAPecm-jQDSxfZj0s9jl0Q32u3oWA0NnDuO3oGoJPQCYWQ_nX3qk4CTFHQ"
  }
}
```

---

### Update user by ID

**PUT:** `/v1/users/:id`

For example to update user with ID 1, make PUT request to `/v1/users/1` route. Required JWT Authorization.

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
    "content": "resource updated"
  },
  "data": {
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2021-07-02T03:23:15.418448805-04:00",
    "deleted_at": "0001-01-01T00:00:00Z",
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "j.doe@go.com",
    "password": "1234567a"
  }
}
```

---

### Get all users

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
      "created_at": "2021-07-02T01:20:19.493928Z",
      "updated_at": "2021-07-02T03:32:19.896399Z",
      "deleted_at": "2021-07-02T03:32:19.896399Z",
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "j.doe@go.com",
      "password": "1234567a"
    },
    {
      "created_at": "2021-07-01T22:16:25.608667Z",
      "updated_at": "0001-01-01T00:00:00Z",
      "deleted_at": "0001-01-01T00:00:00Z",
      "id": 2,
      "first_name": "Adrián",
      "last_name": "Olmedo",
      "email": "adr.ve@a.com",
      "password": "1234567@"
    }
  ]
}
```

---

### Delete user by ID

**DELETE:** `v1/users/:id`

For example to delete user with ID 1, make DELETE request to `/v1/users/1` route. Required JWT Authorization.

Response (200 OK):

```json
{
  "message_ok": {
    "code": "OK002",
    "content": "resource deleted"
  },
  "data": null
}
```
