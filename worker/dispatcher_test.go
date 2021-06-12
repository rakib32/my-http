package worker

import (
	"testing"
)

func TestNewDispatcher(t *testing.T) {
	testCases := []struct {
		input              int
		expectedMaxWorkers int
		got                *Dispatcher
	}{
		{
			input:              10,
			expectedMaxWorkers: 10,
		},
	}

	for _, testCase := range testCases {
		testCase.got = NewDispatcher(testCase.input, make(chan Job))

		if testCase.got.maxWorkers != testCase.expectedMaxWorkers {
			t.Logf("input: %v, expected: %v, got: %v", testCase.input, testCase.expectedMaxWorkers, testCase.got.maxWorkers)
			t.Fail()
		}
	}
}

func TestDispatch(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
		got      int
	}{
		{
			input:    1,
			expected: 1,
		},
	}

	for _, testCase := range testCases {
		q := make(chan Job, 2)
		dispatcher := NewDispatcher(1, q)

		dispatcher.Run()
		dispatcher.AddJob(Job{
			Payload: 1,
			Handler: func(i interface{}) (interface{}, error) {
				testCase.got, _ = i.(int)

				if testCase.got != testCase.expected {
					t.Logf("input: %v, expected: %v, got: %v", testCase.input, testCase.expected, testCase.got)
					t.Fail()
				}
				return i, nil
			},
		})
		dispatcher.Stop()
	}
}
