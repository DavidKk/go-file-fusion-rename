package tray

import (
	"github.com/getlantern/systray"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Initialize() {
	systray.SetTitle("简")
	systray.SetTooltip("中文简繁体转换工具")
}

// 根据转换方向更新托盘标题
func (m *Manager) UpdateTitle(direction string) {
	switch direction {
	case "s2t":
		systray.SetTitle("繁")
	case "t2s":
		systray.SetTitle("简")
	default:
		systray.SetTitle("简")
	}
}

func (m *Manager) Quit() {
	systray.Quit()
}
