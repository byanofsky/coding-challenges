package main

import (
	"testing"
	"time"
)

func TestDictionaryGet(t *testing.T) {
	d := NewDictionary()
	_, ok := d.Get("key")
	if ok != false {
		t.Fatalf("Get before Set. ok=%t. want=%t", ok, false)
	}

	d.Set("key", "value")

	actual, ok := d.Get("key")
	if ok != true {
		t.Fatalf("Get after Set. ok=%t. want=%t", ok, true)
	}
	if actual != "value" {
		t.Fatalf("Get after Set. result=%s. want=%s", actual, "value")
	}
}

func TestDictionaryGetExpire(t *testing.T) {
	d := NewDictionary()

	d.SetWithExpire("key", "value", 1)

	actual, ok := d.Get("key")
	if ok != true {
		t.Fatalf("Get after Set. ok=%t. want=%t", ok, true)
	}
	if actual != "value" {
		t.Fatalf("Get after Set. result=%s. want=%s", actual, "value")
	}

	// TODO: Mock time instead of sleep
	time.Sleep(2 * time.Second)

	_, ok = d.Get("key")
	if ok != false {
		t.Fatalf("Get after expire. ok=%t. want=%t", ok, false)
	}
}
