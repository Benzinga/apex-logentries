# Apex LogEntries
This is an alternative implementation of LogEntries for Apex. It has no external dependencies (other than Apex itself) and uses channels as queues for efficient operation. It attempts to be durable and offer a variety of options for handling errors.

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

## Usage
To use, simply call `logentries.New` with your desired `Config`, then use Apex's `SetHandler` function.

```go
le := logentries.New(logentries.Config{
    UseTLS: true,
    Token: conf.LEToken,
})
log.SetHandler(le)
```
