package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapQueue(t *testing.T) {
	assert := assert.New(t)

	var arr []int

	q := NewCapQueue(2)
	q.CopyTo(&arr)

	assert.Equal([]int{}, arr)

	q.Append(1)
	q.CopyTo(&arr)
	assert.Equal([]int{1}, arr)

	q.Append(2)
	q.CopyTo(&arr)
	assert.Equal([]int{1, 2}, arr)

	q.Append(3)
	q.CopyTo(&arr)
	assert.Equal([]int{2, 3}, arr)
}
