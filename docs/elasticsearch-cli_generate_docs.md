## elasticsearch-cli generate docs

Generates the command tree documentation

### Synopsis


Generates the command tree documentation

```
elasticsearch-cli generate docs [flags]
```

### Options

```
  -h, --help              help for docs
  -l, --location string   Set the location of the generated output (default "./docs")
```

### Options inherited from parent commands

```
      --cluster string      config name, used to have multiple clusters configured in $HOME/.elasticsearch-cli/<env> (default "default")
      --host string         default elasticsearch URL (default "http://localhost")
  -p, --pass string         password to use to authenticate (If not specified, will look for ES_PASS environment variable)
      --poll-interval int   interval on which to poll Elasticsearch to provide index autocompletion (default 10)
      --port int            default elasticsearch port to use (default 9200)
  -t, --timeout int         http client timeout to the remote endpoint (default 10)
  -u, --user string         username to use to authenticate (If not specified look for ES_USER environment variable)
  -v, --verbose             enable verbose mode
```

### SEE ALSO
* [elasticsearch-cli generate](elasticsearch-cli_generate.md)	 - Generates elasticsearch-cli completions and docs

