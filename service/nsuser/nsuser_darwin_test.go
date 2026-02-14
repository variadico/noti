package nsuser

import "testing"

func TestEscapeAS(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"plain", "hello world", "hello world"},
		{"double quotes", `say "hello"`, `say \"hello\"`},
		{"backslash", `path\to\file`, `path\\to\\file`},
		{"both", `"back\slash"`, `\"back\\slash\"`},
		{"tab", "tabs\there", "tabs here"},
		{"newline", "line1\nline2", "line1 line2"},
		{"carriage return", "line1\rline2", "line1 line2"},
		{"mixed whitespace", "a\tb\nc\rd", "a b c d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeAS(tt.in)
			if got != tt.want {
				t.Errorf("escapeAS(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBuildScript(t *testing.T) {
	tests := []struct {
		name string
		n    Notification
		want string
	}{
		{
			name: "message only",
			n:    Notification{InformativeText: "hello"},
			want: `display notification "hello"`,
		},
		{
			name: "all fields",
			n: Notification{
				Title:           "Title",
				Subtitle:        "Sub",
				InformativeText: "Body",
				SoundName:       "Ping",
			},
			want: `display notification "Body" with title "Title" subtitle "Sub" sound name "Ping"`,
		},
		{
			name: "title and message",
			n: Notification{
				Title:           "T",
				InformativeText: "M",
			},
			want: `display notification "M" with title "T"`,
		},
		{
			name: "escaping in fields",
			n: Notification{
				Title:           `He said "hi"`,
				InformativeText: `path\to\file`,
			},
			want: `display notification "path\\to\\file" with title "He said \"hi\""`,
		},
		{
			name: "content image ignored",
			n: Notification{
				InformativeText: "msg",
				ContentImage:    "/some/image.png",
			},
			want: `display notification "msg"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildScript(&tt.n)
			if got != tt.want {
				t.Errorf("buildScript() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
