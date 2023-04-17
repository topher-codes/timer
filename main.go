package main

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/timer"
)

const (
	padding  = 2
	maxWidth = 160
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var minutes int

func main() {
    // Ask the user to enter how many minutes the timer should run for
    fmt.Print("Enter the number of minutes: ")
    fmt.Scanln(&minutes)
    fmt.Println("Timer started for", minutes, "minutes")


	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
        timer: timer.NewWithInterval(time.Minute * time.Duration(minutes), time.Second),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress progress.Model
    timer timer.Model
}

func (m model) InitTimer() tea.Cmd {
    return m.timer.Init()
}

func (m model) Init() tea.Cmd {
    m.InitTimer() 
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case timer.TickMsg:
        var tcmd tea.Cmd
        m.timer, tcmd = m.timer.Update(msg)
        return m, tcmd

	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

        // update the timer
        m.timer, _ = m.timer.Update(msg)


		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
        // Here we take the minutes entered by the user and convert it to a percentage
		cmd := m.progress.IncrPercent(0.01 / float64(minutes) * 1.6)
		return m, tea.Batch(tickCmd(), cmd, )

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
        pad + m.timer.View() + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
