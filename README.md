[![Tests & Linting](https://github.com/Eyewa/eyewa-go-lib/actions/workflows/branches.yaml/badge.svg)](https://github.com/Eyewa/eyewa-go-lib/actions/workflows/branches.yaml)

# eyewa-go-lib
Shared Go Lib for Eyewa's microservices

## List of capabilities
- Packages:
  - Produce/consume events to/from RabbitMQ
  - Metrics instrumentation using OpenTelemetry
  - Tracing instrumentation using OpenTelemetry 
  - Logging
  - Generate  uuids in UUID.v4 format
  - Exposing profiling endpoint using pprof
  - Database drivers
# How to use
This is a private repository, so in order to include it in a microservice or application the following steps need to be carried out:

- Create a `~/.netrc` file on your workstation (if you don't have one).
- Add an entry similar to the following:

```bash
machine github.com
login {YOUR_GITHUB_USER_NAME}
password {YOUR_PERSONAL_ACCESS_TOKEN_FROM_GITHUB}
```

Update the `GOPRIVATE` variable in your Go environment so the private repo can be located as the module's source and not use a proxy.

```bash
go env -w GOPRIVATE=github.com/eyewa/eyewa-go-lib
```

Now you can include the lib into your app:

```bash
go get github.com/eyewa/eyewa-go-lib@dev // pulls the dev branch
```

```bash
go get github.com/eyewa/eyewa-go-lib // pulls the lastest merg
```

Updating dependencies for the lib

```bash
go get -u all
```
