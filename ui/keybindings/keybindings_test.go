package keybindings

import (
	"fmt"
	"io"
	"testing"
)

func init() {
}

func BenchmarkKbMapIterator(b *testing.B) {
	m := make(map[*Keybinding]string)

	kb := defaultConfig()
	for p := range iterProperties(kb, nil) {
		m[p.kb] = p.cmpCfg
	}

	defer test()

	for b.Loop() {
		for k, v := range m {
			kbKey = k.Key
			cmpStr = v
		}
	}
}

var (
	kbKey  string
	cmpStr string
)

func BenchmarkKbStructIterator(b *testing.B) {
	kb := defaultConfig()
	defer test()

	it := iterProperties(kb, nil)

	for b.Loop() {
		for p := range it {
			kbKey = p.kb.Key
			cmpStr = p.cmpCfg
		}
	}
}

func test() {
	fmt.Fprintln(io.Discard, kbKey, cmpStr)
}
