package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	columnStyle = lipgloss.NewStyle().
			Padding(padding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("243"))
	focusedStyle = lipgloss.NewStyle().
			Padding(padding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

const (
	list_fraction = 4
	padding       = 2
)

type Model struct {
	lists         []list.Model
	focused       Status
	ready         bool
	inputing      bool
	inputs        [2]textinput.Model
	focused_input int
	width         int
	height        int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func newModel() *Model {
	m := loadModel()
	return &m
}

func (m *Model) initInputs(width, height int) {

	title_input := textinput.New()
	title_input.Width = width
	title_input.CharLimit = 32
	title_input.Placeholder = "Title"
	title_input.Focus()

	description_input := textinput.New()
	description_input.Width = width
	description_input.CharLimit = 64
	description_input.Placeholder = "Description"

	m.inputs = [2]textinput.Model{
		title_input,
		description_input,
	}
}

func (m *Model) initLists(width, height int) {
	for i, l := range m.lists {
		newList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
		newList.SetItems(l.Items())
		newList.Title = l.Title
		newList.SetShowHelp(false)
		newList.SetShowStatusBar(false)
		newList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("207")).Background(lipgloss.Color("56"))
		// set the style so that the description of the task does not overflow
		//newList.Styles.Item = lipgloss.NewStyle().Width(width/len(m.lists) - 2).Padding(1)

		m.lists[i] = newList
	}
	m.lists[Todo].Title = "Todo"
	m.lists[Doing].Title = "Doing"
	m.lists[Done].Title = "Done"
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if !m.inputing {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			if !m.ready {
				m.initLists(msg.Width, msg.Height-(2+padding*2))
				m.initInputs(msg.Width/2, msg.Height)
				focusedStyle.Width(msg.Width/len(m.lists) - 2)
				columnStyle.Width(msg.Width/len(m.lists) - 2)
				m.ready = true
			}
			m.width = msg.Width
			m.height = msg.Height

			focusedStyle.Width(msg.Width/len(m.lists) - 2)
			columnStyle.Width(msg.Width/len(m.lists) - 2)
			focusedStyle.Height(msg.Height - 2)
			columnStyle.Height(msg.Height - 2)

			for _, list := range m.lists {
				list.SetHeight(msg.Height - 2)
			}

		case tea.KeyMsg:
			switch msg.String() {
			case "left", "h":
				m.focused = PrevStatus(m.focused)
			case "right", "l":
				m.focused = NextStatus(m.focused)

			case "n":
				m.lists[NextStatus(m.focused)].InsertItem(len(m.lists[NextStatus(m.focused)].Items()), m.lists[m.focused].SelectedItem())
				m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())

			case "p":
				m.lists[PrevStatus(m.focused)].InsertItem(len(m.lists[PrevStatus(m.focused)].Items()), m.lists[m.focused].SelectedItem())
				m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())

			case "x":
				m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())
			case "a":
				m.inputs[0].SetValue("")
				m.inputs[1].SetValue("")
				m.inputing = true

			case "ctrl+c", "q":
				saveModel(m)
				return m, tea.Quit
			}

			m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
		}
	} else {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				saveModel(m)
				return m, tea.Quit

			case "tab":
				if m.focused_input == 0 {
					m.focused_input = 1
				} else {
					m.focused_input = 0
				}
				m.inputs[m.focused_input].Focus()

			case "esc":
				m.inputing = false
				m.focused_input = 0
				m.inputs[0].Focus()
				m.inputs[0].SetValue("")
				m.inputs[1].SetValue("")

			case "enter":
				m.inputing = false
				title := m.inputs[0].Value()
				description := m.inputs[1].Value()
				m.lists[Todo].InsertItem(0, Task{title, description, Todo})
			}
		}
		m.inputs[m.focused_input], cmd = m.inputs[m.focused_input].Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.ready {
		if !m.inputing {
			var todo_view, doing_view, done_view string

			if m.focused == Todo {
				todo_view = focusedStyle.Render(m.lists[Todo].View())
				doing_view = columnStyle.Render(m.lists[Doing].View())
				done_view = columnStyle.Render(m.lists[Done].View())
			} else if m.focused == Doing {
				todo_view = columnStyle.Render(m.lists[Todo].View())
				doing_view = focusedStyle.Render(m.lists[Doing].View())
				done_view = columnStyle.Render(m.lists[Done].View())
			} else if m.focused == Done {
				todo_view = columnStyle.Render(m.lists[Todo].View())
				doing_view = columnStyle.Render(m.lists[Doing].View())
				done_view = focusedStyle.Render(m.lists[Done].View())
			}

			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				todo_view,
				doing_view,
				done_view,
			)
		} else {
			s := lipgloss.JoinVertical(lipgloss.Left,
				m.inputs[0].View(),
				m.inputs[1].View(),
			)
			padding_h := (m.width - lipgloss.Width(s)) / 2
			padding_v := (m.height - lipgloss.Height(s)) / 2

			inner_style := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Padding(2)

			style := lipgloss.NewStyle().
				Padding(padding_v, padding_h)

			return fmt.Sprintf("%d", padding_h) + style.Render(inner_style.Render(s))
		}
	} else {
		return "Loading..."
	}
}

func saveModel(m Model) {
	json, err := json.Marshal(ExportModel(m))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile("tasks.json", json, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func cleanModel() Model {
	return Model{
		lists: []list.Model{
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		},
		focused:  Todo,
		ready:    false,
		inputing: false,
	}
}

func loadModel() Model {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return cleanModel()
	}

	var em ExportableModel
	err = json.Unmarshal(data, &em)
	if err != nil {
		fmt.Println(err)
		return Model{}
	}

	m := Model{
		inputing:      false,
		focused_input: 0,
		ready:         false,
		focused:       Todo,
	}

	for i, l := range em.Lists {
		m.lists = append(m.lists, list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0))
		for _, item := range l.Items {
			m.lists[i].InsertItem(len(m.lists[i].Items()), Task{item.Title, item.Description, Todo})
		}
	}

	return m
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
