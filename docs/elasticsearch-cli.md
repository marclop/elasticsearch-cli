## elasticsearch-cli

elasticsearch-cli provides a REPL console-like interface to interact with Elasticsearch

### Synopsis


elasticsearch-cli provides a REPL console-like interface to interact with Elasticsearch

```
elasticsearch-cli [flags]
```

### Options

```
      --cluster string      config name, used to have multiple clusters configured in $HOME/.elasticsearch-cli/<env> (default "default")
  -h, --help                help for elasticsearch-cli
      --host string         default elasticsearch URL (default "http://localhost")
      --insecure            skip tls certificate verification (warning: use for testing or development onlu)
  -p, --pass string         password to use to authenticate (If not specified, will look for ES_PASS environment variable)
      --poll-interval int   interval on which to poll Elasticsearch to provide index autocompletion (default 10)
      --port int            default elasticsearch port to use (default 9200)
  -t, --timeout int         http client timeout to the remote endpoint (default 10)
  -u, --user string         username to use to authenticate (If not specified look for ES_USER environment variable)
  -v, --verbose             enable verbose mode
```

### SEE ALSO
* [elasticsearch-cli delete](elasticsearch-cli_delete.md)	 - Performs a DELETE operation against the remote endpoint
* [elasticsearch-cli generate](elasticsearch-cli_generate.md)	 - Generates elasticsearch-cli docs
* [elasticsearch-cli get](elasticsearch-cli_get.md)	 - Performs a GET operation against the remote endpoint
* [elasticsearch-cli head](elasticsearch-cli_head.md)	 - Performs a HEAD operation against the remote endpoint
* [elasticsearch-cli post](elasticsearch-cli_post.md)	 - Performs a POST operation against the remote endpoint
* [elasticsearch-cli put](elasticsearch-cli_put.md)	 - Performs a PUT operation against the remote endpoint
* [elasticsearch-cli version](elasticsearch-cli_version.md)	 - prints the version

