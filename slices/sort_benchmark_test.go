// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
)

// These benchmarks compare sorting a large slice of int with sort.Ints vs.
// slices.Sort
func makeRandomInts(n int) []int {
	rand.Seed(42)
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = rand.Intn(n)
	}
	return ints
}

func makeSortedInts(n int) []int {
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = i
	}
	return ints
}

func makeReversedInts(n int) []int {
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = n - i
	}
	return ints
}

const N = 100_000

func BenchmarkSortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		sort.Ints(ints)
	}
}

func BenchmarkBubbleSortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		BubbleSort(ints)
	}
}

func BenchmarkSlicesSortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		Sort(ints)
	}
}

func BenchmarkSlicesSortInts_Sorted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeSortedInts(N)
		b.StartTimer()
		Sort(ints)
	}
}

func BenchmarkSlicesSortInts_Reversed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeReversedInts(N)
		b.StartTimer()
		Sort(ints)
	}
}

// Since we're benchmarking these sorts against each other, make sure that they
// generate similar results.
func TestIntSorts(t *testing.T) {
	ints := makeRandomInts(200)
	ints2 := Clone(ints)

	sort.Ints(ints)
	Sort(ints2)

	for i := range ints {
		if ints[i] != ints2[i] {
			t.Fatalf("ints2 mismatch at %d; %d != %d", i, ints[i], ints2[i])
		}
	}
}

// The following is a benchmark for sorting strings.

// makeRandomStrings generates n random strings with alphabetic runes of
// varying lenghts.
func makeRandomStrings(n int) []string {
	rand.Seed(42)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		slen := 2 + rand.Intn(50)
		for j := 0; j < slen; j++ {
			sb.WriteRune(letters[rand.Intn(len(letters))])
		}
		ss[i] = sb.String()
	}
	return ss
}

func TestStringSorts(t *testing.T) {
	ss := makeRandomStrings(200)
	ss2 := Clone(ss)

	sort.Strings(ss)
	Sort(ss2)

	for i := range ss {
		if ss[i] != ss2[i] {
			t.Fatalf("ss2 mismatch at %d; %s != %s", i, ss[i], ss2[i])
		}
	}
}

func BenchmarkSortStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStrings(N)
		b.StartTimer()
		sort.Strings(ss)
	}
}

func BenchmarkSlicesSortStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStrings(N)
		b.StartTimer()
		Sort(ss)
	}
}

// These benchmarks compare sorting a slice of structs with sort.Sort vs.
// slices.SortFunc.
type myStruct struct {
	a, b, c, d string
	n          int
}

type myStructs []*myStruct

func (s myStructs) Len() int           { return len(s) }
func (s myStructs) Less(i, j int) bool { return s[i].n < s[j].n }
func (s myStructs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func makeRandomStructs(n int) myStructs {
	rand.Seed(42)
	structs := make([]*myStruct, n)
	for i := 0; i < n; i++ {
		structs[i] = &myStruct{n: rand.Intn(n)}
	}
	return structs
}

func TestStructSorts(t *testing.T) {
	ss := makeRandomStructs(200)
	ss2 := make([]*myStruct, len(ss))
	for i := range ss {
		ss2[i] = &myStruct{n: ss[i].n}
	}

	sort.Sort(ss)
	SortFunc(ss2, func(a, b *myStruct) bool { return a.n < b.n })

	for i := range ss {
		if *ss[i] != *ss2[i] {
			t.Fatalf("ints2 mismatch at %d; %v != %v", i, *ss[i], *ss2[i])
		}
	}
}

func BenchmarkSortStructs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStructs(N)
		b.StartTimer()
		sort.Sort(ss)
	}
}

func BenchmarkSortFuncStructs(b *testing.B) {
	lessFunc := func(a, b *myStruct) bool { return a.n < b.n }
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStructs(N)
		b.StartTimer()
		SortFunc(ss, lessFunc)
	}
}
