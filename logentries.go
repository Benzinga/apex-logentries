// Package logentries implements a logentries apex log handler.
package logentries

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/apex/log"
)

// ErrorHandling specifies how to deal with errors connecting to and
// delivering logs to LogEntries.
type ErrorHandling int

const (
	// PanicOnError specifies that errors will panic, and connection issues
	// will be logged.
	PanicOnError ErrorHandling = iota
	// LogOnError specifies that errors will be logged to stderr.
	LogOnError
	// IgnoreErrors specifies that errors will not be recorded or handled.
	// Your logs may silently not be delivered.
	IgnoreErrors
)

const (
	// defaultQueueLen specifies the length of the internal queue. If this
	// queue fills up, messages will be discarded.
	defaultQueueLen = 1024
	// dialTimeout specifies how much time is spent waiting for dial to
	// complete.
	dialTimeout = time.Second * 10
	// writeTimeout specifies how much time is spent waiting for a write to
	// complete.
	writeTimeout = time.Second * 10
	// retryDelay specifies how long to wait between reconnections. It is
	// important that this value is not too large.
	retryDelay = time.Second
)

var (
	// defaultAddress contains the default addresses for LogEntries ingestion,
	// keyed by a boolean specifying whether or not TLS is desired.
	defaultAddress = map[bool]string{
		false: "data.logentries.com:80",
		true:  "data.logentries.com:443",
	}
)

// Config holds the configuration options for LogEntries output.
type Config struct {
	// Logentries settings.
	Token     string      // LogEntries token to use.
	UseTLS    bool        // Whether or not to use encryption.
	Address   string      // Address, if you want to override it.
	TLSConfig *tls.Config // TLS configuration to use.

	QueueLen      int           // Length of internal queue.
	Discard       bool          // Specifies if logs are discarded when the queue fills.
	ErrorHandling ErrorHandling // Specifies how errors should be handled.
}

// Handler implementation.
type Handler struct {
	*Config

	pfx    []byte
	ch     chan *log.Entry
	ctx    context.Context
	cancel context.CancelFunc
}

// New handler. This spawns a new Goroutine and connects to logentries
// immediately.
func New(config Config) *Handler {
	if config.Address == "" {
		config.Address = defaultAddress[config.UseTLS]
	}

	if config.QueueLen == 0 {
		config.QueueLen = defaultQueueLen
	}

	// Create context for connectionLoop.
	ctx, cancel := context.WithCancel(context.Background())

	handler := &Handler{
		Config: &config,
		pfx:    []byte(config.Token + " "),
		ch:     make(chan *log.Entry, config.QueueLen),
		ctx:    ctx,
		cancel: cancel,
	}

	go handler.connectionLoop()

	return handler
}

// Close closes the logentries connection and frees associated resources.
// Once this is called, logging with this handler will panic. Most users do
// not need to use this.
func (h *Handler) Close() {
	h.cancel()
}

func (h *Handler) handleError(err error) {
	switch h.ErrorHandling {
	case PanicOnError:
		panic(err)
	case LogOnError:
		fmt.Fprintln(os.Stderr, "apex-logentries:", err)
	default:
		// Ignore.
	}
}

// connectionLoop handles the internal connection that sends the logs to
// LogEntries.
func (h *Handler) connectionLoop() {
	var conn net.Conn
	var err error
	var enc *json.Encoder

	for {
		dialer := net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: writeTimeout,
		}

		if h.UseTLS {
			conn, err = tls.DialWithDialer(&dialer, "tcp", h.Address, h.TLSConfig)
		} else {
			conn, err = dialer.Dial("tcp", h.Address)
		}

		if err != nil {
			h.handleError(err)
			goto Error
		}

		enc = json.NewEncoder(conn)

		for {
			select {
			case log := <-h.ch:
				conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				conn.Write(h.pfx)
				err := enc.Encode(log)
				if err != nil {
					h.handleError(err)
					goto Error
				}
				break

			case <-h.ctx.Done():
				goto Done
			}
		}

	Error:
		select {
		case <-time.After(retryDelay):
			continue
		case <-h.ctx.Done():
			goto Done
		}
	}

Done:
	err = conn.Close()
	if err != nil {
		h.handleError(err)
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	select {
	case h.ch <- e:
		return nil
	default:
		if h.Discard {
			return fmt.Errorf("queue full")
		}

		panic("apex-logentries: queue full")
	}
}
