package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
)

const kbConstPart = `
const (
	_ KeyID = iota
	%s
)
`

const kbDefPart = `
func defaultConfig() *Keybindings {
	k := &Keybindings{}

	%s

	return k
}
`

const kbPart = `
package keybindings

%s

var _emptyCmpCfg = &Configuration{}

%s

type parseKeybindingInfo struct {
	kb *Keybinding
	cmpCfg  string
}

func iterProperties(kb *Keybindings, cfg *Configuration) iter.Seq[parseKeybindingInfo] {
	if cfg == nil {
		cfg = _emptyCmpCfg
	}

	return func(yield func(parseKeybindingInfo) bool) {
		for i := range %d {
			if !yield(getProperty(kb, cfg, i)) {
				return
			}
		}
	}
}

func getProperty(kb *Keybindings, cfg *Configuration, pos int) parseKeybindingInfo {
	p := parseKeybindingInfo{}

	switch pos {
	%s
	}

	return p
}
`

var tagKeysKb = [3]string{"keydef", "keyshorthelp", "keylonghelp"}

type KbGenImpl struct {
	sb, constSb, defSb strings.Builder
	count              int
}

func (k *KbGenImpl) AppendAccessor(s string, tag string) (genTemplRet, error) {
	tags, err := structtag.Parse(tag)
	if err != nil {
		return emptyGenTemplRet(), err
	}

	id := "KeyID" + strings.ReplaceAll(s, ".", "")

	var vals [len(tagKeysKb)]string
	for idx, k := range tagKeysKb {
		t, err := tags.Get(k)
		if err != nil {
			return emptyGenTemplRet(), fmt.Errorf(
				"%w: required key not found for accessor %s (key %s, tag %s)",
				err, s, k, tag,
			)
		}

		vals[idx] = t.Value()
	}

	fmt.Fprintf(&k.constSb, "%s\n", id)

	fmt.Fprintf(
		&k.defSb,
		`k.%s = NewKeybinding(
		%s, %q,
		%q,
		%q,
		)
		`,
		s, id, vals[0], vals[1], vals[2],
	)

	fmt.Fprintf(&k.sb, `
	case %d:
	p.kb = &kb.%s
	p.cmpCfg = cfg.%s

	`, k.count, s, s)

	k.count++

	return newGenTemplRet("Keybinding", true, true), nil
}

func (k *KbGenImpl) GetPartialCode() string {
	consts := fmt.Sprintf(kbConstPart, k.constSb.String())

	k.defSb.WriteString("\n")
	defs := fmt.Sprintf(kbDefPart, k.defSb.String())

	code := fmt.Sprintf(kbPart, consts, defs, k.count, k.sb.String())

	return code
}

func generateKeybindings() error {
	const (
		pkgName = "keybindings"
		kbDir   = "ui/keybindings"
	)

	const (
		configFileName   = "config.go"
		configStructName = "Configuration"
	)

	const (
		kbGenFileName   = "keybindings.gen.go"
		kbStructName    = "Keybindings"
		kbStructComment = "// Keybindings represents the settings for the app's keybindings."
	)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFilePath := filepath.Join(dir, configFileName)
	kbGenFilePath := filepath.Join(dir, kbGenFileName)
	k := &KbGenImpl{}

	return generateStruct(k, &genOptions{
		pkgName:           pkgName,
		currentFilePath:   cfgFilePath,
		currentStructName: configStructName,
		newFilePath:       kbGenFilePath,
		newStructName:     kbStructName,
		newStructComment:  kbStructComment,
	})
}
