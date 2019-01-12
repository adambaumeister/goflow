
# Description
<p align="center">
    <img src="https://i.imgur.com/HdIxEOB.png">
</p>
A golang-based netflow collector with a flexible backend.

A list of upcoming features can be found under the issue tracker for this project.

Currently supported frontend/backend combinations are 

Frontend | Backend
------ | ------ |
Netflow | Mysql 
Netflow | Timescaledb
Netflow | Apache Kafka


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

# Integrations
## Grafana

Goflow integrates natively (i.e - no plugins required) with Grafana when using the Timescale backend type.
 
Grafana transforms the underlying Postgres database into a set of  pretty graphs!

_Note: dummy data shown_

<a href="https://imgur.com/2VTol8w"><img width="300px" src="https://i.imgur.com/2VTol8w.png" title="source: imgur.com" /></a>
<a href="https://imgur.com/DdIhG6k"><img width="300px" src="https://i.imgur.com/DdIhG6k.png" title="source: imgur.com" /></a>
<a href="https://imgur.com/WWpPKch"><img width="300px" src="https://i.imgur.com/WWpPKch.png" title="source: imgur.com" /></a>

For convenience, Goflow provides the Dashboards (/grafana_db/*.json) and some code to setup Grafana correctly.

You need:
* A timescale backend, already configured in config.yml
* A running grafana instance
* <a href="http://docs.grafana.org/http_api/auth/">An API key</a> 

With these requirements met, the dashboards/datasources can be setup as below.
```bash
# Make sure the required env-var for timescale is exported
export SQL_PASSWORD=your-sql-password
./goflow configure-grafana http://[ your-grafana-server ] [ your-api-key ] [ dashboard-directory ]
```

# Performance
## Benchmarks
Each release of Goflow is benchmarked in a test environment. 

Currently the most efficient backend is **timescaledb**.

The environment setup is;
* Type: AWS T2.Micro
* CPU: 1vCPU
* Memory: 1GB
* Storage: EBS SSD (Non-provisioned)


<img src="https://docs.google.com/spreadsheets/d/e/2PACX-1vRzmIcecD3Q-bhAaPSu46EDgxb680ejwWB06Gr9OmabVUFR-GtkVm3PCvUoI6o4Fw0YBW1KTQQjarwn/pubchart?oid=341468645&format=image">

Both network latency and storage have a large impact on performance. The benchmarks above are running with Goflow on the same server as the backends.

## Notes on tuning and hardware requirements
Netflow, unsurprisingly, generates a lot of data.

It's unreasonable to try and estimate compute and storage requirements ahead of time, as this sort of thing 
is hard to quantify, as it's entirely based on how many flows you're exporting which you probably don't already know!

Instead, first decide what your goals are for the data. Specifically, decide how much data you want to _store_, 
then what timeframe you want to be able to query _quickly_ and finally, what constitutes "_quickly_."

After you have that decided, understand how increasing hardware attributes affects each decision:
* More memory allows for more caching, which allows you to run short time range queries very efficiently. In a real environment,
doubling the memory of a timescaledb instance reduced a SELECT query runtime by more than 10x. 
* More cores will make complicated sorts, joins, and other SQL manipulations faster when reading from memory.
* Faster storage improves the speed of queries that cannot be cached or are not yet cached.

In practice? **You should give your SQL server access to an amount of shared memory equal to the amount of data that 
fits in the timeframe you would like to query quickly. If it is unreasonable to fit that amount of data into memory you need to increase storage READ speeds.**

Don't forget to actually tune your database after installation (we've all done it...)! Timescale offers a super good utility for doing it automatically:
https://github.com/timescale/timescaledb-tune

## Example

To illustrate the above points, imagine example.corp wants to store 6 months of netflow data. They would like to query 24 hours worth as quickly as possible
to use on their auto-refreshing wallboards in the office, which refresh once every 30 seconds.

From experimentation, they run at 2k flows per second average with each flow attributed to approximately 150 bytes/flow on disk.

(300B*2000)*86400 = 25GB/day
 
A reccomended hardware setup would be;
CPU: 6-8 cores
Memory: 32GB Minimum 
Storage: 4.5TB of disk benchmarked to at least 100MB/s read.




# Environment variables
Below are a list of all the supported environment variables and the scope in which they are relevant.

Envar | Scope | Purpose
------ | ------ | -------
GOFLOW_CONFIG | * | Path to configuration file (config.yml)
SQL_PASSWORD | Timescaledb, mysql | SQL password
KAFKA_SERVER | Kafka | Kafka server 
KAFKA_TOPIC | Kafka | Topic to publish to
SSL | Kafka | SSL Enabled/disabled
SSL_VERIFY | Kafka | SSL Verfification 
SASL_USER | Kafka  | Kafka username
SASL_PASSWORD | Kafka | Kafka password
