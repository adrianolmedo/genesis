#!/bin/bash

# $ bash compile.sh         Generate a normal binary called 'rest'
# $ bash compile.sh debub   Generate a binary called 'rest' for debug

cd cmd/rest

if [[ $1 == "debug" ]]
then
    go build -gcflags "-N -l" .
else
    go build .
fi

mv rest ../../
