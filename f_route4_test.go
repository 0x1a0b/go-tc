package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoute4(t *testing.T) {
	tests := map[string]struct {
		val  Route4
		err1 error
		err2 error
	}{
		"empty":  {err1: fmt.Errorf("Route4 options are missing")},
		"simple": {val: Route4{ClassID: 0xFFFF, To: 2, From: 3, IIf: 4}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalRoute4(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Route4{}
			err2 := unmarshalRoute4(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Route4 missmatch (want +got):\n%s", diff)
			}
		})
	}
}