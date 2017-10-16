package command

import "testing"

func TestCommandName(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{args: []string{},
			want: "noti"},
		{args: []string{"git"},
			want: "git"},
		{args: []string{"git", "commit"},
			want: "git commit"},
		{args: []string{"ls", "-l"},
			want: "ls"},
		{args: []string{"foo", "bar", "fizz", "buzz"},
			want: "foo bar"},
	}

	for i, c := range cases {
		have := commandName(c.args)
		if have != c.want {
			t.Error("Unexpected command name")
			t.Errorf("%d - have=%q; want=%q", i, have, c.want)
		}
	}
}
