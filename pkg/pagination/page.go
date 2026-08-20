package pagination

import "strconv"

type Request struct {
	Limit  int
	Offset int
	Cursor string
}
type Response[T any] struct {
	Items      []T    `json:"items"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	NextCursor string `json:"next_cursor,omitempty"`
	Total      int    `json:"total"`
}

func Parse(limit, offset string) Request {
	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)
	if l <= 0 || l > 500 {
		l = 100
	}
	if o < 0 {
		o = 0
	}
	return Request{Limit: l, Offset: o}
}
func Slice[T any](items []T, r Request) Response[T] {
	if r.Offset > len(items) {
		r.Offset = len(items)
	}
	end := r.Offset + r.Limit
	if end > len(items) {
		end = len(items)
	}
	next := ""
	if end < len(items) {
		next = strconv.Itoa(end)
	}
	return Response[T]{Items: items[r.Offset:end], Limit: r.Limit, Offset: r.Offset, NextCursor: next, Total: len(items)}
}
func WithCursor[T any](items []T, cursor string, limit int) Response[T] {
	offset, _ := strconv.Atoi(cursor)
	return Slice(items, Request{Offset: offset, Limit: limit})
}
