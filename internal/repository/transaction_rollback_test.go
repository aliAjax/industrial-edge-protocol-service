package repository

import "testing"

func TestReplayFailureDoesNotCommitTransaction(t *testing.T) {
	transaction := NewTransaction()
	transaction.Add(func() {})
	transaction.Rollback()
	if len(transaction.actions) != 0 {
		t.Fatal("rollback retained queued transaction actions")
	}
	transaction.Add(func() {})
	if len(transaction.actions) != 0 {
		t.Fatal("rollback accepted a new transaction action")
	}
	committed := false
	transaction = &Transaction{rolledBack: true, actions: []func(){func() { committed = true }}}
	transaction.Commit()
	if committed {
		t.Fatal("rolled back transaction executed its actions")
	}
}
