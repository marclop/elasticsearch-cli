[![Build Status](https://travis-ci.org/marclop/elasticsearch-cli.svg?branch=master)](https://travis-ci.org/marclop/elasticsearch-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/marclop/elasticsearch-cli)](https://goreportcard.com/report/github.com/marclop/elasticsearch-cli) [![codecov](https://codecov.io/gh/marclop/elasticsearch-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/marclop/elasticsearch-cli)

# Elasticsearch-cli

`elasticsearch-cli` provides a Kibana console-like interface to interact with ElasticSearch

## Features

* Cli interface, this is a one-off execution
* Interactive console-like execution
* REPL autocompletion
* Persistent history
* Authentication support

## Installation

To install `elasticsearch-cli` you can either grab the latest binaries in the [release page](https://github.com/marclop/elasticsearch-cli/releases)
or install the latest and most recent commit from the source code

### Latest build

`elasticsearch-cli` will be placed in ${GOPATH}/bin/elasticsearch-cli

```sh
git clone https://github.com/marclop/elasticsearch-cli
cd elasticsearch-cli
make install
```

## Default configuration

There are a few configuration flags you can pass to `elasticsearch-cli`:

```console
$ elasticsearch-cli --help
elasticsearch-cli provides a REPL console-like interface to interact with Elasticsearch

Usage:
  elasticsearch-cli [flags]
  elasticsearch-cli [command]

Available Commands:
  delete      Performs a DELETE operation against the remote endpoint
  get         Performs a GET operation against the remote endpoint
  head        Performs a HEAD operation against the remote endpoint
  help        Help about any command
  post        Performs a POST operation against the remote endpoint
  put         Performs a PUT operation against the remote endpoint
  version     prints the version

Flags:
      --cluster string      config name, used to have multiple clusters configured in $HOME/.elasticsearch-cli/<env> (default "default")
  -h, --help                help for elasticsearch-cli
      --host string         default elasticsearch URL (default "http://localhost")
  -p, --pass string         password to use to authenticate (If not specified, will look for ES_PASS environment variable)
      --poll-interval int   interval on which to poll Elasticsearch to provide index autocompletion (default 10)
      --port int            default elasticsearch port to use (default 9200)
  -t, --timeout int         http client timeout to the remote endpoint (default 10)
  -u, --user string         username to use to authenticate (If not specified look for ES_USER environment variable)
  -v, --verbose             enable verbose mode

Use "elasticsearch-cli [command] --help" for more information about a command.
```

## Configuration Settings and Precedence

In order for `elasticsearch-cli` to be able to comunicate with an Elasticsearch cluster, it needs to have a set of configuration parameters set, which could either be defined in a configuration file, using environment variables or at runtime or using the cli's global flags. The hirearchy is as follows listed from higher precedence to lower:

1. Command line flags (`--host`, `--user`, `--pass`, `--verbose`, etc).
2. Environment variables.
3. Shared configuration file (`$HOME/.elasticsearch-cli/default.<json|toml|yaml|hcl>`).


## Configuration variables

Before you can start calling your Elasticsearch from the binary, you will need to configure it. Here's an example `YAML` configuration file (`$HOME/.elasticsearch-cli/default.yaml`) that will effectively point and configure the binary for Elastic Cloud:

```yaml
host: https://9a980720ff16d55e3507bfc875551626.us-east-1.aws.found.io:9243
region: us-east-1
user: marc
pass: mypass
```

You could also specify the same settings using environment variables, or to override some settings of the `YAML` file, to do you'll need to prefix the configuration parameter with `ES_` and capitalize the setting, i.e. `ES_HOST` or `ES_USER`.

```sh
export ES_CONFIG=mycluster
```

Last but not least, you can override any of the settings using the CLI flags.

```sh
elasticsearch-cli --config mycluster <COMMAND>
```

## Multple configuration support

`elasticsearch-cli` supports the notion of having multiple cluster configuration files out of the box. It uses those to manage credentials and settings.
By default it will use `$HOME/.elasticsearch-cli/default.<json|toml|yaml|hcl>`, but when the `--config` flag is specified, it will use the `--config` specified value as the file name inside `$HOME/.elasticsearch-cli`:

```console
# Default behaviour
$ elasticsearch-cli version
Using config file: /Users/marc/.elasticsearch-cli/default.yaml
[...]

# when an environment is specified, the configuration file used will change
$ elasticsearch-cli version --config cluster
Using config file: /Users/marc/.elasticsearch-cli/cluster.yaml
[...]
```

# Usage

`elasticsearch-cli`'s usage is very intuitive, the execution is split between non-interactive and interactive, which is composed by 3 request arguments:

1. Method
2. URL
3. Body

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
elasticsearch> GET
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
elasticsearch> exit
$ elasticsearch-cli
elasticsearch> GET
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
elasticsearch> exit
```

### Change configuration

While in interactive mode you an choose to change the application's configuration at any time:

```sh
$ elasticsearch-cli
elasticsearch> get
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
elasticsearch> set port 9201
elasticsearch> get
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
elasticsearch> exit
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
ok  	github.com/marclop/elasticsearch-cli/app	0.010s	coverage: 33.3% of statements
ok  	github.com/marclop/elasticsearch-cli/cli	0.015s	coverage: 63.9% of statements
ok  	github.com/marclop/elasticsearch-cli/client	0.027s	coverage: 83.3% of statements
?   	github.com/marclop/elasticsearch-cli/elasticsearch	[no test files]
ok  	github.com/marclop/elasticsearch-cli/poller	0.008s	coverage: 81.5% of statements
ok  	github.com/marclop/elasticsearch-cli/utils	0.006s	coverage: 100.0% of statements
?   	github.com/marclop/elasticsearch-cli	[no test files]
-> Installing elasticsearch-cli dependencies...
[..]
-> Building elasticsearch-cli...
Number of parallel builds: 7

-->    darwin/amd64: github.com/marclop/elasticsearch-cli
=> Starting Elasticsearch 1.7... Done.
-> Running acceptance tests for elasticsearch-cli in Elasticsearch 1.7...
ok  	github.com/marclop/elasticsearch-cli	1.276s
-> Killing Docker container elasticsearch-cli_es_1.7
=> Starting Elasticsearch 2.4... Done.
-> Running acceptance tests for elasticsearch-cli in Elasticsearch 2.4...
ok  	github.com/marclop/elasticsearch-cli	1.421s
-> Killing Docker container elasticsearch-cli_es_2.4
=> Starting Elasticsearch 5.4... Done.
-> Running acceptance tests for elasticsearch-cli in Elasticsearch 5.4...
ok  	github.com/marclop/elasticsearch-cli	3.566s
-> Killing Docker container elasticsearch-cli_es_5.4
```

### SEE ALSO
* [elasticsearch-cli docs](./docs/elasticsearch-cli.md)
