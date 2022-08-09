package cola_test

import (
	"github.com/smartwalle/cola"
	"sync"
	"testing"
	"time"
)

func TestManager_Tick1(t *testing.T) {
	var m = cola.New()

	var victor string

	m.Add("k1", 1, func(key string) {
		victor = key
	}).Accept()

	m.Add("k2", 2, func(key string) {
		victor = key
	}).Accept()

	m.Add("k3", 3, func(key string) {
		victor = key
	}).Accept()

	var w = &sync.WaitGroup{}
	m.Tick(time.Second, cola.WithWaiter(w))
	w.Wait()

	if victor != "k3" {
		t.Fatal("期望值是: k3, 实际值是:", victor)
	}
}

func TestManager_Tick2(t *testing.T) {
	var m = cola.New()

	var victor string

	m.Add("k1", 1, func(key string) {
		victor = key
	}).Accept()

	m.Add("k2", 2, func(key string) {
		victor = key
	}).Accept()

	m.Add("k3", 3, func(key string) {
		victor = key
	}).Reject()

	var w = &sync.WaitGroup{}
	m.Tick(time.Second, cola.WithWaiter(w))
	w.Wait()

	if victor != "k2" {
		t.Fatal("期望值是: k2, 实际值是:", victor)
	}
}

func TestManager_Tick3(t *testing.T) {
	var m = cola.New()

	var victor string

	m.Add("k1", 1, func(key string) {
		victor = key
	}).Accept()

	m.Add("k2", 2, func(key string) {
		victor = key
	})

	m.Add("k3", 3, func(key string) {
		victor = key
	})

	var w = &sync.WaitGroup{}
	m.Tick(time.Second, cola.WithWaiter(w))
	w.Wait()

	if victor != "k1" {
		t.Fatal("期望值是: k1, 实际值是:", victor)
	}
}

func TestManager_Tick4(t *testing.T) {
	var m = cola.New()

	var victor string

	m.Add("k1", 1, func(key string) {
		victor = key
	})

	m.Add("k2", 2, func(key string) {
		victor = key
	})

	m.Add("k3", 3, func(key string) {
		victor = key
	}).Accept()

	var w = &sync.WaitGroup{}
	m.Tick(time.Second, cola.WithWaiter(w))
	w.Wait()

	if victor != "k3" {
		t.Fatal("期望值是: k3, 实际值是:", victor)
	}
}

func TestManager_Tick5(t *testing.T) {
	var m = cola.New()

	var victor string

	m.Add("k1", 1, func(key string) {
		victor = key
	})

	m.Add("k2", 2, func(key string) {
		victor = key
	})

	m.Add("k3", 3, func(key string) {
		victor = key
	})

	var w = &sync.WaitGroup{}
	m.Tick(time.Second, cola.WithWaiter(w))
	w.Wait()

	if victor != "" {
		t.Fatal("期望值是: , 实际值是:", victor)
	}
}

func TestManager_Tick6(t *testing.T) {
	var m = cola.New()

	var victor string

	m.Add("k1", 1, func(key string) {
		victor = key
	}).Reject()

	m.Add("k2", 2, func(key string) {
		victor = key
	}).Reject()

	m.Add("k3", 3, func(key string) {
		victor = key
	}).Reject()

	var w = &sync.WaitGroup{}
	m.Tick(time.Second, cola.WithWaiter(w))
	w.Wait()

	if victor != "" {
		t.Fatal("期望值是: , 实际值是:", victor)
	}
}

func BenchmarkManager_Tick(b *testing.B) {
	var m = cola.New()
	var w = &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		m.Add("a1", 1, func(key string) {
		}).Accept()

		m.Add("a2", 2, func(key string) {
		}).Accept()

		m.Tick(time.Nanosecond, cola.WithWaiter(w))
	}

	w.Wait()
}
