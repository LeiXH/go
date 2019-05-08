package types

import (
	"encoding/json"
	"testing"
)

func TestSQLNullString_MarshalJSON(t *testing.T) {
	a := struct {
		A SQLNullString `json:"a"`
	}{A: struct {
		String string
		Valid  bool
	}{String: "666", Valid: true}}
	s, err := json.Marshal(a)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(s))
}
