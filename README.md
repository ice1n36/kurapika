# kurapika

central intel

[![Go Report Card](https://goreportcard.com/badge/github.com/ice1n36/kurapika)](https://goreportcard.com/report/github.com/ice1n36/kurapika)
[![Coverage Status](https://img.shields.io/codecov/c/github/ice1n36/kurapika.svg)](https://codecov.io/gh/ice1n36/kurapika)

## Install

```
go get github.com/ice1n36/kurapika
```

## Build & Run

### Locally
```
bazel run :kurapika
```

### Docker
```
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 :kurapika_container_image
docker run --rm -v /Users/tngo/go/src/github.com/ice1n36/kurapika/config:/config -it -p8081:8081 bazel/kurapika_container_image
```

## Publish

### Docker

```
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 :kurapika_container_image_push
```

## Test

### Unit
TODO

### Locally
```
curl -X POST localhost:8081/new_app -d '{"app_id": "cz.digerati.iqtest", "device_codename": "walleye", "os": "android"}'
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
