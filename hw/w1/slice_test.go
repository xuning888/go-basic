package w1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteAtAll(t *testing.T) {
	type testCase struct {
		name      string
		slice     []int
		idx       int
		wantSlice []int
		wantErr   error
	}

	testCases := []testCase{
		{
			name:      "test index out of range, idx = -1",
			slice:     []int{1, 2, 3, 4},
			idx:       -1,
			wantSlice: nil,
			wantErr:   ErrIndexOutOfRange,
		},
		{
			name:      "test index out of range, idx = 4",
			slice:     []int{1, 2, 3, 4},
			idx:       -1,
			wantSlice: nil,
			wantErr:   ErrIndexOutOfRange,
		},
		{
			name:      "success delete idx = 0",
			slice:     []int{1, 2, 3, 4},
			idx:       0,
			wantSlice: []int{2, 3, 4},
			wantErr:   nil,
		},
		{
			name:      "success delete idx = len(s) - 1",
			slice:     []int{1, 2, 3, 4},
			idx:       3,
			wantSlice: []int{1, 2, 3},
			wantErr:   nil,
		},
		{
			name:      "success delete idx = 2",
			slice:     []int{1, 2, 3, 4},
			idx:       2,
			wantSlice: []int{1, 2, 4},
			wantErr:   nil,
		},
		{
			name:      "success delete idx = 0",
			slice:     []int{},
			idx:       0,
			wantSlice: []int{},
			wantErr:   nil,
		},
	}

	valid := func(t *testing.T, tc testCase, actual []int, err error) {
		if err != nil {
			assert.Error(t, err, err.Error())
			assert.Nil(t, actual)
			return
		}
		assert.Equal(t, tc.wantSlice, actual)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := DeleteAt(tc.slice, tc.idx)
			valid(t, tc, actual, err)
			actual, err = DeleteAt1(tc.slice, tc.idx)
			valid(t, tc, actual, err)
			actual, err = DeleteAt3(tc.slice, tc.idx)
			valid(t, tc, actual, err)
			actual, err = DeleteAt2(tc.slice, tc.idx)
			valid(t, tc, actual, err)
		})
	}
}

func TestDeleteAt3(t *testing.T) {
	sli := []int{1, 2, 3, 4}
	actual, err := DeleteAt3(sli, 0)
	if err != nil {
		assert.Error(t, err, err.Error())
	}
	assert.Equal(t, []int{2, 3, 4}, actual)

	sli1 := []float32{1.0, 2.0, 3.0, 4.0}
	actual1, err := DeleteAt3(sli1, 1)
	if err != nil {
		assert.Error(t, err, err.Error())
	}
	assert.Equal(t, []float32{1.0, 3.0, 4.0}, actual1)
}
