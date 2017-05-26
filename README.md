# Elasticsearch-cli

`elasticsearch-cli` provides a Kibana console-like interface to interact with ElasticSearch

# Features

* Cli interface, this is a one-off execution
* Interactive console-like execution
* REPL autocompletion
* Persistent history
* Authentication support

# Installation

To install `elasticsearch-cli` you can either grab the latest binaries in the [release page](https://github.com/elastic/elasticsearch-cli/releases)
or install the latest and most recent commit from the source code

## Latest build

`elasticsearch-cli` will be placed in /usr/local/bin/elasticsearch-cli

```sh
git clone https://github.com/marclop/elasticsearch-cli
cd elasticsearch-cli
make install
```

# Default configuration

There are a few configuration flags you can pass to `elasticsearch-cli`:

```sh
$ elasticsearch-cli -help
Usage of elasticsearch-cli:
  -host string
    	Set the ElasticSearch host url (default "http://localhost")
  -pass string
    	Password for HTTP basic auth
  -poll int
    	Set the poll interval for index autodiscovery (default 10)
  -port int
    	Set the Elasticsearch Port (default 9200)
  -timeout int
    	Set the HTTP client timeout (default 10)
  -user string
    	Username for HTTP basic auth
  -verbose
    	Verbose request/response information
```

# Usage

`elasticsearch-cli`'s usage is very intuitive, the execution is split between non-interactive and interactive, which are composed by 3 request arguments:

1. Method
2. URL
3. Body of the request

Non-interactive example:

```sh
$ elasticsearch-cli GET /
{
  "name": "GNBXbv5",
  "cluster_name": "elasticsearch",
  "cluster_uuid": "g5swow-2SHaCA6zPVvf1dQ",
  "version": {
    "number": "5.2.1",
    "build_hash": "db0d481",
    "build_date": "2017-02-09T22:05:32.386Z",
    "build_snapshot": false,
    "lucene_version": "6.4.1"
  },
  "tagline": "You Know, for Search"
}
$ elasticsearch-cli -verbose GET
Method:       GET
URL:          /
Response:     200 OK
Content-Type: application/json

{
  "name": "GNBXbv5",
  "cluster_name": "elasticsearch",
  "cluster_uuid": "g5swow-2SHaCA6zPVvf1dQ",
  "version": {
    "number": "5.2.1",
    "build_hash": "db0d481",
    "build_date": "2017-02-09T22:05:32.386Z",
    "build_snapshot": false,
    "lucene_version": "6.4.1"
  },
  "tagline": "You Know, for Search"
}
$ elasticsearch-cli GET _cat
=^.^=
/_cat/pending_tasks
/_cat/repositories
/_cat/segments
/_cat/segments/{index}
/_cat/health
/_cat/nodes
/_cat/allocation
/_cat/indices
/_cat/indices/{index}
/_cat/aliases
/_cat/aliases/{alias}
/_cat/templates
/_cat/plugins
/_cat/count
/_cat/count/{index}
/_cat/tasks
/_cat/nodeattrs
/_cat/thread_pool
/_cat/thread_pool/{thread_pools}/_cat/master
/_cat/fielddata
/_cat/fielddata/{fields}
/_cat/snapshots/{repository}
/_cat/recovery
/_cat/recovery/{index}
/_cat/shards
/_cat/shards/{index}
```

## Interactive mode

```sh
$ elasticsearch-cli -verbose
Elasticsearch-cli> GET
Method:       GET
URL:          /
Response:     200 OK
Content-Type: application/json

{
  "name": "GNBXbv5",
  "cluster_name": "elasticsearch",
  "cluster_uuid": "g5swow-2SHaCA6zPVvf1dQ",
  "version": {
    "number": "5.2.1",
    "build_hash": "db0d481",
    "build_date": "2017-02-09T22:05:32.386Z",
    "build_snapshot": false,
    "lucene_version": "6.4.1"
  },
  "tagline": "You Know, for Search"
}
Elasticsearch-cli> exit
$ elasticsearch-cli
Elasticsearch-cli> GET
Method:       GET
URL:          /

{
  "name": "GNBXbv5",
  "cluster_name": "elasticsearch",
  "cluster_uuid": "g5swow-2SHaCA6zPVvf1dQ",
  "version": {
    "number": "5.2.1",
    "build_hash": "db0d481",
    "build_date": "2017-02-09T22:05:32.386Z",
    "build_snapshot": false,
    "lucene_version": "6.4.1"
  },
  "tagline": "You Know, for Search"
}
Elasticsearch-cli> exit
```

### Change configuration

While in interactive mode you an choose to change the application's configuration at any time:

```sh
$ elasticsearch-cli
Elasticsearch-cli> get
Method:       GET
URL:          /

{
  "name": "GNBXbv5",
  "cluster_name": "elasticsearch",
  "cluster_uuid": "g5swow-2SHaCA6zPVvf1dQ",
  "version": {
    "number": "5.2.1",
    "build_hash": "db0d481",
    "build_date": "2017-02-09T22:05:32.386Z",
    "build_snapshot": false,
    "lucene_version": "6.4.1"
  },
  "tagline": "You Know, for Search"
}
Elasticsearch-cli> set port 9201
Elasticsearch-cli> get
Method:       GET
URL:          /

{
  "name": "hIzXUZY",
  "cluster_name": "elasticsearch",
  "cluster_uuid": "g5swow-2SHaCA6zPVvf1dQ",
  "version": {
    "number": "5.2.1",
    "build_hash": "db0d481",
    "build_date": "2017-02-09T22:05:32.386Z",
    "build_snapshot": false,
    "lucene_version": "6.4.1"
  },
  "tagline": "You Know, for Search"
}
Elasticsearch-cli> exit
```

## Usage with jq

Of course if you feel like combining the power of Elasticsearch with `jq` for response filtering you can do so.

```sh
$ elasticsearch-cli GET | jq '.version.number'
"5.2.1"
$ elasticsearch-cli GET | jq '.name'
"GNBXbv5"
$ elasticsearch-cli -port 9201 GET | jq '.name'
"hIzXUZY"
```

# Contributing

## Setting up the environment

Elasticsearch-cli is written in [Go](http://golang.org/), so you'll need the latest version of Golang if you want to contribute.
You will also need the latest version of Docker to be able to run the acceptance tests.

## Running all the tests

Issuing `make test` will run the combination of `unit` and `acceptance` tests. If you want a specific test, just use either target.

```sh
$ make test
-> Running unit tests for elasticsearch-cli...
ok  	github.com/elastic/elasticsearch-cli/app	0.010s	coverage: 33.3% of statements
ok  	github.com/elastic/elasticsearch-cli/cli	0.015s	coverage: 63.9% of statements
ok  	github.com/elastic/elasticsearch-cli/client	0.027s	coverage: 83.3% of statements
?   	github.com/elastic/elasticsearch-cli/elasticsearch	[no test files]
ok  	github.com/elastic/elasticsearch-cli/poller	0.008s	coverage: 81.5% of statements
ok  	github.com/elastic/elasticsearch-cli/utils	0.006s	coverage: 100.0% of statements
?   	github.com/elastic/elasticsearch-cli	[no test files]
-> Installing elasticsearch-cli dependencies...
[..]
-> Building elasticsearch-cli...
Number of parallel builds: 7

-->    darwin/amd64: github.com/elastic/elasticsearch-cli
=> Starting Elasticsearch 1.7... Done.
-> Running acceptance tests for elasticsearch-cli in Elasticsearch 1.7...
ok  	github.com/elastic/elasticsearch-cli	1.276s
-> Killing Docker container elasticsearch-cli_es_1.7
=> Starting Elasticsearch 2.4... Done.
-> Running acceptance tests for elasticsearch-cli in Elasticsearch 2.4...
ok  	github.com/elastic/elasticsearch-cli	1.421s
-> Killing Docker container elasticsearch-cli_es_2.4
=> Starting Elasticsearch 5.4... Done.
-> Running acceptance tests for elasticsearch-cli in Elasticsearch 5.4...
ok  	github.com/elastic/elasticsearch-cli	3.566s
-> Killing Docker container elasticsearch-cli_es_5.4
```

# TODOs

- [X] REPL Auto-discover indices
- [X] Use logger
- [X] Create Unit tests
- [X] Acceptance tests
- [ ] Improve help Flag
- [ ] Configuration files
- [ ] CI
- [ ] Bulk API
