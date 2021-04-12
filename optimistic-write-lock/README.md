## Optimistic Write Lock

This example illustrates use of optimistic write locks to implement a consistent finite state machine.

The example starts an HTTP server which performs queries on GET and proposes updates on PUT.

Any proposed update with an invalid version will be rejected.

Clients must read the version of the key and supply it during update in order to modify the value.

```go
> make run *.go
```

```
> curl -X PUT "http://localhost:8001/testkey?val=testvalue"
{"key":"/testkey","ver":6,"val":"testvalue"}

> curl -X PUT "http://localhost:8001/testkey?val=testvalue2"
Version mismatch (0 != 6)

> curl -X PUT "http://localhost:8001/testkey?val=testvalue2&ver=6"
{"key":"/testkey","ver":8,"val":"testvalue2"}

> curl -X PUT "http://localhost:8001/testkey?val=testvalue3&ver=6"
Version mismatch (6 != 8)

> curl -X PUT "http://localhost:8001/testkey?val=testvalue3&ver=8"
{"key":"/testkey","ver":10,"val":"testvalue3"}
```

Optimistic write locks can be used to implement [CP](https://en.wikipedia.org/wiki/CAP_theorem)
systems using dragonboat.

While this example demonstrates key-level write locks, similar write locks can be implemented
at any level of granularity within a finite state machine to linearize anything from individual
keys to entire datasets.
