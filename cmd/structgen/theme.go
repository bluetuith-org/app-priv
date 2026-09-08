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

var (
	_emptyCmpCfg = &Configuration{}
	_emptyTheme = &Theme{}
)

%s

// parseThemeInfo holds the theme parsing information.
type parseThemeInfo struct {
	rootCfg *string
  cmpCfg string

	style *lipgloss.Style
}

// iterProperties iterates over the configuration's and theme's properties.
func iterProperties(cfg *RootConfiguration, cmpCfg *Configuration, t *Theme) iter.Seq[parseThemeInfo] {
	if cmpCfg == nil {
	  cmpCfg = _emptyCmpCfg
	}

	if t == nil {
	  t = _emptyTheme
	}

	return func(yield func(parseThemeInfo) bool) {
		for i := range %d {
			if !yield(getProperty(cfg, cmpCfg, t, i)) {
				return
			}
		}
	}
}

// getProperty returns the property according to the specified parsing position.
func getProperty(cfg *RootConfiguration, cmpCfg *Configuration, t *Theme, pos int) parseThemeInfo {
	p := parseThemeInfo{}

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

func (t *ThemeGenImpl) AppendAccessor(s string, tag string) (genTemplRet, error) {
	isRoot := true

	k, _ := getAccessor(s)
	if k != "" {
		isRoot = false
	}

	tags, err := structtag.Parse(tag)
	if err != nil {
		return emptyGenTemplRet(), fmt.Errorf("%w: cannot parse tag on accessor %s (tag %s)", err, s, tag)
	}

	var (
		tagValue, tagKey string
		foundCount       int
	)

	for _, tkey := range tagKeys {
		tag, err := tags.Get(tkey)
		if err == nil {
			tagValue = tag.Value()
			tagKey = tkey

			foundCount++
		}
	}
	if tagValue == "" {
		return emptyGenTemplRet(), fmt.Errorf(
			"no tag value was found for accessor %s (no tags keys '%s' were found)",
			s, strings.Join(tagKeys[:], ", "),
		)
	}
	if foundCount >= len(tagKeys) {
		return emptyGenTemplRet(), fmt.Errorf(
			"only one of %s must be specified for accessor %s (tag %s)",
			strings.Join(tagKeys[:], ", "),
			s, tag,
		)
	}

	fmt.Fprintf(&t.defSb, "r.%s = %q\n", s, strings.ReplaceAll(tagValue, " ", ""))

	if isRoot && tagKey == tagKeys[0] {
		return newGenTemplRet("", false, true), nil
	}

	fmt.Fprintf(&t.sb, `
		case %d:
		p.rootCfg = &cfg.%s
		p.cmpCfg = cmpCfg.%s
		p.style = &t.%s

		`, t.count, s, s, s)

	t.count++

	return newGenTemplRet("lipgloss.Style", true, true), nil
}

func (t *ThemeGenImpl) GetPartialCode() string {
	t.defSb.WriteString("\n")
	def := fmt.Sprintf(defaultConfigPart, t.defSb.String())

	code := fmt.Sprintf(themeGenPart, def, t.count, t.sb.String())

	return code
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
