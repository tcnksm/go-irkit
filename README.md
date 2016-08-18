# go-irkit 

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[license]: https://github.com/tcnksm/go-irkit/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/tcnksm/go-irkit/v1

`go-irkit` is the unofficial golang client for [IRKit](http://getirkit.com/en/) (IRKit is a Wi-Fi enabled Open Source Infrared Remote Controller device). Note that API is not completed (and currently it's only support internet HTTP API). It has only what [I](https://github.com/tcnksm) need. If you want to use other feature, PR is always welcome. 

To install, use `go get` (it uses [`context`](https://golang.org/pkg/context/) pacakge, so you need Go1.7 or later),

```bash
$ go get github.com/tcnksm/go-irkit/v1
```

Full documentation is available at http://godoc.org/github.com/tcnksm/go-irkit/v1 . See example usage on [`v1/_example`](v1/_example) directory.
