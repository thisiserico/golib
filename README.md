# golib
> An extremely opinionated set of modules

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/thisiserico/golib/v2?tab=packages)


## 🧐 Motivation
When kicking off a new project, often times engineers decide not to care about consistency, debuggability and other little details that are not completely necessary to "deliver".
This set of modules provide exactly that: a way to keep things consistent without getting in your way.

At the same time, it provides an opinionated way on how certain elements should look like. Examples are the [`errors`][errors] or [`logger`][logger] packages,
which expose a simpler interface from what we're used to.


## 👩‍💻 Provided modules
[`github.com/thisiserico/golib/v2/cntxt`][cntxt]

The `cntxt` package lets you interact with known attributes that are required to keep in a context.

[`github.com/thisiserico/golib/v2/errors`][errors]

The `errors` package lets you create contextual errors using a simplified contract.

[`github.com/thisiserico/golib/v2/halt`][halt]

The `halt` package lets you handle graceful shutdowns.

[`github.com/thisiserico/golib/v2/kv`][kv]

The `kv` package lets you define key-value pairs to be used in multiple situations.

[`github.com/thisiserico/golib/v2/logger`][logger]

The `logger` package lets you log as you'd normally do, only a simplified contract is used.

[`github.com/thisiserico/golib/v2/pubsub`][pubsub]

The `pubsub` package lets you publish and subscribe to messages.

[`github.com/thisiserico/golib/v2/trace`][trace]

The `trace` package lets you trace operations, wrapping `opentracing` underneath.


## 🥺 Missing packages
Existing packages are subject to change.
[Semantic versioning][semver] is used, backwards compatibility will be kept.


[cntxt]: tree/master/cntxt
[errors]: tree/master/errors
[halt]: tree/master/halt
[kv]: tree/master/kv
[logger]: tree/master/logger
[pubsub]: tree/master/pubsub
[trace]: tree/master/trace
[semver]: https://semver.org

