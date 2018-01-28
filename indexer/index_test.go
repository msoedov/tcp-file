package indexer

import (
	"sync"
	"testing"

	. "github.com/franela/goblin"
)

func TestNewIndexer(t *testing.T) {
	idx, err := BuildIndex("const.go")
	if err != nil {
		t.Fatalf("Failed to create index %v\n", err)
	}
	g := Goblin(t)
	g.Describe("Indexer", func() {
		g.It("Should return the first line", func() {
			lineStr, _ := idx.GetLine(int64(1))
			g.Assert(lineStr).Equal("package indexer\n")
		})
		g.It("Should handle neg indexes", func() {
			lineStr, err := idx.GetLine(int64(-3))
			g.Assert(lineStr).Equal("")
			g.Assert(err).Equal(ouchErr)
		})
		g.It("Should handle out of boundaries line", func() {
			lineStr, err := idx.GetLine(int64(30000))
			g.Assert(lineStr).Equal("")
			g.Assert(err).Equal(offsetErr)
		})
	})
}

func BenchmarkTest(b *testing.B) {
	idx, err := BuildIndex("const.go")
	if err != nil {
		b.Fatalf("Failed to create index %v\n", err)
	}
	for n := 0; n < b.N; n++ {
		idx.GetLine(int64(n % 10))
	}
}

func BenchmarkConcurency(b *testing.B) {
	idx, err := BuildIndex("const.go")
	if err != nil {
		b.Fatalf("Failed to create index %v\n", err)
	}
	var wg sync.WaitGroup
	for thread := 1; thread <= 10; thread++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				idx.GetLine(int64(n % 10))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
