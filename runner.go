package main

import (
	"log"
	"time"

	gocb "gopkg.in/couchbase/gocb.v1"
)

const changelogDocKey = "dbchangelog"

// Runner is a runner of changesets on couchbase cluster/bucket.
type Runner struct {
	cluster *gocb.Cluster
}

// New instantiates a new isntance of Runner. Returns an error if connSpecStr is not valid.
func New(connSpecStr string) (*Runner, error) {
	cluster, err := gocb.Connect(connSpecStr)
	if err != nil {
		return nil, err
	}

	r := Runner{
		cluster: cluster,
	}
	return &r, nil
}

// Execute executes changeset on a bucket.
func (r *Runner) Execute(bucket, password string, changeset ChangeSet) error {
	b, err := r.cluster.OpenBucket(bucket, password)
	if err != nil {
		return err
	}
	defer b.Close()

	m := b.Manager(bucket, password)
	return execute(b, m, changeset)
}

func execute(b *gocb.Bucket, m *gocb.BucketManager, changeset ChangeSet) error {
	changelogDoc := ChangeLogDocument{}
	b.Get(changelogDocKey, &changelogDoc)
	changeSetDoc, ok := changelogDoc[changeset.ID]
	if !ok {
		changeSetDoc = ChangeSetDocument{}
		changelogDoc[changeset.ID] = changeSetDoc
	}

	log.Printf("Executing changeset \"%s\"...", changeset.ID)
	for _, change := range changeset.Changes {
		if _, ok = changeSetDoc[change.UID()]; ok {
			log.Printf("Skipping change \"%s\" from \"%s\"...", change.ID, changeset.ID)
		} else {
			log.Printf("Executing change \"%s\" from \"%s\"...", change.ID, changeset.ID)
			err := change.Execute(b, m)
			if err != nil {
				log.Printf("Executing change \"%s\" from \"%s\" has failed. %s", change.ID, changeset.ID, err)
				return err
			}
			changeSetDoc[change.UID()] = time.Now().UnixNano()
			log.Printf("Change \"%s\" execution from \"%s\" has finished", change.ID, changeset.ID)
		}
	}
	b.Upsert(changelogDocKey, &changelogDoc, 0)
	log.Printf("Changeset \"%s\" execution has finished", changeset.ID)
	return nil
}
