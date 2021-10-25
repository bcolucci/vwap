package util

import "encoding/json"

// CapQueue is a simple sliding queue (FIFO)
type CapQueue struct {
	elems []interface{}
	cap   int
}

func NewCapQueue(cap int) *CapQueue {
	return &CapQueue{
		elems: []interface{}{},
		cap:   cap,
	}
}

// Append appends an element and remove the oldest one if necessary
func (q *CapQueue) Append(elem interface{}) {
	if len(q.elems) == q.cap {
		q.elems = q.elems[1:]
	}
	q.elems = append(q.elems, elem)
}

// CopyTo copies all the values into a custom slice (with a type, it's better ^_^)
func (q *CapQueue) CopyTo(dst interface{}) error {
	b, _ := json.Marshal(q.elems)
	return json.Unmarshal(b, dst)
}
