#!/bin/sh

./rest -port 1323 -cors 127.0.0.1 \
    -dbengine postgres \
    -dbhost 127.0.0.1 \
    -dbport 5432 \
    -dbuser postgres \
    -dbpass 1234567@ \
    -dbname go_practice_restapi