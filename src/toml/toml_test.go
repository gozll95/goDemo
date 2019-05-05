package conf

import "testing"

func TestReadConf(t *testing.T) {
	tt, err := ReadConf("./config.toml")
	if err != nil {
		t.Logf("%v", err)
	}
	t.Logf("result %v", tt)
	t.Logf("result %v", tt.Servers)
}
