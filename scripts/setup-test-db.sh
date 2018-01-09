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
# User
#############################
# Not supported in mysql 5.5!
# run_mysql_stmt "DROP USER IF EXISTS 'pinchito-test'@'localhost';"

run_mysql_stmt "CREATE USER 'pinchito-test'@'127.0.0.1' IDENTIFIED BY 'p1nt3st';"
run_mysql_stmt "GRANT ALL PRIVILEGES ON \`pinchito-test\`.* TO 'pinchito-test'@'127.0.0.1'"

#############################
# Application
#############################
echo '{ "Db" : { "Name" : "pinchito-test", "User" : "pinchito-test", "Pass" : "p1nt3st", "Protocol": "tcp", "Charset" : "latin1" }, "Auth": { "Token": "auth-token-test" } }' > ../conf/settings.json
