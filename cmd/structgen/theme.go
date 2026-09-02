package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
)

const themeGenPart = `
package theme

%s

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

const defaultConfigPart = `
func defaultConfig() *RootConfiguration {
	r := &RootConfiguration{}

	%s

	return r
}
`

var tagKeys = [2]string{"themetype", "themedef"}

type ThemeGenImpl struct {
	sb, defSb strings.Builder
	count     int
}

func getAccessor(s string) (k string, v string) {
	idx := strings.LastIndex(s, ".")
	if idx == -1 {
		return "", s
	}

	return s[:idx], s[idx:]
}

func (t *ThemeGenImpl) AppendAccessor(s string, tag string) (genTemplRet, error) {
	path := "root"
	isRoot := true

	k, v := getAccessor(s)
	if k != "" {
		path += "/" + strings.ReplaceAll(k, ".", "/")
		isRoot = false
	}
	path += ":" + v

	tags, err := structtag.Parse(tag[min(1, len(tag)):max(0, len(tag)-1)])
	if err != nil {
		return emptyGenTemplRet(), fmt.Errorf("%w: cannot parse tag on accessor %s (tag %s)", err, s, tag)
	}

	tagValue := ""
	for _, tagKey := range tagKeys {
		tag, err := tags.Get(tagKey)
		if err == nil {
			tagValue = tag.Value()
		}
	}
	if tagValue == "" {
		return emptyGenTemplRet(), fmt.Errorf(
			"no tag value was found for accessor %s (no tags keys '%s' were found)",
			s, strings.Join(tagKeys[:], ", "),
		)
	}

	fmt.Fprintf(&t.defSb, "r.%s = \"%s\"\n", s, tagValue)

	if isRoot && s == "Tint" {
		return newGenTemplRet("", false, true), nil
	}

	fmt.Fprintf(&t.sb, `
		case %d:
		p.path = "%s"
		p.rootCfg = &cfg.%s
		p.cmpCfg = cmpCfg.%s
		p.style = &t.%s

		`, t.count, path, s, s, s)

	t.count++

	return newGenTemplRet("lipgloss.Style", true, true), nil
}

func (t *ThemeGenImpl) GetPartialCode() string {
	t.defSb.WriteString("\n")
	def := fmt.Sprintf(defaultConfigPart, t.defSb.String())

	v := fmt.Sprintf(themeGenPart, def, t.count, t.sb.String())

	return v
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

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFilePath := filepath.Join(dir, configFileName)
	themeGenPath := filepath.Join(dir, themeGenFileName)
	t := &ThemeGenImpl{}

	return generateStruct(t, &genOptions{
		pkgName:           pkgName,
		currentFilePath:   cfgFilePath,
		currentStructName: configStructName,
		newFilePath:       themeGenPath,
		newStructName:     themeStructName,
		newStructComment:  themeStructComment,
	})
}
