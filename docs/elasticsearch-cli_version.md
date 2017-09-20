## elasticsearch-cli version

prints the version

### Synopsis


prints the version

```
elasticsearch-cli version [flags]
```

### Options

```
  -h, --help   help for version
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
* [elasticsearch-cli](elasticsearch-cli.md)	 - elasticsearch-cli provides a REPL console-like interface to interact with Elasticsearch

