package httpx

import "testing"

func TestMockClient(t *testing.T) {
	Mock()
	{
		const url = "https://antonz.org/example.txt"
		ok := Exists(url)
		if !ok {
			t.Errorf("Exists(%s) expected true, got false", url)
		}
	}
	{
		const url = "https://antonz.org/missing.txt"
		ok := Exists(url)
		if ok {
			t.Errorf("Exists(%s) expected false, got true", url)
		}
	}
	{
		const url = "https://antonz.org/example.txt"
		data, err := GetBytes(url)
		if err != nil {
			t.Errorf("GetBytes: unexpected error %v", err)
		}
		if string(data) != "example.txt" {
			t.Errorf("GetBytes: unexpected value %q", string(data))
		}
	}
	{
		const url = "https://antonz.org/missing.txt"
		_, err := GetBytes(url)
		if err == nil {
			t.Error("GetBytes: expected error, got nil")
		}
	}
}
