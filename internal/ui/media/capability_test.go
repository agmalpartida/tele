package media_test

import (
	"testing"

	"github.com/sorokin-vladimir/tele/internal/ui/media"
	"github.com/stretchr/testify/require"
)

func TestDetectMode(t *testing.T) {
	cases := []struct {
		name     string
		override string
		env      map[string]string
		want     media.Mode
	}{
		{"kitty term", "auto", map[string]string{"TERM": "xterm-kitty"}, media.ModeKitty},
		{"kitty window id", "auto", map[string]string{"KITTY_WINDOW_ID": "1"}, media.ModeKitty},
		{"ghostty term", "auto", map[string]string{"TERM": "xterm-ghostty"}, media.ModeKitty},
		{"ghostty program", "auto", map[string]string{"TERM_PROGRAM": "ghostty"}, media.ModeKitty},
		{"plain xterm", "auto", map[string]string{"TERM": "xterm-256color"}, media.ModeBlocks},
		{"kitty under tmux", "auto", map[string]string{"TERM": "xterm-kitty", "TMUX": "/tmp/x"}, media.ModeBlocks},
		{"kitty under screen", "auto", map[string]string{"TERM": "xterm-kitty", "STY": "1.x"}, media.ModeBlocks},
		{"empty override falls through", "", map[string]string{"TERM": "xterm-kitty"}, media.ModeKitty},
		{"override kitty on plain", "kitty", map[string]string{"TERM": "xterm-256color"}, media.ModeKitty},
		{"override blocks on kitty", "blocks", map[string]string{"TERM": "xterm-kitty"}, media.ModeBlocks},
		{"override auto falls through", "auto", map[string]string{"TERM": "xterm-kitty"}, media.ModeKitty},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := media.DetectMode(tc.override, func(k string) string { return tc.env[k] })
			require.Equal(t, tc.want, got)
		})
	}
}
