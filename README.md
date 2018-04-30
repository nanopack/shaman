[![shaman logo](http://nano-assets.gopagoda.io/readme-headers/shaman.png)](http://nanobox.io/open-source#shaman)  
[![Build Status](https://travis-ci.org/nanopack/shaman.svg)](https://travis-ci.org/nanopack/shaman)
[![GoDoc](https://godoc.org/github.com/nanopack/shaman?status.svg)](https://godoc.org/github.com/nanopack/shaman)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanopack/shaman)](https://goreportcard.com/report/github.com/nanopack/shaman)
<!-- [![GoCover](https://gocover.io/_badge/github.com/nanopack/shaman)](https://gocover.io/github.com/nanopack/shaman) -->

# Shaman

Small, clusterable, lightweight, api-driven dns server.


## Quickstart:
```sh
# Start shaman with defaults (requires admin privileges (port 53))
shaman -s

# register a new domain
shaman add -d nanopack.io -A 127.0.0.1

# perform dns lookup
# OR `nslookup -port=53 nanopack.io 127.0.0.1`
dig @localhost nanopack.io +short
# 127.0.0.1

# Congratulations!
```


## Usage:

### As a CLI
Simply run `shaman <COMMAND>`

`shaman` or `shaman -h` will show usage and a list of commands:

```
shaman - api driven dns server

Usage:
  shaman [flags]
  shaman [command]

Available Commands:
  add         Add a domain to shaman
  delete      Remove a domain from shaman
  list        List all domains in shaman
  get         Get records for a domain
  update      Update records for a domain
  reset       Reset all domains in shaman

Flags:
  -C, --api-crt string            Path to SSL crt for API access
  -a, --api-domain string         Domain of generated cert (if none passed) (default "shaman.nanobox.io")
  -k, --api-key string            Path to SSL key for API access
  -p, --api-key-password string   Password for SSL key
  -H, --api-listen string         Listen address for the API (ip:port) (default "127.0.0.1:1632")
  -c, --config-file string        Configuration file to load
  -O, --dns-listen string         Listen address for DNS requests (ip:port) (default "127.0.0.1:53")
  -d, --domain string             Parent domain for requests (default ".")
  -f, --fallback-dns              Fallback dns server address (ip:port), if not specified fallback is not used
  -i, --insecure                  Disable tls key checking (client) and listen on http (api). Also disables auth-token
  -2, --l2-connect string         Connection string for the l2 cache (default "scribble:///var/db/shaman")
  -l, --log-level string          Log level to output [fatal|error|info|debug|trace] (default "INFO")
  -s, --server                    Run in server mode
  -t, --token string              Token for API Access (default "secret")
  -T, --ttl int                   Default TTL for DNS records (default 60)
  -v, --version                   Print version info and exit

Use "shaman [command] --help" for more information about a command.
```

For usage examples, see [api](api/README.md) and/or [cli](commands/README.md) readme  

### As a Server
To start shaman as a server run:  
`shaman --server`  
An optional config file can also be passed on startup:  
`shaman -c config.json`  

>config.json
>```json
>{
>  "api-domain": "shaman.nanobox.io",
>  "api-crt": "",
>  "api-key": "",
>  "api-key-password": "",
>  "api-listen": "127.0.0.1:1632",
>  "token": "secret",
>  "insecure": false,
>  "l2-connect": "scribble:///var/db/shaman",
>  "ttl": 60,
>  "domain": ".",
>  "dns-listen": "127.0.0.1:53",
>  "log-level": "info",
>  "server": true
>}
>```

#### L2 connection strings

##### Scribble Cacher
The connection string looks like `scribble://localhost/path/to/data/store`.

<!--
#### Redis Cacher
The connection string looks like `redis://[user:password@]host:port/`.

#### Postgresql Cacher
The connection string looks like `postgres://[user@]host/database`.
 -->


## API:

| Route | Description | Payload | Output |
| --- | --- | --- | --- |
| **POST** /records | Adds the domain and full record | json domain object | json domain object |
| **PUT** /records | Update all domains and records (replaces all) | json array of domain objects | json array of domain objects |
| **GET** /records | Returns a list of domains we have records for | nil | string array of domains |
| **PUT** /records/{domain} | Update domain's records (replaces all) | json domain object | json domain object |
| **GET** /records/{domain} | Returns the records for that domain | nil | json domain object |
| **DELETE** /records/{domain} | Delete a domain | nil | success message |

**note:** The API requires a token to be passed for authentication by default and is configurable at server start (`--token`). The token is passed in as a custom header: `X-AUTH-TOKEN`.  

For examples, see [the api's readme](api/README.md)  


## Overview

```sh
+------------+     +----------+     +-----------------+
|            +----->          +----->                 |
| API Server |     |          |     |   Short-Term    |
|            <-----+ Caching  <-----+   (in-memory)   |
+------------+     | And      |     +-----------------+
                   | Database |
+------------+     | Manager  |     +-----------------+
|            +----->          +----->                 |
| DNS Server |     |          |     | Long-Term (L2)  |
|            <-----+          <-----+                 |
+------------+     +----------+     +-----------------+
```


## Data types:
### Domain (Resource):
json:
```json
{
  "domain": "nanopack.io.",
  "records": [
    {
      "ttl": 60,
      "class": "IN",
      "type": "A",
      "address": "127.0.0.1"
    },
    {
      "ttl": 60,
      "class": "IN",
      "type": "A",
      "address": "127.0.0.2"
    }
  ]
}
```

Fields:
- **domain**: Domain name to resolve
- **records**: Array of address records
  - **ttl**: Seconds a client should cache for
  - **class**: Record class
  - **type**: Record type
    - A - Address record
    - CNAME - Canonical name record
    - MX - Mail exchange record
    - [Many more](https://en.wikipedia.org/wiki/List_of_DNS_record_types) - may or may not work as is
  - **address**: Address domain resolves to
    - <sup>note: Special rules apply in some cases. E.g. MX records require a number "10 mail.google.com"</sup>

### Error:
json:
```json
{
  "err": "exit status 2: unexpected argument"
}
```

Fields:
 - **err**: Error message

### Message:
json:
```json
{
  "msg": "Success"
}
```

Fields:
 - **msg**: Success message


## Todo
- atomic local cache updates
- export in hosts file format
- improve scribble add (adding before stored in cache overwrites)


## Changelog
- v0.0.2 (May 11, 2016)
  - Refactor to allow multiple records per domain and more fully utilize dns library
- v0.0.3 (May 12, 2016)
  - Tests for DNS server
  - Start Server Insecure
- v0.0.4 (Aug 16, 2016)
  - Postgresql as a backend


## Contributing
Contributions to shaman are welcome and encouraged. Shaman is a [Nanobox](https://nanobox.io) project and contributions should follow the [Nanobox Contribution Process & Guidelines](https://docs.nanobox.io/contributing/).


[![oss logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
