package whoapi

import (
	"encoding/json"
	"testing"
)

func TestGet(t *testing.T) {
	c := &Client{}
	data, err := c.Get("myaccount", "")
	if err != nil {
		t.Error(err)
	}
	var r map[string]json.RawMessage
	err = json.Unmarshal(data, &r)
	if err != nil {
		t.Error(err)
	}
	// Validate there's really no error
	if string(r["status"]) != `0` {
		t.Errorf("Unexpected status: %s", r["status"])
	}
}

func TestGetKeyError(t *testing.T) {
	c := &Client{Key: "_"}
	_, err := c.Get("myaccount", "")
	if err != ErrInvalidAPIAccount {
		t.Error("Expected ErrInvalidAPIAccount, got:", err)
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		data string
		ok   bool
		val  int
	}{
		{"", false, 0},
		{`""`, false, 0},
		{"0", true, 0},
		{`"0"`, true, 0},
		{"123", true, 123},
		{`"123"`, true, 123},
		{"-10", true, -10},
		{`"-10"`, true, -10},
		{`"1""`, false, 0},
	}
	for _, c := range cases {
		i := Int(1)
		err := i.UnmarshalJSON([]byte(c.data))
		if (err == nil) != c.ok || int(i) != c.val {
			t.Error(c, i, err)
		}
	}
}
