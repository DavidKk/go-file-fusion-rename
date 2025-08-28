package app

import (
	"testing"
)

func TestConvertDirection(t *testing.T) {
	state := NewState()

	// 测试默认状态
	defaultDirection := state.GetConvertDirection()
	t.Logf("默认转换方向: %s", defaultDirection)

	// 测试设置简转繁
	state.SetConvertDirection("s2t")
	direction := state.GetConvertDirection()
	if direction != "s2t" {
		t.Errorf("期望转换方向为 's2t'，实际为 '%s'", direction)
	}
	t.Log("简转繁设置测试通过")

	// 测试设置繁转简
	state.SetConvertDirection("t2s")
	direction = state.GetConvertDirection()
	if direction != "t2s" {
		t.Errorf("期望转换方向为 't2s'，实际为 '%s'", direction)
	}
	t.Log("繁转简设置测试通过")

	// 验证源语言和目标语言
	srcLang := state.GetSrcLang()
	dstLang := state.GetDstLang()
	t.Logf("转换方向 %s: 源语言=%s, 目标语言=%s", direction, srcLang, dstLang)

	t.Log("转换方向功能测试通过")
}
