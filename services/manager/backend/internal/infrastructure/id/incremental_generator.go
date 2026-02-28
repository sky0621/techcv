package id

import (
	"fmt"
	"sync/atomic"
)

type IncrementalGenerator struct {
	counter atomic.Uint64
}

func NewIncrementalGenerator() *IncrementalGenerator {
	return &IncrementalGenerator{}
}

func (g *IncrementalGenerator) NewID() string {
	n := g.counter.Add(1)
	return fmt.Sprintf("profile_%d", n)
}
