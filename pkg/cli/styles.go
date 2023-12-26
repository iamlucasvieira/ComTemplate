package cli

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	List = lipgloss.NewStyle().
		BorderForeground(subtle).
		MarginTop(2).
		MarginBottom(1)

	ListHeader = lipgloss.NewStyle().
			Foreground(special).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			Render

	ListItem = lipgloss.NewStyle().PaddingLeft(2).Render

	tick = lipgloss.NewStyle().SetString("-").
		Foreground(highlight).
		PaddingRight(1).
		String()

	ListItemTick = func(s string) string {
		return tick + lipgloss.NewStyle().
			Render(s)
	}
)

func RenderList(title string, items []string) string {
	// Transform []string into []ListItemTick
	parts := make([]string, 0, len(items)+1)
	parts = append(parts, ListHeader(title))

	// Add each item, transformed by ListItemTick, to the slice
	for _, item := range items {
		parts = append(parts, ListItemTick(item))
	}

	l := lipgloss.JoinVertical(
		lipgloss.Top,
		parts...,
	)

	return List.Render(l)
}
