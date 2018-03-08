# Apex LogEntries
This is an alternative implementation of LogEntries for Apex. It has no external dependencies (other than Apex itself) and uses channels as queues for efficient operation. It attempts to be durable and offer a variety of options for handling errors.

# Getting Started
`apex-logentries` is a handler for the Apex logger. To make use of it, you need to install it and use it with Apex.

## Prerequisites
`apex-logentries` is built in the Go programming language. If you are new to Go, you will need to [install Go](https://golang.org/dl/). This is the only dependency.

## Acquiring
Next, you'll want to `go get` apex-logentries, like so:

```sh
go get github.com/Benzinga/apex-logentries
```

If your `$GOPATH` is configured, and git is setup to know your credentials, in a few moments the command should complete with no output. The repository will exist under `$GOPATH/src/github.com/Benzinga/apex-logentries`. It cannot be moved from this location.

Hint: If you've never used Go before, your `$GOPATH` will be under the `go` folder of your user directory.

## Demo
A documented demo is included in the repository as `demo.go`. It is pretty simple to invoke. In a shell inside the project folder, do:

```
go run logentries -token [logentries token here]
```

If it worked, you'll see a few messages in the log pointed to by the token.

# Usage
To use, call `logentries.New` with your desired `Config`, then use Apex's `SetHandler` function.

```go
le := logentries.New(logentries.Config{
    UseTLS: true,
    Token: "token",
})

log.SetHandler(le)
```

## Configuration
The following configuration options are available:

- `Token` (`string`): Specifies the LogEntries token to use.
- `UseTLS` (`bool`): Whether or not to use TLS to connect to LogEntries.
- `Address` (`string`): Address to use to connect to LogEntries.
- `TLSConfig` (`*tls.Config`): TLS configuation to use, if using TLS.
- `QueueLen` (`int`): Length of queue to use for queueing messages to LogEntries.
- `Discard` (`bool`) Whether or not it is OK to silently discard log messages if the queue fills (for example, if there is very high log volume, or the connection is unstable.)
- `ErrorHandling`: Determines how to handle errors that occur on the connection and otherwise. The following values are possible:
    - `PanicOnError`: Every error will panic. This is the default.
    - `LogOnError`: Every error will be logged to stderr.
    - `IgnoreErrors`: Errors will be silently ignored.