[![shaman logo](http://nano-assets.gopagoda.io/readme-headers/shaman.png)](http://nanobox.io/open-source#shaman)
[![Build Status](https://travis-ci.org/nanopack/shaman.svg)](https://travis-ci.org/nanopack/shaman)

# shaman

Small, lightweight, api-driven dns server.

## Status

Incomplete

## Todo

- Get Started
- Documentation
- Tests

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
