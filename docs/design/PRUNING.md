# Pruning

In order to manage the size of data stored on disk, topics are
configurable with a retention period. Any messages that fall outside this
period can then be safely removed from the store.

The same is also true for consumer data. It may be the case that one-off
consumers are created for testing or once-only jobs that read from a topic.
After a certain interval, the index of a consumer within a topic can also
be safely removed to release storage on disk.

To achieve this, operators are able to specify retention periods for both
messages and consumer indexes when creating a topic. It is also possible for
a period to be infinite, so that a topic can act as a single source-of-truth
for the state of a resource.

While the server is running. It routinely checks the configured retention
periods on a topic for both messages and consumer indexes. Each message is
stored alongside a timestamp of when it was received by the server, each
consumer also has a timestamp of when it last updated its index within the
topic. The pruner starts at the earliest message index, deleting messages it
finds that have a timestamp that falls outside the retention period. Once
it reaches the first message that is not outside the retention period, it
stops pruning further messages.

The same action is performed for consumer indexes. For each topic, the
consumer retention period is checked and any indexes that have not been
updated within that retention period are removed.

The default interval for pruning checks is every minute. This can be configured
via the `--prune-interval` command-line flag and expects the string
representation of a go `time.Duration` type. The [documentation](https://pkg.go.dev/time#Duration.String)
for the `time.Duration` type explains the format:

> String returns a string representing the duration in the form "72h3m0.5s".
> Leading zero units are omitted. As a special case, durations less than one
> second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
> that the leading digit is non-zero. The zero duration formats as 0s.

The current implementation of storage uses [boltdb](https://github.com/etcd-io/bbolt).
Because of this, the disk space may not appear to be freed. Rather, that space
is made available for the server to store data in within the file. Subsequent
writes will not take up additional space on disk. It may be necessary in the
future to provide operators with a means of compacting the boltdb database.
This is currently possible using the boltdb command-line but the server does
not currently perform compaction.

The implementation for this can be found in the `internal/prune` package within
the source code.
