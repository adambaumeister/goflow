# Timescale
Timescale is a Postgres-based database that does automatic partitioning based on timestamptz fields.

Because of this, it's well suited to Netflow workloads.

https://docs.timescale.com/v1.1/introduction

## Quick start
Goflow source includes a setup script that stands up a basic Timescaledb, useful for testing.

../setup_scripts/timescaledb.sh

Eventually this readme will include optimal settings for tweaking the db specifically for goflow.

