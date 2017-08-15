package main

import (
	"fmt"

	gocb "gopkg.in/couchbase/gocb.v1"
)

// ChangeSet represents a set of changes for DB.
type ChangeSet struct {
	ID      string
	Changes []Change
}

// Change represents a change for DB.
type Change struct {
	ID      string
	Author  string
	Execute func(b *gocb.Bucket, m *gocb.BucketManager) error
}

// UID returns unique ID for the change based on its content.
func (change Change) UID() string {
	return fmt.Sprintf("%s::%s", change.ID, change.Author)
}
