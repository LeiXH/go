package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5Encode(t *testing.T) {
	_in := "0"
	expect := "cfcd208495d565ef66e7dff9f98764da"
	assert.Equal(t, expect, MD5Sum(_in), "md5sum failed")
}

func TestPseudo_uuid(t *testing.T) {
	t.Log(PseudoUUID())
}

func TestCountUnicodeString(t *testing.T) {
	assert.Equal(t, 2, CountUnicodeString("世界"), "count unicode-string length failed")
	assert.Equal(t, 10, CountUnicodeString("小牛牛 tester"), "count unicode-string length failed")
	assert.Equal(t, 6, CountUnicodeString("Tester"), "count unicode-string length failed")
}

func TestCutUnicodeString(t *testing.T) {
	assert.Equal(t, "世", CutUnicodeString("世界", 1), "cut unicode-string length failed")
	assert.Equal(t, "小牛牛 t", CutUnicodeString("小牛牛 tester", 5), "cut unicode-string length failed")
}

func TestCopyAttrs(t *testing.T) {
	a := struct {
		A string
		B string
	}{A: "1", B: "2"}

	b := struct {
		B string
		C string
	}{B: "3"}

	CopyAttrs(&a, &b)
	t.Log(a)
}

func TestIsAllAscii(t *testing.T) {
	assert.Equal(t, true, IsAllAscii("1 a B"), "ascii test failed")
	assert.Equal(t, false, IsAllAscii("哈哈 a"), "ascii test failed")
}
func TestIfHaveAscii(t *testing.T) {
	assert.Equal(t, true, IfHaveAscii("1 a B"), "ascii test failed")
	assert.Equal(t, true, IfHaveAscii("哈哈 a"), "ascii test failed")
	assert.Equal(t, false, IfHaveAscii("哈哈 "), "ascii test failed")
}
