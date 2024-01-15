#!/bin/bash

# options:
DEBUG=false
CREATEDB=true
DBNAME="genesis"
DBUSER="johndoe"
SQLDIR="sql"
CONTAINER="postgres"

echo "DEBUG: $DEBUG"
echo "DBNAME: $DBNAME"
echo "DBUSER: $DBUSER"
echo "SQLDIR: "$SQLDIR""
echo "CONTAINER: $CONTAINER"
if [[ "$CREATEDB" == true ]]; then
    echo -e "Create database: Yes\n"
else
    echo -e "Create database: No\n"
fi

# security question
if [[ "$BASH_ARGV" != "-y" ]]; then
    read -r -p "Do you want to continue? [y/n] " response
fi

# execute migration if response is equeal to y/Y or the last argument is "-y"
if [[ "$response" =~ ^[yY]$ ]] || [[ "$BASH_ARGV" == "-y" ]]; then

    # create data base
    if [[ "$CREATEDB" == true ]]; then
        if [[ -z "$CONTAINER" ]]; then
            if [[ "$DEBUG" == true ]]; then
                echo "Creating db without docker"
            else
                sudo -u postgres createdb $DBNAME
            fi
        else
            if [[ "$DEBUG" == true ]]; then
                echo "Creating db with docker"
            else
                echo "Creating database..."
                docker exec -it $CONTAINER createdb --username=$DBUSER --owner=$DBUSER $DBNAME
            fi
        fi
    fi

    # execute sql files from a sql directory
    echo "Running additional migration scripts..."

    # running for local postgres installation
    if [[ -z "$CONTAINER" ]]; then
        for file in "$SQLDIR"/*.sql
        do
            if [[ -f "$file" && "$file" == *".sql" ]]; then
                if [[ "$DEBUG" == true ]]; then
                    echo "Executing SQL file without docker: $file"
                else
                    echo "Executing SQL file: $file"
                    sudo -u postgres psql -d $DBNAME -f $SQLDIR/$file
                fi

            fi
        done
    else
    # running for a docker container with postgres
        for file in "$SQLDIR"/*.sql
        do
            if [[ -f "$file" && "$file" == *".sql" ]]; then
                if [[ "$DEBUG" == true ]]; then
                    echo "Executing SQL file with docker: $file"
                else
                    echo "Executing SQL file: $file"
                    cat $file | docker exec -i $CONTAINER psql -U $DBUSER -d $DBNAME
                fi

            fi
        done
    fi

    echo "Migration executed!"
else
    echo "Migration canceled!"
fi
