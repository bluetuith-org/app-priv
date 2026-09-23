package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/fatih/structtag"
	"github.com/gdamore/tcell/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const kbConstPart = `
const (
	KeyNone KeyID = iota
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

	id := "Key" + strings.ReplaceAll(s, ".", "")

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

	kb, err := genTcellKey(vals[0])
	if err != nil {
		return emptyGenTemplRet(), err
	}

	fmt.Fprintf(&k.constSb, "%s\n", id)

	fmt.Fprintf(
		&k.defSb,
		`k.%s = newKeybinding(
		%s, %s, // %s
		%q,
		%q,
		)
		`,
		s, id, kb, vals[0], vals[1], vals[2],
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

var translateKeys = map[string]string{
	"Pgup":      "PgUp",
	"Pgdown":    "PgDn",
	"Upright":   "UpRight",
	"Downright": "DownRight",
	"Upleft":    "UpLeft",
	"Downleft":  "DownLeft",
	"Prtsc":     "Print",
	"Backspace": "Backspace2",
}

type tcellKey struct {
	key tcell.Key
	str string
	mod tcell.ModMask
}

func genTcellKey(key string) (string, error) {
	var mods []string

	var tkey string
	var keyCount int
	var found bool

	kb := tcellKey{
		key: tcell.KeyRune,
		mod: tcell.ModNone,
	}

	tokens := strings.FieldsFuncSeq(key, func(c rune) bool {
		return unicode.IsSpace(c) || c == '+'
	})

	titleCase := cases.Title(language.Und, cases.NoLower)

	for token := range tokens {
		if keyCount > 1 {
			break
		}

		if len(token) == 1 {
			kb.str = token
			keyCount++

			continue
		}

		token = titleCase.String(token)

		if translated, ok := translateKeys[token]; ok {
			token = translated
		}

		switch token {
		case "Ctrl":
			kb.mod |= tcell.ModCtrl
			mods = append(mods, "tcell.ModCtrl")

		case "Alt":
			kb.mod |= tcell.ModAlt
			mods = append(mods, "tcell.ModAlt")

		case "Hyper":
			kb.mod |= tcell.ModHyper
			mods = append(mods, "tcell.ModHyper")

		case "Shift":
			kb.mod |= tcell.ModShift
			mods = append(mods, "tcell.ModShift")

		case "Space", "Plus":
			kb.str = " "
			if token == "Plus" {
				kb.str = "+"
			}

			keyCount++

		default:
			tkey = token
			keyCount++
		}
	}

	if keyCount > 1 {
		return "", fmt.Errorf("config: More than one key entered for %s (%s)", "keybinding", key)
	}

	if kb.mod == tcell.ModCtrl && kb.str != "" && kb.str != " " && len(kb.str) == 1 {
		tkey = "Ctrl-" + strings.ToUpper(kb.str)
	}

	if tkey == "" {
		goto Print
	}

	for tk, ts := range tcell.KeyNames {
		if ts == tkey {
			kb.key = tk
			if kb.mod == tcell.ModCtrl {
				kb.str = ""
			}

			found = true

			break
		}
	}

	if !found {
		return "", fmt.Errorf("config: Invalid keybinding key %s (%s)", tkey, key)
	}

Print:
	if keyCount == 0 {
		return "", fmt.Errorf("config: No key specified or invalid keybinding for %s (%s)", "keybinding", key)
	}

	modStr := "tcell.ModNone"
	if kb.mod != 0 {
		modStr = strings.Join(mods, "|")
	}

	return fmt.Sprintf("newTcellKey(tcell.Key(%d), %q, %s)", kb.key, kb.str, modStr), nil
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
