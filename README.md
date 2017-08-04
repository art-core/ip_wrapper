# ip_wrapper - Execute a sensu check for all IPs from hostname lookup
[![GitHub version](https://badge.fury.io/gh/wywygmbh%2Fip_wrapper.svg)](https://badge.fury.io/gh/wywygmbh%2Fip_wrapper)
[![Build Status](https://travis-ci.org/wywygmbh/ip_wrapper.svg?branch=master)](https://travis-ci.org/wywygmbh/ip_wrapper)
### Motivation

If a service has multiple IPs (e.g. behind a load balancer) and you want to run a sensu
check for all the IPs and only let it succeed if all of them run.


### Usage

  * Example:
    * `ip_wrapper -host example.org -- check-http.rb -u 'http://%%IP%%/'`
    * Execute `check-http.rb` for every ip that gets resolved from `example.org`
  * Help
    * `ip_wrapper -h`

### License

Copyright (C) 2017 wywy GmbH

Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
