package media

import "strings"

// Mode is the selected image rendering strategy.
type Mode int

const (
	// ModeBlocks renders photos as ANSI half-block art (universal fallback).
	ModeBlocks Mode = iota
	// ModeKitty renders photos via the Kitty graphics protocol.
	ModeKitty
)

// DetectMode picks a rendering mode. override is the config value
// (photos.mode): "kitty" or "blocks" force a mode, "auto" (or empty) runs
// heuristics. getenv is usually os.Getenv; it is injected for testing.
//
// Kitty is enabled only for terminals with confirmed Unicode-placeholder
// support (kitty, Ghostty) and never inside a tmux/screen session (no
// passthrough this pass).
func DetectMode(override string, getenv func(string) string) Mode {
	switch strings.ToLower(override) {
	case "kitty":
		return ModeKitty
	case "blocks":
		return ModeBlocks
	}
	if getenv("TMUX") != "" || getenv("STY") != "" {
		return ModeBlocks
	}
	term := getenv("TERM")
	if strings.Contains(term, "kitty") || strings.Contains(term, "ghostty") {
		return ModeKitty
	}
	if getenv("KITTY_WINDOW_ID") != "" {
		return ModeKitty
	}
	if strings.EqualFold(getenv("TERM_PROGRAM"), "ghostty") {
		return ModeKitty
	}
	return ModeBlocks
}
