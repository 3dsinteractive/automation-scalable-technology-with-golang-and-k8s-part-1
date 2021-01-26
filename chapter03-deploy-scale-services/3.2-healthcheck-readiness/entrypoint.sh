#!/bin/sh -e

# exec env file if exists, .env might be in /.env or /config/.env
FILE=/.env ; [ -f $FILE ] && . $FILE
FILE=/config/.env ; [ -f $FILE ] && . $FILE

# start main program
exec /main