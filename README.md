README
======

Use [Hetzner DNS API](https://dns.hetzner.com/api-docs/) from [Golang](https://golang.org/).

## Install

```console
go get github.com/panta/go-hetzner-dns
```

## Usage

### Authentication

Hetzner DNS API uses an API token. You can either explicitly set the
token in the `ApiKey` field of the `Client` object, or you can set the
environment variable `HETZNER_API_KEY`.

### API

Almost all Hetzner DNS APIs are supported. See the example in `cmd/example` or
the tests in `hetzner_dns_test.go` for a more complete list.

```go
import "github.com/panta/go-hetzner-dns"

client := hetzner_dns.Client{}

// ...

zonesResponse, err := client.GetZones(context.Background(), "", "", 1, 100)
if err != nil {
    log.Fatal(err)
}

// ...

recordsResponse, err := client.GetRecords(context.Background(), "zone-id", 0, 0)
if err != nil {
    log.Fatal(err)
}
```

### Overriding the http.Client

It's possible to use a custom `http.Client` object by setting the
`HttpClient` field in the `Client` object. If not set, the library
will create one for you.

### Example program

To build the example program on a unix-like:

```shell
$ make
```

then:

```shell
$ export HETZNER_API_KEY="....."
$ ./go-hetzner-dns list
$ ./go-hetzner-dns add-record -zone=ZONEID RECORD_NAME TYPE RECORD_VALUE
```

## Author

By [Marco Pantaleoni](https://github.com/panta).

## LICENSE

```
Copyright 2021 Marco Pantaleoni.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
