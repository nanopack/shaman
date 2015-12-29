[![shaman logo](http://nano-assets.gopagoda.io/readme-headers/shaman.png)](http://nanobox.io/open-source#shaman)
[![Build Status](https://travis-ci.org/nanopack/shaman.svg)](https://travis-ci.org/nanopack/shaman)

# shaman

Small, lightweight, api-driven dns server.

## Status

Working

## Todo
- Logging
- Tests
- Read configuration from file

## Server
```
Usage:
   [flags]
   [command]

Available Commands:
  add         Add entry into shaman database
  remove      Remove entry from shaman database
  show        Show entry in shaman database
  update      Update entry in shaman database
  list        List entries in shaman database

Flags:
  -c, --api-crt="": Path to SSL crt for API access
  -H, --api-host="127.0.0.1": Listen address for the API
  -k, --api-key="": Path to SSL key for API access
  -p, --api-key-password="": Password for SSL key
  -P, --api-port="8443": Listen address for the API
  -t, --api-token="": Token for API Access
  -d, --domain=".": Parent domain for requests
  -h, --help[=false]: help for 
  -O, --host="127.0.0.1": Listen address for DNS requests
  -i, --insecure[=false]: Disable tls key checking
  -1, --l1-connect="map://127.0.0.1/": Connection string for the l1 cache
  -e, --l1-expires=120: TTL for the L1 Cache (0 = never expire)
  -2, --l2-connect="map://127.0.0.1/": Connection string for the l2 cache
  -E, --l2-expires=0: TTL for the L2 Cache (0 = never expire)
  -l, --log-file="": Log file (blank = log to console)
  -L, --log-level="INFO": Log level to use
  -o, --port="8053": Listen port for DNS requests
  -s, --server[=false]: Run in server mode
  -T, --ttl=60: Default TTL for DNS records

Use " [command] --help" for more information about a command.
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

### Commands

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
