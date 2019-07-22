package gocontainer

import (
	"testing"
)

func TestRegister(t *testing.T) {
	Register("test", 1)

	if MustGet("test") != 1 {
		t.Fail()
	}
}

func TestDeregister(t *testing.T) {
	Register("test", 1)
	Deregister("test")

	_, ok := Get("test")
	if ok {
		t.Fail()
	}
}

func TestHas(t *testing.T) {
	Register("test", 1)

	if !Has("test") {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	Register("test", 1)
	o, ok := Get("test")

	if !ok {
		t.Fail()
	}

	if o != 1 {
		t.Fail()
	}
}

func TestMustGet(t *testing.T) {
	Register("test", 1)
	o := MustGet("test")

	if o != 1 {
		t.Fail()
	}
}

func TestInvoke(t *testing.T) {
	Register("test", 1)
	Invoke("test", func(i int, ok bool) {
		if !ok {
			t.Fail()
		}

		if i != 1 {
			t.Fail()
		}
	})
}

func TestMustInvoke(t *testing.T) {
	Register("test", 1)
	MustInvoke("test", func(i int) {
		if i != 1 {
			t.Fail()
		}
	})
}

func TestMustInvokeMany(t *testing.T) {
	Register("test1", 1)
	Register("test2", 2)
	MustInvokeMany("test1", "test2")(func(x int, y int) {
		if x != 1 {
			t.Fail()
		}
		if y != 2 {
			t.Fail()
		}
	})
}
