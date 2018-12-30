# **This is pre-release, Alpha software. Use at your own risk!**
# Description
<p align="center">
    <img src="https://i.imgur.com/HdIxEOB.png">
</p>
This is a very early release of Goflow, a golang-based netflow collector with a flexible backend.

A list of upcoming features can be found under the issue tracker for this project.

Currently supported frontend/backend combinations are 

Frontend | Backend
------ | ------ |
Netflow | Mysql 
Netflow | Timescaledb


# Prereqs

You need a running backend and associated connection information;

- Server fqdn 
- Username
- Password
- Database name

For Mysql, you could use the free tier of Amazon RDS to get started.

See the SETUP.md instructions within the backend directories in this repo for help
specific to each backend. (i.e ./backends/timescale/SETUP.MD)

# Installation
Goflow requires two files to run;
 - goflow, the binary itself
 - config.yml, the configuration file

The tar releases contain both these.

```bash
# Extract and set perms
tar -xzvf goflow.Linux.AMD64.tar.gz
chmod +x goflow
mv config_example.yml config.yml
# Edit the config.yml file to make specific to your environment
vi config.yml
# Export the required environment variables
export SQL_PASSWORD=your_sql_pw_here
# Run
./goflow
```

In the future, an installation script will be packaged for most systems but for now, you will need to create your own systemd or init scripts to start it.

# Monitoring and utilities
The goflow binary doubles as a client interface, and a JSON API is started at the same time as the daemon.

The API listens on localhost by default, but this can be tuned (see the configuration example).

The Goflow API is not for retrieving flow data, but performing ongoing maintenence and ops on the collector itself.

Goflow help displays a list of options.
```bash
./goflow help
```
