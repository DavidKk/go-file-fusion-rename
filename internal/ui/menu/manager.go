package menu

import (
	"file-fusion-rename/internal/app"
	"file-fusion-rename/internal/ui/dialog"
	"file-fusion-rename/internal/ui/tray"

	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
)

// 菜单管理器
type Manager struct {
	state   *app.State
	dialog  *dialog.Manager
	tray    *tray.Manager
	renamer *Renamer

	// 菜单项
	s2tItem    *systray.MenuItem
	t2sItem    *systray.MenuItem
	selectPath *systray.MenuItem
	quitItem   *systray.MenuItem
}

// 创建菜单管理器
func NewManager(state *app.State, dialogMgr *dialog.Manager, trayMgr *tray.Manager) *Manager {
	return &Manager{
		state:   state,
		dialog:  dialogMgr,
		tray:    trayMgr,
		renamer: NewRenamer(state, dialogMgr),
	}
}

// 创建托盘菜单
func (m *Manager) CreateMenus() {
	m.t2sItem = systray.AddMenuItem("繁转简", "繁体中文转换为简体中文")
	m.s2tItem = systray.AddMenuItem("简转繁", "简体中文转换为繁体中文")

	systray.AddSeparator()
	m.selectPath = systray.AddMenuItem("选择文件", "选择文件或文件夹进行转换（文件名+内容）")

	systray.AddSeparator()
	m.quitItem = systray.AddMenuItem("退出", "退出应用程序")

	m.updateDirectionMenus()

	go m.handleEvents()
}

// 更新转换方向菜单的选中状态和托盘标题
func (m *Manager) updateDirectionMenus() {
	m.s2tItem.Uncheck()
	m.t2sItem.Uncheck()

	direction := m.state.GetConvertDirection()

	switch direction {
	case "s2t":
		m.s2tItem.Check()
	case "t2s":
		m.t2sItem.Check()
	}

	m.tray.UpdateTitle(direction)
}

// 处理菜单事件
func (m *Manager) handleEvents() {
	for {
		select {
		case <-m.s2tItem.ClickedCh:
			m.state.SetConvertDirection("s2t")
			m.updateDirectionMenus()

		case <-m.t2sItem.ClickedCh:
			m.state.SetConvertDirection("t2s")
			m.updateDirectionMenus()

		case <-m.selectPath.ClickedCh:
			m.handlePathSelection()

		case <-m.quitItem.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// 处理文件或文件夹选择，选择后直接显示确认对话框
func (m *Manager) handlePathSelection() {
	path, err := m.dialog.SelectFileOrFolder()
	if err != nil {
		if err == zenity.ErrCanceled {
			return
		}

		m.dialog.ShowError("选择失败: " + err.Error())
		return
	}

	err = m.renamer.ShowConfirmationForPath(path)
	if err != nil {
		if err == zenity.ErrCanceled {
			return
		}

		m.dialog.ShowError("确认对话框错误: " + err.Error())
		return
	}

	go m.renamer.PerformRenameForPath(path)
}
