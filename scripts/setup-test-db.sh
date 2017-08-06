#!/bin/bash -ex

# Run this script inside scripts directory

function run_mysql_stmt()
(
    stmt=$1
    mysql -u root -e "$stmt"
)

#############################
# Populate database
#############################
run_mysql_stmt "DROP DATABASE \`pinchito-test\` IF EXISTS;"

TMPFILE=$(mktemp)
bunzip2 -c ../tests/test-db.sql.bz2 > ${TMPFILE}

run_mysql_stmt "source ${TMPFILE}"

#############################
# User
#############################
run_mysql_stmt "DROP USER IF EXISTS 'pinchito-test'@'localhost';"

run_mysql_stmt "CREATE USER 'pinchito-test'@'localhost' IDENTIFIED BY 'p1nt3st';"
run_mysql_stmt "GRANT ALL PRIVILEGES ON \`pinchito-test\`.* TO 'pinchito-test'@'localhost'"

#############################
# Application
#############################
echo '{ "Db" : { "Name" : "pinchito-test", "User" : "pinchito-test", "Pass" : "p1nt3st" } }' > ../conf/settings.json
