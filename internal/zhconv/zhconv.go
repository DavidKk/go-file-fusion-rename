package zhconv

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"unicode/utf8"

	"github.com/longbridgeapp/opencc"
)

// Language represents a Chinese script.
// Supported values:
//   - "s": Simplified Chinese
//   - "t": Traditional Chinese
type Language string

const (
	LangSimplified  Language = "s"
	LangTraditional Language = "t"
)

var (
	initOnce sync.Once
	convT2S  *opencc.OpenCC
	convS2T  *opencc.OpenCC
	initErr  error
)

func ensureConverters() error {
	initOnce.Do(func() {
		var err error
		convT2S, err = opencc.New("t2s")
		if err != nil {
			initErr = err
			return
		}
		convS2T, err = opencc.New("s2t")
		if err != nil {
			initErr = err
			return
		}
	})
	return initErr
}

// DetectLanguage attempts to determine whether the given text is primarily
// Simplified or Traditional Chinese. Returns "s", "t", or empty string if unknown.
func DetectLanguage(text string) Language {
	if err := ensureConverters(); err != nil {
		return ""
	}

	// Heuristic: apply both conversions and count how many runes differ from
	// the original. The direction that results in fewer changes likely
	// corresponds to the source script of the input text.
	toS, errS := convT2S.Convert(text)
	toT, errT := convS2T.Convert(text)
	if errS != nil || errT != nil {
		return ""
	}

	dS := runeDiffCount(text, toS)
	dT := runeDiffCount(text, toT)

	if dS < dT {
		return LangSimplified
	}
	if dT < dS {
		return LangTraditional
	}
	return ""
}

func runeDiffCount(a, b string) int {
	// Compare rune by rune to handle different lengths safely
	ar := []rune(a)
	br := []rune(b)
	n := len(ar)
	if len(br) > n {
		n = len(br)
	}
	diff := 0
	for i := 0; i < n; i++ {
		var ra, rb rune
		if i < len(ar) {
			ra = ar[i]
		}
		if i < len(br) {
			rb = br[i]
		}
		if ra != rb {
			diff++
		}
	}
	return diff
}

// ConvertString converts text from src to dst ("t"<->"s"). If src==dst, the
// original text is returned. Only Simplified/Traditional are supported.
func ConvertString(text string, src, dst Language) (string, error) {
	if src == dst {
		return text, nil
	}

	if err := ensureConverters(); err != nil {
		return "", err
	}

	switch {
	case src == LangTraditional && dst == LangSimplified:
		return convT2S.Convert(text)
	case src == LangSimplified && dst == LangTraditional:
		return convS2T.Convert(text)
	default:
		return "", errors.New("unsupported language conversion; only 't'<->'s' supported")
	}
}

// ConvertPath walks the directory tree rooted at rootPath and converts files
// whose detected language matches src into dst. Only UTF-8 text files are
// processed. Files that do not change after conversion are left untouched.
func ConvertPath(rootPath string, src, dst Language) error {
	if src == dst {
		return nil
	}

	if err := ensureConverters(); err != nil {
		return err
	}

	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Read file content
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		if !utf8.Valid(data) {
			// Skip non-UTF8 files
			return nil
		}

		original := string(data)
		detected := DetectLanguage(original)
		if detected != src {
			return nil
		}

		converted, convErr := ConvertString(original, src, dst)
		if convErr != nil {
			return convErr
		}

		if converted == original {
			return nil
		}

		// Preserve permissions
		info, statErr := d.Info()
		if statErr != nil {
			return statErr
		}

		// Only write when content actually changed
		if !bytes.Equal([]byte(converted), data) {
			if writeErr := os.WriteFile(path, []byte(converted), info.Mode()); writeErr != nil {
				return writeErr
			}
		}

		return nil
	})
}
