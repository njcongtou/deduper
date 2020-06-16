package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lrucache := New(int64(0), nil)
	lrucache.Add("key1", String("1234"))
	if v, ok := lrucache.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("Test Get key1=1234 failed.")
	}
	if _, ok := lrucache.Get("key2"); ok {
		t.Fatalf("Test Get cache miss key2 failed.")
	}
}

func TestDeleteOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2) // set maxBytes to 2 items
	lrucache := New(int64(cap), nil)
	lrucache.Add(k1, String(v1))
	lrucache.Add(k2, String(v2))
	lrucache.Add(k3, String(v3))

	if _, ok := lrucache.Get(k2); !ok || lrucache.Len() != 2 {
		t.Fatalf("Test DeleteOldest item key2 failed.")
	}

	if _, ok := lrucache.Get(k3); !ok || lrucache.Len() != 2 {
		t.Fatalf("Test DeleteOldest item key3 failed.")
	}

	if _, ok := lrucache.Get(k1); ok || lrucache.Len() != 2 {
		t.Fatalf("Test DeleteOldest item key1 failed.")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callbackFunc := func(key string, value Value) {
		keys = append(keys, key)
	}
	lrucache := New(int64(10), callbackFunc)
	lrucache.Add("key1", String("123456"))
	lrucache.Add("k2", String("v2"))
	lrucache.Add("k3", String("v3"))
	lrucache.Add("k4", String("v4"))

	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("OnEvicted invocation failed, expect keys equal to %s, actual is %s", expect, keys)
	}
}
