# golib
> An extremely opinionated set of modules

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/thisiserico/golib?tab=packages)


## 🧐 Motivation

When kicking off a new project, often times engineers decide not to care about consistency, debuggability and other little details that are not completely necessary to "deliver".
This set of modules provide exactly that: a way to keep things consistent without getting in your way.

At the same time, it provides an opinionated way on how certain elements should look like. Examples are the [`oops`][oops] or [`logger`][logger] packages,
which expose a simpler interface from what we're used to.


## 👩‍💻 Provided modules

The [`halt`][halt] package lets you handle graceful shutdowns.

The [`kv`][kv] package lets you define key-value pairs to be used in multiple situations.

The [`logger`][logger] package lets you log as you'd normally do, only a simplified contract is used.

The [`o11y`][o11y] package contains functionality that [`opentelemetry`][opentelemetry] uses to ingest telemetry data.

The [`oops`][oops] package lets you create contextual errors using a simplified contract.

The [`pubsub`][pubsub] package lets you publish and subscribe to messages.


## 🥺 What's next

Existing packages are subject to change.
[Semantic versioning][semver] is used, backwards compatibility will be kept.
Different concrete implementations or packages will be added when needed.


[opentelemetry]: https://pkg.go.dev/go.opentelemetry.io
[halt]: https://pkg.go.dev/github.com/thisiserico/golib/halt
[kv]: https://pkg.go.dev/github.com/thisiserico/golib/kv
[logger]: https://pkg.go.dev/github.com/thisiserico/golib/logger
[o11y]: https://pkg.go.dev/github.com/thisiserico/golib/o11y
[oops]: https://pkg.go.dev/github.com/thisiserico/golib/oops
[pubsub]: https://pkg.go.dev/github.com/thisiserico/golib/pubsub
[semver]: https://semver.org

