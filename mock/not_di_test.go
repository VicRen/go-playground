package mock

import "testing"

var a *A

func setupTestA(t *testing.T) func(t *testing.T) {
	a = NewA()
	return func(t *testing.T) {
		a = nil
	}
}

func TestA_add(t *testing.T) {
	tearDownTestA := setupTestA(t)
	defer  tearDownTestA(t)

	ret := a.add(4, 4)
	if ret != 8 {
		t.Fatalf("expected result %d; got %d", 7, ret)
	}
}

var b *B

func setupTestB(t *testing.T) func(t *testing.T) {
	b = NewSomeStruct(*NewA())
	return func(t *testing.T) {
		b = nil
	}
}

func TestB_SomeMethod(t *testing.T) {
	tearDownTestB := setupTestB(t)
	defer tearDownTestB(t)

	ret := b.SomeMethod()
	if ret != "7" {
		t.Fatalf("expected result %s; got %s", "7", ret)
	}
}