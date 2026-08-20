package repository

import "sync"

type Transaction struct {
	mu         sync.Mutex
	committed  bool
	rolledBack bool
	actions    []func()
}

func NewTransaction() *Transaction { return &Transaction{} }
func (t *Transaction) Add(action func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.committed || t.rolledBack {
		return
	}
	t.actions = append(t.actions, action)
}
func (t *Transaction) Commit() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.committed || t.rolledBack {
		return
	}
	for _, a := range t.actions {
		a()
	}
	t.committed = true
}
func (t *Transaction) Rollback() { t.mu.Lock(); t.rolledBack = true; t.actions = nil; t.mu.Unlock() }
func (t *Transaction) Done() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.committed || t.rolledBack
}
