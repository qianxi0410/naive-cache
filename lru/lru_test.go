package lru_test

import (
	"github.com/qianxi0410/naive-lru/lru"
	"testing"
	"unsafe"
)

type Value string

func (v Value) Len() int {
	return int(unsafe.Sizeof(v))
}

func TestLruGet(t *testing.T) {
	cache := lru.New(100, nil)
	cache.Add("k1", Value("v1"))
	cache.Add("k2", Value("v2"))

	if v, ok := cache.Get("k1"); !ok || string(v.(Value)) != "v1" {
		t.Fail()
	}

	if v, ok := cache.Get("k2"); !ok || string(v.(Value)) != "v2" {
		t.Fail()
	}
}

func TestLruOverSize(t *testing.T) {
	cache := lru.New(16, nil)
	// remove
	cache.Add("k1", Value("v1"))
	cache.Add("k2", Value("v2"))
	cache.Add("k3", Value("v3"))

	if v, ok := cache.Get("k1"); ok || string(v.(Value)) == "v1" {
		t.Fail()
	}

	if cache.Len() != 2 {
		t.Fail()
	}
}