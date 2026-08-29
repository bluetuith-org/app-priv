package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	pkgName  = "theme"
	themeDir = "ui/theme"
)

const (
	configFileName   = "config.go"
	configStructName = "Configuration"
)

const (
	themeGenFileName   = "theme.gen.go"
	themeStructName    = "Theme"
	themeStructComment = "// Theme represents the settings for the app's theme."
)

const themeGenPart = `
package theme

var (
	_emptyCmpCfg = &Configuration{}
	_emptyTheme = &Theme{}
)

// parseThemeInfo holds the theme parsing information.
type parseThemeInfo struct {
	path string

	rootCfg *string
  cmpCfg string

	style *lipgloss.Style
}

// iterProperties iterates over the configuration's and theme's properties.
func iterProperties(cfg *RootConfiguration, cmpCfg *Configuration, t *Theme) iter.Seq[*parseThemeInfo] {
	if cmpCfg == nil {
	  cmpCfg = _emptyCmpCfg
	}

	if t == nil {
	  t = _emptyTheme
	}

	return func(yield func(*parseThemeInfo) bool) {
		parseInfo := &parseThemeInfo{}

		for i := range %d {
			if !yield(getProperty(cfg, cmpCfg, t, parseInfo, i)) {
				return
			}
		}
	}
}

// getProperty returns the property according to the specified parsing position.
func getProperty(cfg *RootConfiguration, cmpCfg *Configuration, t *Theme, p *parseThemeInfo, pos int) *parseThemeInfo {
	switch pos {
	%s
	}

	return p
}
`

type ThemeGenImpl struct {
	sb    strings.Builder
	count int
}

func (t *ThemeGenImpl) AppendAccessor(s string) (string, bool) {
	path := "root"
	isRoot := true

	treePath := strings.ReplaceAll(s, ".", "/")
	pathFrag := filepath.Dir(treePath)
	if pathFrag != "." {
		path += "/" + pathFrag
		isRoot = false
	}

	path += ":" + filepath.Base(treePath)

	if isRoot && s == "Tint" {
		return "", false
	}

	fmt.Fprintf(&t.sb, `
		case %d:
		p.path = "%s"
		p.rootCfg = &cfg.%s
		p.cmpCfg = cmpCfg.%s
		p.style = &t.%s

		`, t.count, path, s, s, s)

	t.count++

	return "lipgloss.Style", true
}

func (t *ThemeGenImpl) GetPartialCode() string {
	return fmt.Sprintf(themeGenPart, t.count, t.sb.String())
}

func generateTheme() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFilePath := filepath.Join(dir, configFileName)
	themeGenPath := filepath.Join(dir, themeGenFileName)

	return generateStruct(&ThemeGenImpl{}, &GenOptions{
		pkgName:           pkgName,
		currentFilePath:   cfgFilePath,
		currentStructName: configStructName,
		newFilePath:       themeGenPath,
		newStructName:     themeStructName,
		newStructComment:  themeStructComment,
	})
}
