package main

// ChangeLogDocument represents a document for recording history of applied changesets.
type ChangeLogDocument map[string]ChangeSetDocument

// ChangeSetDocument represents a document for recording history of applied changes.
type ChangeSetDocument map[string]int64
