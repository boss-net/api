package plugin

import (
	"fmt"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/boss-net/api/boss-plugin/pkg/entities/manifest_entities"
)

type versionRequire struct {
	minimalBossVersion ti.Model

	warning string
}

func newVersionRequire() versionRequire {
	minimalBossVersion := ti.New()
	minimalBossVersion.Placeholder = "Minimal Boss version"
	minimalBossVersion.CharLimit = 128
	minimalBossVersion.Prompt = "Minimal Boss version (press Enter to next step): "
	minimalBossVersion.Focus()

	return versionRequire{
		minimalBossVersion: minimalBossVersion,
	}
}

func (p versionRequire) MinimalBossVersion() string {
	return p.minimalBossVersion.Value()
}

func (p versionRequire) View() string {
	s := fmt.Sprintf("Edit minimal Boss version requirement, leave it blank by default\n%s\n", p.minimalBossVersion.View())
	if p.warning != "" {
		s += fmt.Sprintf("\033[31m%s\033[0m\n", p.warning)
	}
	return s
}

func (p *versionRequire) checkRule() bool {
	if p.minimalBossVersion.Value() == "" {
		p.warning = ""
		return true
	}

	_, err := manifest_entities.NewVersion(p.minimalBossVersion.Value())
	if err != nil {
		p.warning = "Invalid minimal Boss version"
		return false
	}

	p.warning = ""
	return true
}

func (p versionRequire) Update(msg tea.Msg) (subMenu, subMenuEvent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return p, SUB_MENU_EVENT_NONE, tea.Quit
		case "enter":
			// check if empty
			if !p.checkRule() {
				return p, SUB_MENU_EVENT_NONE, nil
			}
			return p, SUB_MENU_EVENT_NEXT, nil
		}
	}

	// update view
	var cmd tea.Cmd
	p.minimalBossVersion, cmd = p.minimalBossVersion.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return p, SUB_MENU_EVENT_NONE, tea.Batch(cmds...)
}

func (p versionRequire) Init() tea.Cmd {
	return nil
}

func (p *versionRequire) SetMinimalBossVersion(version string) {
	p.minimalBossVersion.SetValue(version)
}
