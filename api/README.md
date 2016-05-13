[![shaman logo](http://nano-assets.gopagoda.io/readme-headers/shaman.png)](http://nanobox.io/open-source#shaman)  
[![Build Status](https://travis-ci.org/nanopack/shaman.svg)](https://travis-ci.org/nanopack/shaman)

# Shaman

Small, lightweight, api-driven dns server.

## Routes:

| Route | Description | Payload | Output |
| --- | --- | --- | --- |
| **POST** /records | Adds the domain and full record | json domain object | json domain object |
| **PUT** /records | Update all domains and records (replaces all) | json array of domain objects | json array of domain objects |
| **GET** /records | Returns a list of domains we have records for | nil | string array of domains |
| **PUT** /records/{id} | Update domain's records (replaces all) | json domain object | json domain object |
| **GET** /records/{id} | Returns the records for that domain | nil | json domain object |
| **DELETE** /records/{id} | Delete a domain | nil | success message |

## Usage Example:

#### add domain
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records -d \
       '{"domain":"nanopack.io","records":[{"ttl":60,"class":"IN","type":"A","address":"127.0.0.2"}]}'
# {"domain":"nanopack.io.","records":[{"ttl":60,"class":"IN","type":"A","address":"127.0.0.2"}]}
```

#### list domains
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records
# ["nanopack.io"]
```
or add `?full=true` for the full records
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records?full=true
# [{"domain":"nanopack.io.","records":[{"ttl":60,"class":"IN","type":"A","address":"127.0.0.2"}]}]
```

#### update domains
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records -d \
       '[{"domain":"nanobox.io","records":[{"address":"127.0.0.1"}]}]' \
       -X PUT
# [{"domain":"nanobox.io.","records":[{"ttl":60,"class":"IN","type":"A","address":"127.0.0.1"}]}]
```

#### update domain
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records/nanobox.io -d \
       '{"domain":"nanobox.io","records":[{"address":"127.0.0.2"}]}' \
       -X PUT
# {"domain":"nanobox.io.","records":[{"ttl":60,"class":"IN","type":"A","address":"127.0.0.2"}]}
```

#### delete domain
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records/nanobox.io \
       -X DELETE
# {"msg":"success"}
```

#### get domain
```sh
$ curl -k -H "X-AUTH-TOKEN: secret" https://localhost:1632/records/nanobox.io
# {"err":"failed to find record for domain - 'nanobox.io'"}
```

[![oss logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
