package cbcm

import "errors"

var (
	// ErrNotAllowedChangesetID error returnes by runner.Execute if changeset.ID has not allowed value.
	ErrNotAllowedChangesetID = errors.New("changeset.ID is not allowed")

	// ErrNilChangesetExecute error returnes by runner.Execute if changeset.Execute has a nil value.
	ErrNilChangesetExecute = errors.New("changeset.Execute is not allowed")

	// ErrNotUniqueChangeSets error returnes by runner.Execute if changes aren't unique (id).
	ErrNotUniqueChangeSets = errors.New("changesets are not unique")
)
