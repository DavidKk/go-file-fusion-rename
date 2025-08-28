package app

import (
	"file-fusion-rename/internal/zhconv"
)

// State 应用程序状态管理
type State struct {
	convertDirection string
	filePath         string
}

// 创建新的应用状态
func NewState() *State {
	return &State{
		convertDirection: "t2s",
	}
}

// 获取转换方向
func (s *State) GetConvertDirection() string {
	return s.convertDirection
}

// 设置转换方向
func (s *State) SetConvertDirection(direction string) {
	s.convertDirection = direction
}

// 获取源语言
func (s *State) GetSrcLang() zhconv.Language {
	switch s.convertDirection {
	case "s2t":
		return zhconv.LangSimplified
	case "t2s":
		return zhconv.LangTraditional
	default:
		return zhconv.LangSimplified
	}
}

// 获取目标语言
func (s *State) GetDstLang() zhconv.Language {
	switch s.convertDirection {
	case "s2t":
		return zhconv.LangTraditional
	case "t2s":
		return zhconv.LangSimplified
	default:
		return zhconv.LangTraditional
	}
}

// 获取文件路径
func (s *State) GetFilePath() string {
	return s.filePath
}

// 设置文件路径
func (s *State) SetFilePath(path string) {
	s.filePath = path
}

// 获取转换方向的显示文本
func GetDirectionText(direction string) string {
	switch direction {
	case "s2t":
		return "简转繁"
	case "t2s":
		return "繁转简"
	default:
		return "繁转简"
	}
}
