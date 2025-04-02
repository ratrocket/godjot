package testx

import (
	"errors"
	"reflect"
	"testing"
)

func AssertPanic(t *testing.T, label string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("no panic: %s", label)
		}
	}()
	f()
}

// AssertEmpty asserts that the specified object is empty. I.e. nil, "",
// false, 0 or either a slice or a channel with len == 0.
func AssertEmpty(t *testing.T, msg string, object any) bool {
	t.Helper()

	pass := isEmpty(object)
	if !pass {
		if msg != "" {
			t.Errorf("%s: should be empty, but was %v", msg, object)
		} else {
			t.Errorf("should be empty, but was %v", object)
		}
	}
	return pass
}

// AssertEqual ensures that want "equals" actual, in whatever sense
// reflect.DeepEqual holds.
func AssertEqual(t *testing.T, msg string, want, actual any) bool {
	t.Helper()

	if !reflect.DeepEqual(want, actual) {
		if msg != "" {
			t.Errorf("%s: want %v, got %v", msg, want, actual)
		} else {
			t.Errorf("want %v, got %v", want, actual)
		}
		return false
	}
	return true
}

func AssertNotEqual(t *testing.T, msg string, want, actual any) bool {
	t.Helper()

	if reflect.DeepEqual(want, actual) {
		if msg != "" {
			t.Errorf("%s: want %v, got %v", msg, want, actual)
		} else {
			t.Errorf("want %v, got %v", want, actual)
		}
	}
	return true
}

// AssertFalse ensures value is false.
func AssertFalse(t *testing.T, msg string, value bool) bool {
	t.Helper()
	if value {
		if msg != "" {
			t.Errorf("%s: want false, got true", msg)
		} else {
			t.Errorf("want false, got true")
		}
		return false
	}
	return true
}

// AssertIsType checks that the types of expectedType and object are the
// same.
func AssertIsType(t *testing.T, msg string, expectedType any, object any) bool {
	t.Helper()

	if !objectsAreEqual(reflect.TypeOf(object), reflect.TypeOf(expectedType)) {
		if msg != "" {
			t.Errorf("%s: object expected to be of type %v, but was %v",
				msg, reflect.TypeOf(expectedType), reflect.TypeOf(object))
		} else {
			t.Errorf("object expected to be of type %v, but was %v",
				reflect.TypeOf(expectedType), reflect.TypeOf(object))
		}
		return false
	}
	return true
}

// AssertNil checks that object is nil.
func AssertNil(t *testing.T, msg string, object any) bool {
	t.Helper()
	if object != nil {
		if msg != "" {
			t.Errorf("%s: want nil, got not nil", msg)
		} else {
			t.Errorf("want nil, got not nil")
		}
		return false
	}
	return true
}

// AssertNotNil checks that object is not nil.
func AssertNotNil(t *testing.T, msg string, object any) bool {
	t.Helper()
	if object == nil {
		if msg != "" {
			t.Errorf("%s: want not nil, got nil", msg)
		} else {
			t.Errorf("want not nil, got nil")
		}
		return false
	}
	return true
}

// AssertTrue checks that value is true.
func AssertTrue(t *testing.T, msg string, value bool) bool {
	t.Helper()
	if !value {
		if msg != "" {
			t.Errorf("%s: want true, got false", msg)
		} else {
			t.Errorf("want true, got false")
		}
		return false
	}
	return true
}

// AssertNilError checks that err is nil.
func AssertNilError(t *testing.T, msg string, err error) bool {
	t.Helper()
	if err != nil {
		if msg != "" {
			t.Errorf("%s: want nil error, got %v", msg, err)
		} else {
			t.Errorf("want nil error, got %v", err)
		}
		return false
	}
	return true
}

// AssertErrorIs checks that err is wantErr.
//
// You can use this instead of AssertNilError, just set wantErr==nil.
func AssertErrorIs(t *testing.T, msg string, wantErr, err error) bool {
	t.Helper()
	if !errors.Is(err, wantErr) {
		if msg != "" {
			t.Errorf("%s: want error %v, got %v", msg, wantErr, err)
		} else {
			t.Errorf("want error %v, got %v", wantErr, err)
		}
		return false
	}
	return true
}
