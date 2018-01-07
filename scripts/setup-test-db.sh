#!/bin/bash -ex

# Run this script inside scripts directory

function run_mysql_stmt()
(
    stmt=$1
    mysql -h localhost --protocol=TCP -u root -ppassword -e "$stmt"
)

#############################
# Populate database
#############################
run_mysql_stmt "DROP DATABASE IF EXISTS \`pinchito-test\`;"

TMPFILE=$(mktemp)
bunzip2 -c ../tests/test-db.sql.bz2 > ${TMPFILE}

run_mysql_stmt "source ${TMPFILE}"

#############################
# Application
#############################
echo '{ "Db" : { "Name" : "pinchito-test", "User" : "root", "Pass" : "password", "Protocol": "tcp", "Charset" : "latin1" }, "Auth": { "Token": "auth-token-test" } }' > ../conf/settings.json
