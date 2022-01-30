#!/bin/sh

openssl genrsa -out app.sra 1024
openssl rsa -in app.sra -pubout > app.sra.pub
