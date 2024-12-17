package tui

import (
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"regexp"
)

// Style definitions
var (
	aliasStyle = lipgloss.NewStyle().Align(lipgloss.Center).PaddingTop(2)

	inputContainerStyle = lipgloss.NewStyle().Width(30).Align(lipgloss.Left).BorderForeground(lipgloss.Color("62")).Padding(0, 1)

	titleStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#FFFFFF"}).MarginBottom(1).Padding(4, 0, 2, 0)
)

// Register represents the registration terminal UI
type Register struct {
	width        int
	height       int
	txt          textinput.Model
	placeholder  string
	bindURL      string
	err          error
	registerFunc doReg
}

type doReg func(string) (string, error)

// newRegister creates a new registration terminal instance
func newRegister(reg doReg) *Register {
	ti := textinput.New()
	ti.Placeholder = "typing your connect alias"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 28 // Slightly less than container width to account for padding
	ti.TextStyle = lipgloss.NewStyle().Align(lipgloss.Center)

	return &Register{
		txt:          ti,
		registerFunc: reg,
	}
}

// Init implements tea.Model
func (r *Register) Init() tea.Cmd {
	return textinput.Blink
}

var isAlphabetic = regexp.MustCompile(`^[a-z0-9]{6,20}$`).MatchString

// Update implements tea.Model
func (r *Register) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return r, tea.Quit
		case tea.KeyEnter:
			value := r.txt.Value()
			var err error
			if "" != value && isAlphabetic(value) {
				r.bindURL, err = r.registerFunc(value)
				if nil != err {
					r.err = err
				}
			} else {
				r.err = errors.New("!! Alias validation failed ")
			}
			return r, nil
		}
		r.err = nil
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		return r, nil
	}
	var cmd tea.Cmd
	r.txt, cmd = r.txt.Update(msg)
	return r, cmd
}

// View implements tea.Model
func (r *Register) View() string {

	ok := isAlphabetic(r.txt.Value())

	if "" != r.bindURL && ok {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Width(r.width).Render("Registration Information"),
			lipgloss.NewStyle().Align(lipgloss.Center).
				PaddingLeft(4).
				Render("Please open the following URL through your browser for authentication"),
			aliasStyle.Render(r.bindURL),
		)
	}

	tips := lipgloss.NewStyle().Align(lipgloss.Center)

	if !ok {
		tips.Foreground(lipgloss.Color("31"))
	}

	if nil != r.err {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Width(r.width).Render("Registration Information"),
			tips.Render("Enter your alias, length range 6-20, can only contain a-z,0-9"),
			lipgloss.NewStyle().PaddingTop(1).Render(r.err.Error()),
			aliasStyle.Render(
				inputContainerStyle.Render(r.txt.View()),
			),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Width(r.width).Render("Registration Information"),
		tips.Render("Enter your alias, length range 6-20, can only contain a-z,0-9"),
		aliasStyle.Render(
			inputContainerStyle.Render(r.txt.View()),
		),
	)
}
