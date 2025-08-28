package menu

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"file-fusion-rename/internal/app"
	"file-fusion-rename/internal/ui/dialog"
	"file-fusion-rename/internal/zhconv"

	"github.com/ncruces/zenity"
)

// 文件夹重命名器
type Renamer struct {
	state  *app.State
	dialog *dialog.Manager
}

// 创建重命名器
func NewRenamer(state *app.State, dialogMgr *dialog.Manager) *Renamer {
	return &Renamer{
		state:  state,
		dialog: dialogMgr,
	}
}

// 显示重命名确认对话框
func (r *Renamer) ShowConfirmation() error {
	directionText := app.GetDirectionText(r.state.GetConvertDirection())
	pathName := filepath.Base(r.state.GetFilePath())
	message := "确认要转换吗？\n\n路径: " + pathName + "\n转换方向: " + directionText + "\n\n将同时转换文件名和文件内容"

	return zenity.Question(message,
		zenity.Title("确认转换"),
		zenity.OKLabel("开始转换"),
		zenity.CancelLabel("取消"),
	)
}

// 为指定路径显示确认对话框
func (r *Renamer) ShowConfirmationForPath(path string) error {
	directionText := app.GetDirectionText(r.state.GetConvertDirection())
	pathName := filepath.Base(path)
	message := "确认要转换吗？\n\n路径: " + pathName + "\n转换方向: " + directionText + "\n\n将同时转换文件名和文件内容"

	return zenity.Question(message,
		zenity.Title("确认转换"),
		zenity.OKLabel("开始转换"),
		zenity.CancelLabel("取消"),
	)
}

// 执行实际的重命名操作
func (r *Renamer) PerformRename() {
	path := r.state.GetFilePath()
	srcLang := r.state.GetSrcLang()
	dstLang := r.state.GetDstLang()

	info, err := os.Stat(path)
	if err != nil {
		r.dialog.ShowError("读取路径信息失败: " + err.Error())
		return
	}

	err = zhconv.ConvertPath(path, srcLang, dstLang)
	if err != nil {
		r.dialog.ShowError("转换文件内容失败: " + err.Error())
		return
	}

	var renameList []RenameItem

	if info.IsDir() {
		renameList, err = r.collectRenameItemsFromFolder(path, srcLang, dstLang)
	} else {
		renameList, err = r.collectRenameItemsFromFile(path, srcLang, dstLang)
	}

	if err != nil {
		r.dialog.ShowError("扫描文件失败: " + err.Error())
		return
	}

	finalPath := path
	for _, item := range renameList {
		err := os.Rename(item.oldPath, item.newPath)
		if err != nil {
			r.dialog.ShowError("重命名失败: " + item.oldPath + " -> " + err.Error())
			continue
		}
		if item.oldPath == path {
			finalPath = item.newPath
		}
	}

	r.openFolder(finalPath)
}

// 为指定路径执行重命名操作
func (r *Renamer) PerformRenameForPath(path string) {
	srcLang := r.state.GetSrcLang()
	dstLang := r.state.GetDstLang()

	info, err := os.Stat(path)
	if err != nil {
		r.dialog.ShowError("读取路径信息失败: " + err.Error())
		return
	}

	err = zhconv.ConvertPath(path, srcLang, dstLang)
	if err != nil {
		r.dialog.ShowError("转换文件内容失败: " + err.Error())
		return
	}

	var renameList []RenameItem

	if info.IsDir() {
		renameList, err = r.collectRenameItemsFromFolder(path, srcLang, dstLang)
	} else {
		renameList, err = r.collectRenameItemsFromFile(path, srcLang, dstLang)
	}

	if err != nil {
		r.dialog.ShowError("扫描文件失败: " + err.Error())
		return
	}

	finalPath := path
	for _, item := range renameList {
		err := os.Rename(item.oldPath, item.newPath)
		if err != nil {
			r.dialog.ShowError("重命名失败: " + item.oldPath + " -> " + err.Error())
			continue
		}
		if item.oldPath == path {
			finalPath = item.newPath
		}
	}

	r.openFolder(finalPath)
}

// RenameItem 重命名项目
type RenameItem struct {
	oldPath string
	newPath string
	isDir   bool
}

