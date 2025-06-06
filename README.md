# Goava

My own personal Go toolset. Something like Google's Guava, but at a small scale and in Go.

* Package `cache`: Cache implementations.
* Package `casing`: String casing conversion tools. For example, from `CamelCase` or `DromedaryCase` to `dot.case` or `snake_case`.
* Package `err`: Error handling tools.
* Package `maps`: Esoteric data structures based on maps.
* Package `msg`: Messaging tools.
* Package `pipe`: Tools for pipeline creation and comunication.
* Package `rate`: Tools for rate calculation and limiting.
* Package `svc`: Tools for service management.
* Package `test`: Some tools used in testing.
  -`eventually.go`: repeats a test until it finally succeeds (then the test is marked as succeeded),
    or until it fails (then the test is marked as failed).
  -`ports.go`: find a free port for spawning test services.
  -`chan.go`: test tools for channels.

