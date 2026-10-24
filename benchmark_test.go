package gocontainer_test

import (
	"strconv"
	"testing"

	"github.com/vardius/gocontainer/v2"
)

// The benchmarks read from a container of 32 services, about what an app wires at startup,
// and the key read is one of them, not the first registered.

var sink *service

func BenchmarkGet(b *testing.B) {
	c, keys := fill(32)
	k := keys[16]
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = k.Get(c)
	}
}

func BenchmarkGetNotRegistered(b *testing.B) {
	c, _ := fill(32)
	k := gocontainer.NewKey[*service]("missing")
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = k.Get(c)
	}
}

func BenchmarkMustGet(b *testing.B) {
	c, keys := fill(32)
	k := keys[16]
	b.ReportAllocs()
	for b.Loop() {
		sink = k.MustGet(c)
	}
}

func BenchmarkHas(b *testing.B) {
	c, keys := fill(32)
	k := keys[16]
	var has bool
	b.ReportAllocs()
	for b.Loop() {
		has = k.Has(c)
	}
	_ = has
}

func BenchmarkGetParallel(b *testing.B) {
	c, keys := fill(32)
	k := keys[16]
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var s *service
		for pb.Next() {
			s, _ = k.Get(c)
		}
		_ = s
	})
}

// BenchmarkRegister replaces a service in a container of n: Register copies the registry, so it grows with n.
func BenchmarkRegister(b *testing.B) {
	for _, n := range []int{1, 32, 1024} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			c, keys := fill(n)
			k, db := keys[n/2], &service{}
			b.ReportAllocs()
			for b.Loop() {
				k.Register(c, db)
			}
		})
	}
}
