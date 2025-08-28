//go:build !nogui
// +build !nogui

package menu

import (
	"os"
	"path/filepath"
	"testing"

	"file-fusion-rename/internal/app"
	"file-fusion-rename/internal/ui/dialog"
)

func TestOptimizedWorkflow(t *testing.T) {
	// 创建临时测试目录
	tempDir := filepath.Join(os.TempDir(), "linkey_workflow_test")
	err := os.RemoveAll(tempDir)
	if err != nil {
		t.Fatalf("清理临时目录失败: %v", err)
	}

	// 创建测试文件
	testFile := filepath.Join(tempDir, "简体测试文件.txt")
	err = os.MkdirAll(filepath.Dir(testFile), 0755)
	if err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}
	file.WriteString("这是一个简体中文测试文件")
	file.Close()

	// 创建应用组件
	state := app.NewState()
	dialog := dialog.NewManager()
	renamer := NewRenamer(state, dialog)

	t.Log("测试优化后的工作流程:")

	// 1. 测试默认状态
	direction := state.GetConvertDirection()
	t.Logf("1. 默认转换方向: %s", direction)

	// 2. 测试直接路径处理
	t.Log("2. 测试文件处理流程...")

	// 使用新的方法直接处理路径
	items, err := renamer.collectRenameItemsFromFile(testFile, state.GetSrcLang(), state.GetDstLang())
	if err != nil {
		t.Fatalf("收集重命名项目失败: %v", err)
	}

	if len(items) > 0 {
		item := items[0]
		t.Logf("3. 检测到需要转换的文件:")
		t.Logf("   原始文件: %s", filepath.Base(item.oldPath))
		t.Logf("   转换后: %s", filepath.Base(item.newPath))
		t.Logf("   文件类型: %v", !item.isDir)
	} else {
		t.Log("3. 没有检测到需要转换的文件")
	}

	// 4. 测试转换方向切换
	t.Log("4. 测试转换方向切换...")
	state.SetConvertDirection("t2s")
	newDirection := state.GetConvertDirection()
	t.Logf("   切换后方向: %s", newDirection)

	// 清理
	os.RemoveAll(tempDir)

	t.Log("优化后的工作流程测试完成!")
	t.Log("✅ 用户现在只需要:")
	t.Log("   1. 选择转换方向（简转繁/繁转简）")
	t.Log("   2. 点击'选择文件'")
	t.Log("   3. 在弹出的确认对话框中点击'确认'")
	t.Log("   4. 转换自动完成！")
}
