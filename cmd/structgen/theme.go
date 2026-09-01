package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dave/dst"
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

var acc, val []string

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

	acc = append(acc, fmt.Sprintf("r.%s", s))

	t.count++

	return "lipgloss.Style", true
}

func (t *ThemeGenImpl) GetPartialCode() string {
	return fmt.Sprintf(themeGenPart, t.count, t.sb.String())
}

func generateTheme() error {
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

	const (
		replaceVar = "_rootCfg"
	)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFilePath := filepath.Join(dir, configFileName)
	themeGenPath := filepath.Join(dir, themeGenFileName)

	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	var p []varReplaceOptions

	scanner := bufio.NewScanner(f)
	lineNum := 1

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		line := scanner.Text()

		if strings.Contains(line, replaceVar) {
			out, err := execute("go", "tool", "fillstruct", "-file", cfgFilePath, "-line", strconv.Itoa(lineNum))
			if err != nil {
				return err
			}

			err = json.Unmarshal(out, &p)
			if err != nil {
				return err
			}

			break
		}

		lineNum++
	}

	replOpts := p[0]
	replOpts.set(replaceVar, pkgName, cfgFilePath, func(n dst.Node) bool {
		if lit, ok := n.(*dst.BasicLit); ok {
			val = append(val, lit.Value)

			if lit.Value == `""` {
				lit.Kind = token.VAR
				lit.Value = `__UNDEFINED__`
			}
		}
		return true
	})

	if err := replaceStructVar(replOpts); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	defer func() {
		for i := range acc {
			fmt.Printf("%s = %s\n", acc[i], val[i])
		}
	}()

	return generateStruct(&ThemeGenImpl{}, &genOptions{
		pkgName:           pkgName,
		currentFilePath:   cfgFilePath,
		currentStructName: configStructName,
		newFilePath:       themeGenPath,
		newStructName:     themeStructName,
		newStructComment:  themeStructComment,
	})
}
