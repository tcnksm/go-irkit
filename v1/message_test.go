package irkit

import "testing"

func TestMessage_Validate(t *testing.T) {

	cases := []struct {
		msg     *Message
		success bool
	}{

		{
			&Message{
				Format: "raw",
				Freq:   38,
				Data:   []int{1, 2, 3, 4},
			},
			true,
		},

		{
			&Message{},
			false,
		},

		{
			&Message{
				Format: "raw",
				Freq:   38,
				Data:   []int{},
			},
			false,
		},

		{
			&Message{
				Format: "raw",
				Freq:   3,
				Data:   []int{},
			},
			false,
		},
	}

	for i, tc := range cases {
		actual := tc.msg.validate()
		if (actual != nil) == tc.success {
			t.Fatalf("#%d error: %v", i, actual)
		}
	}
}
