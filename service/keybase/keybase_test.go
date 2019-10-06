package keybase

import (
	"reflect"
	"testing"
	"time"
)

func TestKeybase_prepareArgs(t *testing.T) {
	// Initial args to the keybase cli
	base := []string{"chat", "send"}

	cases := []struct {
		name string
		// input
		n *Notification
		// expected outputs
		args []string
		err  error
	}{
		{
			name: "Minimal keybase notification",
			n:    &Notification{Conversation: "capn_picard", Message: "salutations"},
			args: append(base, "capn_picard", "salutations"),
			err:  nil,
		},
		{
			name: "Channel specified notification",
			n: &Notification{
				Conversation:      "golang-dev",
				ChannelName:       "off-topic",
				Message:           "salutations",
				ExplodingLifetime: 30 * time.Minute,
			},
			args: append(base, "--channel", "off-topic", "--exploding-lifetime", "30m0s", "golang-dev", "salutations"),
			err:  nil,
		},
		{
			name: "Missing required `Conversation`",
			n:    &Notification{Conversation: "", Message: "salutations"},
			args: nil,
			err:  ErrorMissingConversation,
		},
		{
			name: "Missing required `Message`",
			n:    &Notification{Conversation: "homer", Message: ""},
			args: nil,
			err:  ErrorMissingMessage,
		},
		{
			name: "Invalid exploding time",
			n:    &Notification{Conversation: "biggie", Message: "salutations", ExplodingLifetime: -1},
			args: nil,
			err:  ErrorBadExplodingTime,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			args, err := prepareArgs(tt.n)
			if !reflect.DeepEqual(args, tt.args) {
				t.Errorf("got wrong args: have=%v; want=%v;", args, tt.args)
			}
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("got wrong error: have=%v; want=%v;", err, tt.err)
			}
		})
	}

}
