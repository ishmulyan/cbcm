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

// changeLogDocument represents a document for recording history of applied changesets.
type changeLogDocument map[string]changeSetInfo

// changeSetInfo represents an information about applied changeset.
type changeSetInfo struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	AppliedAt uint64 `json:"appliedAt"`
}
