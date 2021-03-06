package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPolice(t *testing.T) {
	tests := map[string]struct {
		val  Police
		err1 error
		err2 error
	}{
		"simple":          {val: Police{AvRate: 1337, Result: 42}},
		"invalidArgument": {val: Police{AvRate: 1337, Result: 42, Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"tbfOnly": {val: Police{Tbf: &Policy{
			Index: 0x0, Action: 0x2, Limit: 0x0, Burst: 0x4c4b40, Mtu: 0x2400,
			Rate:     RateSpec{CellLog: 0x6, Linklayer: 0x1, Overhead: 1, CellAlign: 0xffff, Mpu: 1, Rate: 0x7d},
			PeakRate: RateSpec{CellLog: 1, Linklayer: 1, Overhead: 1, CellAlign: 1, Mpu: 1, Rate: 1},
		}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalPolice(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Police{}
			err2 := unmarshalPolice(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Police missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalPolice(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
