package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/sorokin-vladimir/tele/internal/store"
	"github.com/sorokin-vladimir/tele/internal/ui/keys"
	"github.com/sorokin-vladimir/tele/internal/ui/layout"
)

// FolderSelectedMsg is emitted when the user confirms a folder selection.
// Filter is nil when "All Chats" is selected.
type FolderSelectedMsg struct {
	Filter *store.FolderFilter
}

var (
	selectedFolderStyle = lipgloss.NewStyle().Background(lipgloss.Color("63")).Foreground(lipgloss.Color("0"))
	normalFolderStyle   = lipgloss.NewStyle()
	activeFolderStyle   = lipgloss.NewStyle().Bold(true)
)

var allChatsFilter = store.FolderFilter{ID: 0, Title: "All Chats"}

type FoldersModel struct {
	folders      []store.FolderFilter // index 0 is always allChatsFilter
	cursor       int
	activeIdx    int
	width        int
	height       int
	focused      bool
	unreadCounts map[int]int
}

func NewFoldersModel() *FoldersModel {
	return &FoldersModel{
		folders:      []store.FolderFilter{allChatsFilter},
		unreadCounts: make(map[int]int),
	}
}

func (m *FoldersModel) SetFolders(folders []store.FolderFilter) {
	m.folders = make([]store.FolderFilter, 0, len(folders)+1)
	m.folders = append(m.folders, allChatsFilter)
	m.folders = append(m.folders, folders...)
	if m.cursor >= len(m.folders) {
		m.cursor = 0
	}
	if m.activeIdx >= len(m.folders) {
		m.activeIdx = 0
	}
}

func (m *FoldersModel) SetFocused(focused bool)            { m.focused = focused }
func (m *FoldersModel) Focused() bool                      { return m.focused }
func (m *FoldersModel) SetSize(width, height int)          { m.width = width; m.height = height }
func (m *FoldersModel) SetUnreadCounts(counts map[int]int) { m.unreadCounts = counts }
func (m *FoldersModel) Cursor() int                        { return m.cursor }
func (m *FoldersModel) Folders() []store.FolderFilter      { return m.folders }
func (m *FoldersModel) Context() keys.Context              { return keys.ContextFolders }

// SelectedFilter returns the currently active filter. Nil means All Chats.
func (m *FoldersModel) SelectedFilter() *store.FolderFilter {
	if m.activeIdx == 0 {
		return nil
	}
	f := m.folders[m.activeIdx]
	return &f
}

func (m *FoldersModel) HasFolders() bool {
	return len(m.folders) > 1
}

func (m *FoldersModel) Init() tea.Cmd { return nil }

func (m FoldersModel) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	if !m.focused {
		return &m, nil
	}
	action, ok := msg.(keys.ActionMsg)
	if !ok {
		return &m, nil
	}
	switch action.Action {
	case keys.ActionDown:
		if m.cursor < len(m.folders)-1 {
			m.cursor++
		}
	case keys.ActionUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.ActionConfirm:
		m.activeIdx = m.cursor
		var f *store.FolderFilter
		if m.activeIdx > 0 {
			ff := m.folders[m.activeIdx]
			f = &ff
		}
		return &m, func() tea.Msg { return FolderSelectedMsg{Filter: f} }
	}
	return &m, nil
}

func (m FoldersModel) View() string {
	var lines []string
	for i, f := range m.folders {
		label := m.formatEntry(f, i == m.activeIdx)
		style := normalFolderStyle
		if i == m.cursor && m.focused {
			style = selectedFolderStyle
		} else if i == m.activeIdx {
			style = activeFolderStyle
		}
		lines = append(lines, style.Render(label))
	}
	return joinLines(lines)
}

func (m FoldersModel) formatEntry(f store.FolderFilter, active bool) string {
	badge := ""
	if f.ID != 0 {
		if n, ok := m.unreadCounts[f.ID]; ok && n > 0 {
			if n > 99 {
				badge = "[99+]"
			} else {
				badge = fmt.Sprintf("[%d]", n)
			}
		}
	}
	prefix := "  "
	if active {
		prefix = "▸ "
	}
	nameWidth := m.width - runewidth.StringWidth(prefix) - runewidth.StringWidth(badge)
	if nameWidth < 1 {
		nameWidth = 1
	}
	name := truncateTo(f.Title, nameWidth)
	line := prefix + padRight(name, nameWidth)
	if badge != "" {
		line += badge
	}
	return line
}

func truncateTo(s string, maxW int) string {
	if runewidth.StringWidth(s) <= maxW {
		return s
	}
	// Reserve 1 column for the "…" ellipsis.
	targetW := maxW - 1
	if targetW <= 0 {
		return "…"
	}
	w := 0
	for i, r := range s {
		rw := runewidth.RuneWidth(r)
		if w+rw > targetW {
			return s[:i] + "…"
		}
		w += rw
	}
	return s
}

func padRight(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	for i := w; i < width; i++ {
		s += " "
	}
	return s
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
