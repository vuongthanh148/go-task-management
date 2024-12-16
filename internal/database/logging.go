package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/kwanpham2195/go-gcp-boilerplate/internal/derrors"
	"github.com/kwanpham2195/go-gcp-boilerplate/internal/log"
)

// QueryLoggingDisabled stops logging of queries when true.
// For use in tests only: not concurrency-safe.
var QueryLoggingDisabled bool

var queryCounter int64 // atomic: per-process counter for unique query IDs

type queryEndLogEntry struct {
	ID              string
	Query           string
	Args            string
	DurationSeconds float64
	Error           string `json:",omitempty"`
}

func logQuery(ctx context.Context, query string, args []any, instanceID string, retryable bool) func(*error) {
	if QueryLoggingDisabled {
		return func(*error) {}
	}
	const maxlen = 300 // maximum length of displayed query

	// To make the query more compact and readable, replace newlines with spaces
	// and collapse adjacent whitespace.
	var r []rune
	for _, c := range query {
		if c == '\n' {
			c = ' '
		}
		if len(r) == 0 || !unicode.IsSpace(r[len(r)-1]) || !unicode.IsSpace(c) {
			r = append(r, c)
		}
	}
	query = string(r)
	if len(query) > maxlen {
		query = query[:maxlen] + "..."
	}

	uid := generateLoggingID(instanceID)

	// Construct a short string of the args.
	const (
		maxArgs   = 20
		maxArgLen = 50
	)
	var argStrings []string
	for i := 0; i < len(args) && i < maxArgs; i++ {
		s := fmt.Sprint(args[i])
		if len(s) > maxArgLen {
			s = s[:maxArgLen] + "..."
		}
		argStrings = append(argStrings, s)
	}
	if len(args) > maxArgs {
		argStrings = append(argStrings, "...")
	}
	argString := strings.Join(argStrings, ", ")

	log.Debugf(ctx, "%s %s args=%s", uid, query, argString)
	start := time.Now()
	return func(errp *error) {
		dur := time.Since(start)
		if errp == nil { // happens with queryRow
			log.Debugf(ctx, "%s done", uid)
		} else {
			derrors.Wrap(errp, "DB running query %s", uid)
			entry := queryEndLogEntry{
				ID:              uid,
				Query:           query,
				Args:            argString,
				DurationSeconds: dur.Seconds(),
			}
			if *errp == nil {
				log.Debug(ctx, entry)
			} else {
				entry.Error = (*errp).Error()
				logf := log.Error
				if errors.Is(ctx.Err(), context.Canceled) ||
					strings.Contains(entry.Error, "pq: canceling statement due to user request") {
					logf = log.Debug
				}
				// If the transaction is retryable and this is a serialization error,
				// then it's not really an error at all. Log it as debug, so if
				// we get a "failed due to max retries" error, we can find
				// these easily. However, these errors can also be noisy, so we
				// can also hide them by setting GO_DISCOVERY_LOG_LEVEL=info.
				if retryable && isSerializationFailure(*errp) {
					logf = log.Debug
				}
				logf(ctx, entry)
			}
		}
	}
}

func (db *DB) logTransaction(ctx context.Context) func(*error) {
	if QueryLoggingDisabled {
		return func(*error) {}
	}
	uid := generateLoggingID(db.instanceID)
	log.Debugf(ctx, "%s transaction (isolation %s) started", uid, db.opts.Isolation)
	start := time.Now()
	return func(errp *error) {
		log.Debugf(ctx, "%s transaction (isolation %s) finished in %s with error %v",
			uid, db.opts.Isolation, time.Since(start), *errp)
	}
}

func generateLoggingID(instanceID string) string {
	if instanceID == "" {
		instanceID = "local"
	} else if len(instanceID) > 8 {
		// App Engine instance IDs are long strings. The low-order part seems
		// quite random, so shortening the ID will still likely result in
		// something unique.
		instanceID = instanceID[len(instanceID)-4:]
	}
	n := atomic.AddInt64(&queryCounter, 1)
	return fmt.Sprintf("%s-%d", instanceID, n)
}
