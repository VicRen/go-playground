package mock

import (
	"errors"
	"github.com/VicRen/go-play-ground/mock/mock_gen"
	"testing"

	"github.com/golang/mock/gomock"
)

var dib *DiB

func setupTestDib(t *testing.T) func(t *testing.T) {
	dib = NewSomeStructDiB(NewDiA())
	return func(t *testing.T) {
		dib = nil
	}
}

func TestDiB_SomeMethod(t *testing.T) {
	tearDownTestDib := setupTestDib(t)
	defer tearDownTestDib(t)

	ret := dib.SomeMethod()
	if ret != "7" {
		t.Fatalf("expected result %s; got %s", "7", ret)
	}
}

var (
	dib2   *DiB
	mocDia *mock_mock.MockDiA
)

func setupTestDib2(t *testing.T) func(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mocDia = mock_mock.NewMockDiA(mockCtrl)
	dib2 = NewSomeStructDiB(mocDia)
	return func(t *testing.T) {
		dib2 = nil
		mocDia = nil
	}
}

func TestDiB_SomeMethod2(t *testing.T) {
	tearDownTestDib := setupTestDib2(t)
	defer tearDownTestDib(t)

	mocDia.EXPECT().SomeMethodDiA().Return(8, nil)
	ret := dib2.SomeMethod()
	if ret != "8" {
		t.Fatalf("expected result %s; got %s", "8", ret)
	}
}

func TestDiB_SomeMethod3(t *testing.T) {
	tt := []struct {
		name        string
		outOfA int
		expectedRet string
		expectedErr error
	}{
		{"No Error", 8, "8", nil},
		{"Error", 8, "", errors.New("error occur")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tearDownTestDib := setupTestDib2(t)
			defer tearDownTestDib(t)

			mocDia.EXPECT().SomeMethodDiA().Return(tc.outOfA, tc.expectedErr)
			ret := dib2.SomeMethod()
			if ret != tc.expectedRet {
				t.Fatalf("expected result %s; got %s", tc.expectedRet, ret)
			}
		})
	}
}
