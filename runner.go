package cbcm

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

// Execute executes changesets on a bucket.
func (r *Runner) Execute(bucket, password string, changes []ChangeSet) error {
	b, err := r.cluster.OpenBucket(bucket, password)
	if err != nil {
		return err
	}
	defer b.Close()

	m := b.Manager(bucket, password)
	return execute(b, m, changes)
}

func execute(b *gocb.Bucket, m *gocb.BucketManager, changes []ChangeSet) error {
	if err := validate(changes); err != nil {
		return err
	}

	changelog := changeLogDocument{}
	if _, err := b.Get(changelogDocKey, &changelog); err != gocb.ErrKeyNotFound {
		return err
	}
	defer b.Upsert(changelogDocKey, &changelog, 0)

	for _, changeset := range changes {
		if _, ok := changelog[changeset.ID]; ok {
			log.Printf("Skipping changeset: \"%s\"", changeset)
		} else {
			log.Printf("Executing changeset \"%s\"", changeset)
			err := changeset.Execute(b, m)
			if err != nil {
				log.Printf("Changeset \"%s\" execution has failed. %s", changeset, err)
				return err
			}
			info := changeSetInfo{
				ID:        changeset.ID,
				Author:    changeset.Author,
				AppliedAt: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
			}
			changelog[changeset.ID] = info
			log.Printf("Changeset \"%s\" execution has finished", changeset)
		}
	}
	return nil
}

func validate(changes []ChangeSet) error {
	m := make(map[string]bool)
	for _, changeset := range changes {
		if changeset.ID == "" {
			log.Printf("ID is not allowed: \"%s\"", changeset)
			return ErrNotAllowedChangesetID
		}
		if changeset.Execute == nil {
			log.Printf("Execute is nil: \"%s\"", changeset)
			return ErrNilChangesetExecute
		}

		if _, ok := m[changeset.ID]; ok {
			log.Printf("Duplicate change has found: \"%s\"", changeset.ID)
			return ErrNotUniqueChangeSets
		}
		m[changeset.ID] = true
	}
	return nil
}
