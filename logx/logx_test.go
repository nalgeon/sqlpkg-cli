package logx

import "testing"

func TestSetOutput(t *testing.T) {
	mem := NewMemory("log")
	SetOutput(mem)
	if logger.out != mem {
		t.Errorf("SetOutput: unexpected value %v", logger.out)
	}
}

func TestVerbose(t *testing.T) {
	mem := NewMemory("log")
	SetOutput(mem)
	{
		if IsVerbose() {
			t.Errorf("IsVerbose: expected false, got true")
		}
	}
	{
		SetVerbose(true)
		if !IsVerbose() {
			t.Errorf("IsVerbose: expected true, got false")
		}
	}
	{
		SetVerbose(false)
		if IsVerbose() {
			t.Errorf("IsVerbose: expected false, got true")
		}
	}
}

func TestLog(t *testing.T) {
	mem := NewMemory("log")
	SetOutput(mem)
	{
		Log("value: %d", 42)
		if len(mem.Lines) != 1 {
			t.Errorf("Log: expected line count %v", len(mem.Lines))
		}
		if !mem.Has("value: 42") {
			t.Errorf("Log: expected output: %v", mem.Lines)
		}
	}
	{
		Log("value: %d", 84)
		if len(mem.Lines) != 2 {
			t.Errorf("Log: expected line count %v", len(mem.Lines))
		}
		if !mem.Has("value: 42") || !mem.Has("value: 84") {
			t.Errorf("Log: expected output: %v", mem.Lines)
		}
	}
}

func TestDebug(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		mem := NewMemory("log")
		SetOutput(mem)
		SetVerbose(true)
		{
			Debug("value: %d", 42)
			if len(mem.Lines) != 1 {
				t.Errorf("Log: expected line count %v", len(mem.Lines))
			}
			if !mem.Has("value: 42") {
				t.Errorf("Log: expected output: %v", mem.Lines)
			}
		}
		{
			Debug("value: %d", 84)
			if len(mem.Lines) != 2 {
				t.Errorf("Log: expected line count %v", len(mem.Lines))
			}
			if !mem.Has("value: 42") || !mem.Has("value: 84") {
				t.Errorf("Log: expected output: %v", mem.Lines)
			}
		}
	})
	t.Run("disabled", func(t *testing.T) {
		mem := NewMemory("log")
		SetOutput(mem)
		SetVerbose(false)
		Debug("value: %d", 42)
		if len(mem.Lines) != 0 {
			t.Errorf("Log: expected line count %v", len(mem.Lines))
		}
	})
}
