[![shaman logo](http://nano-assets.gopagoda.io/readme-headers/shaman.png)](http://nanobox.io/open-source#shaman)
[![Build Status](https://travis-ci.org/nanopack/shaman.svg)](https://travis-ci.org/nanopack/shaman)

# shaman

Small, lightweight, api-driven dns server.

## Status

Working

## Todo
- Logging
- Tests
    - api
    - cli
- Read configuration from file

## Server
```
Usage of ./shaman:
  -address string
    	Listen address for DNS requests (default "127.0.0.1:8053")
  -api-address string
    	Listen address for the API (default "127.0.0.1:8443")
  -api-crt string
    	Path to SSL crt for API access
  -api-key string
    	Path to SSL key for API access
  -api-key-password string
    	Password for SSL key
  -api-token string
    	Token for API Access
  -domain string
    	Parent domain for requests (default "example.com")
  -l1-connect string
    	Connection string for the l1 cache (default "map://127.0.0.1/")
  -l1-expires int
    	TTL for the L1 Cache (0 = never expire) (default 120)
  -l2-connect string
    	Connection string for the l2 cache (default "map://127.0.0.1/")
  -l2-expires int
    	TTL for the L2 Cache (0 = never expire)
  -log-file string
    	Log file (blank = log to console)
  -log-level string
    	Log level to use (default "INFO")
  -ttl int
    	Default TTL for DNS records (default 60)
```
### L1 and L2 connection strings

#### In-Memory Map Cacher
This is the default cacher. If the connection string doesn't match any of the other's, it will use this one.

#### Postgresql Cacher
The connection string looks like `postgres://user@host/database` and more [docs here](https://godoc.org/github.com/lib/pq). This string gets passed into the sql driver without modification.

#### Redis Cacher
The connection string looks like `redis://user:password@host:port/`. The user is not really used, but only there if there is a password on the redis-server.

#### Scribble Cacher
The connection string looks like `scribble://localhost/path/to/data/store`. Scribble only cares about the path part of the URI to determine where it should place the files.

## API
The API is a web based API. The API uses TLS and a token for security and authentication.

### API token
The API requires a token to be passed for authentication. This token is set when the server is started. The token is passed in the header as `X-NANOBOX-TOKEN`.

#### Add
POST to `/records/[record type]/[domain]`
A `value` must be posted. Currently it has to be past as a query string rather than part of the post body like `/records/[record type]/[domain]?value=[value]`. This is an issue that should be fixed.

#### Remove
DELETE to `/records/[record type]/[domain]`

#### Show
GET to `/records/[record type]/[domain]`

#### Update
PUT to `/records/[record type]/[domain]`
A `value` must be put. Currently it has to be past as a query string rather than part of the put body like `/records/[record type]/[domain]?value=[value]`. This is an issue that should be fixed.

#### List
GET to `/records`

## CLI
```
Usage:
   [command]

Available Commands:
  add         Add entry into shaman database
  remove      Remove entry from shaman database
  show        Show entry in shaman database
  update      Update entry in shaman database
  list        List entries in shaman database

Flags:
  -A, --auth="": Shaman auth token
  -h, --help[=false]: help for 
  -H, --host="127.0.0.1": Shaman hostname/IP
  -i, --insecure[=false]: Disable tls key checking
  -p, --port=8443: Shaman admin port

Use " [command] --help" for more information about a command.
```

#### add
`add [Record Type] [Domain] [Value]`

#### remove
`remove [Record Type] [Domain]`

#### show
`show [Record Type] [Domain]`

#### update
`update [Record Type] [Domain] [Value]`

#### list 
`list`

### Notes

#### Using nslookup to test
The port can be set with `set port=8053` and the server with `server 127.0.0.1`
```
$ nslookup
> set port=8053
> server 127.0.0.1
Default server: 127.0.0.1
Address: 127.0.0.1#8053
> test.com
Server:		127.0.0.1
Address:	127.0.0.1#8053

Non-authoritative answer:
*** Can't find test.com: No answer
> exit
```

#### Overview

```
+------------+     +----------+     +-----------------+
|            +----->          +----->                 |
| API Server |     |          |     | Short-Term (L1) |
|            <-----+ Caching  <-----+                 |
+------------+     | And      |     +-----------------+
                   | Database |
+------------+     | Manager  |     +-----------------+
|            +----->          +----->                 |
| DNS Server |     |          |     | Long-Term (L2)  |
|            <-----+          <-----+                 |
+------------+     +----------+     +-----------------+
```

[![shaman logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
