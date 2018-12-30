#!/usr/bin/env bash
apt install mysql-server
systemctl start mysql
systemctl enable mysql
DBNAME="testgoflow"
USER="remoteuser"
pw=`< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c${1:-32};echo;`
mysql --execute="CREATE DATABASE $DBNAME;"
mysql --execute="GRANT ALL PRIVILEGES ON *.* TO '$USER'@'%' IDENTIFIED BY '$pw';"
mysql --execute="flush privileges;"
echo "Remoteuser password: $pw"
sed -i 's/127.0.0.1/0.0.0.0/' /etc/mysql/mysql.conf.d/mysqld.cnf
service mysql restart