# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# How to use
This is private repository so in order to be able to include it in your microservice or Go application the following steps needs to be carried out:

- Create a `~/.netrc` file on your workstation (if you don't have one).
- Add an entry similar to the following:

```bash
machine github.com
login {YOUR_GITHUB_USER_NAME}
password {YOUR_PERSONAL_ACCESS_TOKEN_FROM_GITHUB}
```

Update the `GOPRIVATE` variable for your Go environment to allow the Go to look for private repo in the location as a source and not use a proxy.

```bash
go env -w GOPRIVATE=github.com/eyewa/eyewa-go-lib
```

Now you can include the lib your app:

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