package scheduler

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestScheduler(t *testing.T) {
	j := Job{
		ID:    1,
		Name:  "Upload report",
		Every: 10,
		Unit:  Seconds,
		Task: func() {
			log.Println("Uploading report...")
		},
	}
	tests := []struct {
		name     string
		job      Job
		expected Job
	}{
		{
			name:     "test_add",
			job:      j,
			expected: j,
		},
	}

	sch := New()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			sch.AddJob(test.job)

			jobs := sch.GetJobs()
			j := sch.GetJob(jobs[0].ID)
			if !cmp.Equal(test.expected, j, cmpopts.IgnoreFields(Job{}, "Task", "RegisteredAt")) {
				t.Errorf("expected %+v got %+v", test.expected, j)
			}
		})
	}

}
