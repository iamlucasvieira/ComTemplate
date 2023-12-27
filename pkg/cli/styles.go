package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	marginLeft   = 2
	marginTop    = 1
	marginBottom = 1
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	Shell = lipgloss.NewStyle().
		MarginLeft(marginLeft).
		Render

	ShellMargin = lipgloss.NewStyle().
			MarginTop(marginTop).
			MarginBottom(marginBottom).
			MarginLeft(marginLeft).
			Render

	List = lipgloss.NewStyle().
		BorderForeground(subtle)

	Header = lipgloss.NewStyle().
		Foreground(special).
		Bold(true).
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

	TextHighlight = lipgloss.NewStyle().
			Foreground(highlight).
			Render
)

// RenderList renders a list of strings
func RenderList(title string, items []string) string {
	// Transform []string into []ListItemTick
	parts := make([]string, 0, len(items)+1)
	parts = append(parts, Header(title))

	// Add each item, transformed by ListItemTick, to the slice
	for _, item := range items {
		parts = append(parts, ListItemTick(item))
	}

	l := lipgloss.JoinVertical(
		lipgloss.Top,
		parts...,
	)

	return ShellMargin(List.Render(l))
}

// Write prints a list of strings to the terminal
func Write(items ...string) {
	vertical := lipgloss.JoinVertical(lipgloss.Top, items...)
	fmt.Println(ShellMargin(vertical))
}

// WriteNoMargin prints a list of strings to the terminal
func WriteNoMargin(items ...string) {
	vertical := lipgloss.JoinVertical(lipgloss.Top, items...)
	fmt.Println(Shell(vertical))
}
