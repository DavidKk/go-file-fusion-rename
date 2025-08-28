package dialog

import (
	"github.com/ncruces/zenity"
)

// 对话框管理器
type Manager struct{}

// 创建对话框管理器
func NewManager() *Manager {
	return &Manager{}
}

// 显示文件或文件夹选择对话框
func (m *Manager) SelectFileOrFolder() (string, error) {
	path, err := zenity.SelectFile(
		zenity.Title("选择文件或文件夹进行转换（文件名+内容）"),
		zenity.Directory(),
	)
	if err == nil {
		return path, nil
	}

	if err == zenity.ErrCanceled {
		return zenity.SelectFile(
			zenity.Title("选择文件或文件夹进行转换（文件名+内容）"),
		)
	}

	return "", err
}

// 显示错误对话框
func (m *Manager) ShowError(message string) {
	go func() {
		zenity.Error(message, zenity.Title("错误"))
	}()
}
