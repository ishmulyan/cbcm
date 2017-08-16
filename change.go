package cbcm

import (
	"fmt"

	gocb "gopkg.in/couchbase/gocb.v1"
)

// ChangeSet represents a change for DB.
type ChangeSet struct {
	ID      string
	Author  string
	Execute func(b *gocb.Bucket, m *gocb.BucketManager) error
}

func (changeset ChangeSet) String() string {
	return fmt.Sprintf("id: %s, author: %s", changeset.ID, changeset.Author)
}

// ChangeLogDocument represents a document for recording history of applied changesets.
type ChangeLogDocument map[string]ChangeSetInfo

// ChangeSetInfo represents an information about applied changeset.
type ChangeSetInfo struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	AppliedAt uint64 `json:"appliedAt"`
}
