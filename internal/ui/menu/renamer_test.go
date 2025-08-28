//go:build !nogui
// +build !nogui

package menu

import (
	"os"
	"path/filepath"
	"testing"

	"file-fusion-rename/internal/app"
	"file-fusion-rename/internal/ui/dialog"
	"file-fusion-rename/internal/zhconv"
)

func TestRenameOrder(t *testing.T) {
	// 创建临时测试目录
	tempDir := filepath.Join(os.TempDir(), "linkey_test")
	err := os.RemoveAll(tempDir)
	if err != nil {
		t.Fatalf("清理临时目录失败: %v", err)
	}

	// 创建测试文件夹结构
	testStructure := []string{
		"测试文件夹/简体文件.txt",
		"测试文件夹/子文件夹/深层简体文件.txt",
		"测试文件夹/另一个子文件夹/另一个简体文件.txt",
	}

	for _, path := range testStructure {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}

		if filepath.Ext(path) != "" { // 如果是文件
			file, err := os.Create(fullPath)
			if err != nil {
				t.Fatalf("创建文件失败: %v", err)
			}
			file.Close()
		}
	}

	// 创建重命名器
	state := app.NewState()
	state.SetConvertDirection("s2t") // 简转繁
	dialog := dialog.NewManager()
	renamer := NewRenamer(state, dialog)

	// 测试收集重命名项目
	items, err := renamer.collectRenameItems(tempDir, zhconv.LangSimplified, zhconv.LangTraditional)
	if err != nil {
		t.Fatalf("收集重命名项目失败: %v", err)
	}

	// 验证顺序：文件应该在文件夹之前
	var fileCount int
	foundFirstFolder := false

	for i, item := range items {
		if item.isDir {
			if !foundFirstFolder {
				foundFirstFolder = true
			}
		} else {
			fileCount++
			if foundFirstFolder {
				t.Errorf("发现文件在文件夹之后: %s (位置 %d)，这会导致路径变化问题", item.oldPath, i)
			}
		}
	}

	t.Logf("共收集到 %d 个重命名项目", len(items))
	t.Logf("文件数量: %d，文件夹数量: %d", fileCount, len(items)-fileCount)

	// 打印重命名顺序以供检查
	t.Log("重命名顺序:")
	for i, item := range items {
		itemType := "文件"
		if item.isDir {
			itemType = "文件夹"
		}
		t.Logf("%d. %s (%s): %s -> %s", i+1, itemType, item.oldPath, filepath.Base(item.oldPath), filepath.Base(item.newPath))
	}

	// 清理
	os.RemoveAll(tempDir)
}
