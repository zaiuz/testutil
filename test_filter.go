package testutil

import "time"
import z "github.com/zaiuz/zaiuz"

type TestFilter struct {
	z.Filter
	Called     bool
	CallTime   time.Time
	Finished   bool
	FinishTime time.Time
}

func NewTestFilter() *TestFilter {
	zero := time.Time{}

	testFilter := &TestFilter{nil, false, zero, false, zero}
	testFilter.Filter = func(action z.Action) z.Action {
		return func(context *z.Context) z.Result {
			testFilter.Called = true
			testFilter.CallTime = time.Now()
			result := action(context)

			time.Sleep(time.Millisecond) // ensure finish time > call time
			testFilter.Finished = true
			testFilter.FinishTime = time.Now()
			return result
		}
	}

	return testFilter
}

func (tf *TestFilter) Reset() {
	tf.Called = false
	tf.CallTime = time.Time{}
	tf.Finished = false
	tf.FinishTime = time.Time{}
}

