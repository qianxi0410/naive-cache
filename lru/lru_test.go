package lru_test

import (
	"github.com/qianxi0410/naive-lru/lru"
	"reflect"
	"testing"
)

type Value string

func (v Value) Len() int {
	return len(v)
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
	cache := lru.New(10, nil)
	// remove
	cache.Add("k1", Value("v1"))
	cache.Add("k2", Value("v2"))
	cache.Add("k3", Value("v3"))

	if _, ok := cache.Get("k1"); ok {
		t.Fail()
	}

	if cache.Len() != 2 {
		t.Fail()
	}
}

func  TestLruOnEvicted(t *testing.T)  {
	keys := make([]string, 0)
	cache := lru.New(int64(10), func(k string, value lru.Value) {
		keys = append(keys, k)
	})
	cache.Add("key1", Value("123456"))
	cache.Add("k2", Value("k2"))
	cache.Add("k3", Value("k3"))
	cache.Add("k4", Value("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}