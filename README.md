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
bazel run :kurapika_container_image
docker run --rm -it -p8081:8081 bazel/kurapika_container_image
```

## Publish

### Docker
```
bazel run :kurapika_container_push
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
