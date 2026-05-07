package domain

// TaskOp type for redis operations commands.
type TaskOp string

const (
	// OpSet is a set command for redis.
	OpSet TaskOp = "set"
	// OpDel is a delete command for redis.
	OpDel TaskOp = "del"
)

// NotifyStatusTask represents a task to update notification status in a cache or store.
type NotifyStatusTask struct {
	Op    TaskOp
	Key   string
	Value string
}
