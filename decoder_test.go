package scraper

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

type a struct {
	Name string `html:".name"`
}

type b struct {
	Names []a `html:".value"`
}

type bb struct {
	Names []*a `html:".value"`
}

type c struct {
	B *b `html:".set"`
}

type d int

func TestDecoder(t *testing.T) {
	errFoo := errors.New("Foo")

	isErr := func(want, got error) bool {
		if want == got {
			return true
		}

		if want == nil || got == nil {
			return false
		}

		if w, ok1 := want.(*InvalidUnmarshalError); ok1 {
			if g, ok2 := got.(*InvalidUnmarshalError); ok2 {
				return *w == *g
			}
		}
		return false
	}

	tests := []struct {
		name      string
		inputFile string
		options   []Option
		want      interface{}
		got       interface{}
		wantErr   error
	}{
		{"A", "testdata/a.html", nil, &a{"  Hello World  "}, &a{}, nil},
		{"A (trim)", "testdata/a.html", []Option{TrimSpace()}, &a{"Hello World"}, &a{}, nil},
		{"B", "testdata/b.html", nil, &b{[]a{{"one"}, {"two"}, {"three"}}}, &b{}, nil},
		{"BB", "testdata/b.html", nil, &bb{[]*a{{"one"}, {"two"}, {"three"}}}, &bb{}, nil},
		{"C", "testdata/c.html", nil, &c{&b{[]a{{"one"}, {"two"}, {"three"}}}}, &c{}, nil},
		{"quick fail", "testdata/a.html", []Option{func(*Unmarshaler) error { return errFoo }}, &a{}, &a{}, errFoo},
		{"InvalidUnmarshalError 1", "testdata/a.html", nil, &a{}, nil, &InvalidUnmarshalError{reflect.TypeOf(nil), reflect.Ptr}},
		{"InvalidUnmarshalError 2", "testdata/a.html", nil, new(d), new(d), &InvalidUnmarshalError{reflect.TypeOf(d(0)), reflect.Struct}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.Open(test.inputFile)
			if err != nil {
				t.Fatalf("Failed to open %s: %v", test.inputFile, err)
			}

			dec := NewDecoder(file, test.options...)
			gotErr := dec.Decode(test.got)
			if isErr(test.wantErr, gotErr) {
				if gotErr == nil {
					if diff := deep.Equal(test.want, test.got); diff != nil {
						t.Error(diff)
					}
				}
			} else {
				t.Errorf("Wanted error %v got %v", test.wantErr, gotErr)
			}
		})
	}
}
