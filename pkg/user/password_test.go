package user

import (
	"testing"
)

func TestAddSalt(t *testing.T) {
	enStr := "87k1uh9GZDApe8GSlr/fuxHa8Wkj0/vNhtG93nO/+Qs="
	r := addSalt("123456")
	if r != enStr {
		t.Errorf("want: <%s>\nget : <%s>\n", enStr, r)
	}
}
