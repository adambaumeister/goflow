# **This is pre-release software. Use at your own risk!**
# Description

![Imgur](https://i.imgur.com/HdIxEOB.png)

This is a very early release of Goflow, a golang-based netflow collector with a flexible backend.

Right now, there are many features not implemented and only one supported frontend/backend (netflow/mysql, respectively).

A list of upcoming features can be found under the issue tracker for this project.

# Prereqs

For this release, you need a running mysql server and the following details from it

- IP address or hostname
- Database name
- Username
- Password

A cost-effective option is Amazon RDS (https://aws.amazon.com/rds/) if you don't have spare servers to run a mysql instance on. Keep in mind that SSL for the backend is not yet implemented so uhhh did I mention this is pre-release software?

# Installation
Goflow requires two files to run;
 - goflow, the binary itself
 - config.yml, the configuration file

The tar releases contain both these.

```bash
# Extract and set perms
tar -xzvf goflow.tar.gz.Linux.AMD64.gz
chmod +x goflow
# Edit the config.yml file to make specific to your environment
vi  config.yml
# Export the required environment variables
export SQL_PASSWORD=your_sql_pw_here
# Run
./goflow
```

In the future, an installation script will be packaged for most systems but for now, you will need to create your own systemd or init scripts to start it.

