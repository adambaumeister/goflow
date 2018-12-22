#!/usr/bin/env bash
sh -c "echo 'deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -c -s`-pgdg main' >> /etc/apt/sources.list.d/pgdg.list"
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
apt-get update

add-apt-repository ppa:timescale/timescaledb-ppa
apt-get update

# To install for PG 10.2+
apt install timescaledb-postgresql-10

# Prepare for timescaledb
echo "shared_preload_libraries = 'timescaledb'" >> /etc/postgresql/10/main/postgresql.conf
echo "listen_addresses = '*'" >> /etc/postgresql/10/main/postgresql.conf
service postgresql restart

# Setup postgresql with a default database
pw=`< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c${1:-32};echo;`
su postgres -c "psql -c \"ALTER USER postgres PASSWORD '$pw'\""
echo "Generated postgres user pw is $pw"
echo "host    all             remoteuser      203.206.36.44/32        md5" >> /etc/postgresql/10/main/pg_hba.conf
su postgres -c "createuser --pwprompt remoteuser"
su postgres -c "psql -c \"CREATE DATABASE testgoflow\""
su postgres -c "psql testgoflow -c \"CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE\""