// 收集单个文件的重命名项目
func (r *Renamer) collectRenameItemsFromFile(filePath string, srcLang, dstLang zhconv.Language) ([]RenameItem, error) {
	var items []RenameItem

	if srcLang == "" {
		fileName := filepath.Base(filePath)
		srcLang = zhconv.DetectLanguage(fileName)
		if srcLang == "" {
			return items, nil
		}
	}

	oldName := filepath.Base(filePath)

	newName, err := zhconv.ConvertString(oldName, srcLang, dstLang)
	if err != nil {
		return items, nil
	}

	if newName == oldName {
		return items, nil
	}

	dir := filepath.Dir(filePath)
	newPath := filepath.Join(dir, newName)

	items = append(items, RenameItem{
		oldPath: filePath,
		newPath: newPath,
		isDir:   false,
	})

	return items, nil
}

// 从文件夹收集重命名项目
func (r *Renamer) collectRenameItemsFromFolder(folderPath string, srcLang, dstLang zhconv.Language) ([]RenameItem, error) {
	if srcLang == "" {
		srcLang = r.detectLanguageFromFolder(folderPath)
		if srcLang == "" {
			return nil, fmt.Errorf("无法检测文件夹内文件名的语言，请手动选择源语言")
		}
	}

	items, err := r.collectRenameItems(folderPath, srcLang, dstLang)
	if err != nil {
		return nil, err
	}

	oldName := filepath.Base(folderPath)
	newName, convErr := zhconv.ConvertString(oldName, srcLang, dstLang)
	if convErr == nil && newName != oldName {
		dir := filepath.Dir(folderPath)
		newPath := filepath.Join(dir, newName)
		rootItem := RenameItem{
			oldPath: folderPath,
			newPath: newPath,
			isDir:   true,
		}
		items = append(items, rootItem)
	}

	return items, nil
}

// collectRenameItems 收集需要重命名的项目
// 重要：先处理所有文件，再处理文件夹，避免路径变化问题
func (r *Renamer) collectRenameItems(folderPath string, srcLang, dstLang zhconv.Language) ([]RenameItem, error) {
	var files []RenameItem
	var folders []RenameItem

	// 第一次遍历：收集所有文件
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == folderPath {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		oldName := info.Name()

		newName, convErr := zhconv.ConvertString(oldName, srcLang, dstLang)
		if convErr != nil {
			return nil
		}

		if newName == oldName {
			return nil
		}

		dir := filepath.Dir(path)
		newPath := filepath.Join(dir, newName)

		files = append(files, RenameItem{
			oldPath: path,
			newPath: newPath,
			isDir:   false,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 第二次遍历：收集所有文件夹
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == folderPath {
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		oldName := info.Name()

		newName, convErr := zhconv.ConvertString(oldName, srcLang, dstLang)
		if convErr != nil {
			return nil
		}

		if newName == oldName {
			return nil
		}

		dir := filepath.Dir(path)
		newPath := filepath.Join(dir, newName)

		folders = append(folders, RenameItem{
			oldPath: path,
			newPath: newPath,
			isDir:   true,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(folders, func(i, j int) bool {
		return len(strings.Split(folders[i].oldPath, string(filepath.Separator))) >
			len(strings.Split(folders[j].oldPath, string(filepath.Separator)))
	})

	var result []RenameItem
	result = append(result, files...)
	result = append(result, folders...)

	return result, nil
}

// 从文件夹中检测语言
func (r *Renamer) detectLanguageFromFolder(folderPath string) zhconv.Language {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return ""
	}

	for i, entry := range entries {
		if i >= 5 {
			break
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		detected := zhconv.DetectLanguage(name)
		if detected != "" {
			return detected
		}
	}

	return ""
}

// 在不同操作系统中打开文件夹
func (r *Renamer) openFolder(path string) {
	var cmd *exec.Cmd

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	var targetPath string
	if info.IsDir() {
		targetPath = path
	} else {
		targetPath = filepath.Dir(path)
	}

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", targetPath)
	case "windows":
		cmd = exec.Command("explorer", targetPath)
	case "linux":
		cmd = exec.Command("xdg-open", targetPath)
	default:
		return
	}

	cmd.Start()
}
