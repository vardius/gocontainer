package gocontainer

import (
	"testing"
)

func TestRegister(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)

	if MustGet("test") != 1 {
		t.Fail()
	}

	GlobalContainer = nil

	assertPanic(t, func() {
		Register("test", 1)
	})
}

func TestDeregister(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)
	Deregister("test")

	_, ok := Get("test")
	if ok {
		t.Fail()
	}

	GlobalContainer = nil

	assertPanic(t, func() {
		Deregister("test")
	})
}

func TestHas(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)

	if !Has("test") {
		t.Fail()
	}

	GlobalContainer = nil

	assertPanic(t, func() {
		Has("test")
	})
}

func TestGet(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)
	o, ok := Get("test")

	if !ok {
		t.Fail()
	}

	if o != 1 {
		t.Fail()
	}

	GlobalContainer = nil

	assertPanic(t, func() {
		Get("test")
	})
}

func TestMustGet(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)
	o := MustGet("test")

	if o != 1 {
		t.Fail()
	}

	GlobalContainer = nil

	assertPanic(t, func() {
		MustGet("test")
	})
}

func TestInvoke(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)

	Invoke("test", func(i int, ok bool) {
		if !ok {
			t.Fail()
		}

		if i != 1 {
			t.Fail()
		}
	})

	GlobalContainer = nil

	assertPanic(t, func() {
		Invoke("test", func(i int) {})
	})
}

func TestMustInvoke(t *testing.T) {
	GlobalContainer = New()

	Register("test", 1)

	MustInvoke("test", func(i int) {
		if i != 1 {
			t.Fail()
		}
	})

	GlobalContainer = nil

	assertPanic(t, func() {
		MustInvoke("test", func(i int) {})
	})
}

func TestMustInvokeMany(t *testing.T) {
	GlobalContainer = New()

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

	GlobalContainer = nil

	assertPanic(t, func() {
		MustInvokeMany("test1", "test2")
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	f()
}
