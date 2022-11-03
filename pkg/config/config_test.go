package config

import "testing"

func Test_NewRecordConfig(t *testing.T) {
	target := "ws://some-host:8080"
	file := "./some/path.gob"
	config, err := NewRecordConfig(&target, 30, &file)

	if err != nil {
		t.Error("NewRecordConfig failed unexpectedly", err)
	}
	if config.Target != target {
		t.Errorf("config.Target = %s; expected %s", config.Target, target)
	}
	if config.File != file {
		t.Errorf("config.File = %s; expected %s", config.File, file)
	}
	if config.Duration != 30 {
		t.Errorf("config.Duration = %d; expected %d", config.Duration, 30)
	}
}

func Test_NewRecordConfig_Empty_Target(t *testing.T) {
	target := ""
	file := "./some/path.gob"
	config, err := NewRecordConfig(&target, 30, &file)

	if config != nil {
		t.Error("NewRecordConfig with missing data produced a struct. Expected nil")
	}

	if err != ErrMissingTargetParam {
		t.Errorf("Expected %v; got %v", ErrMissingTargetParam, err)
	}
}

func Test_NewRecordConfig_Empty_File(t *testing.T) {
	target := "ws://localhost:8000"
	file := ""
	config, err := NewRecordConfig(&target, 30, &file)

	if config != nil {
		t.Error("NewRecordConfig with missing data produced a struct. Expected nil")
	}

	if err != ErrMissingFileParam {
		t.Errorf("Expected %v; got %v", ErrMissingFileParam, err)
	}
}

func Test_NewPlaybackConfig(t *testing.T) {
	addr := ":8001"
	file := "./some/path.gob"
	config, err := NewPlaybackConfig(&file, &addr)

	if err != nil {
		t.Error("NewPlaybackConfig failed unexpectedly", err)
	}
	if config.File != file {
		t.Errorf("config.File = %s; expected %s", config.File, file)
	}
	if config.ServerAddr != addr {
		t.Errorf("config.ServerAddr = %s; expected %s", config.ServerAddr, addr)
	}
}

func Test_NewPlaybackConfig_Empty_File(t *testing.T) {
	addr := ""
	file := "./some/path.gob"
	config, err := NewPlaybackConfig(&addr, &file)

	if config != nil {
		t.Error("NewPlaybackConfig with missing data produced a struct. Expected nil")
	}

	if err != ErrMissingFileParam {
		t.Errorf("Expected %v; got %v", ErrMissingFileParam, err)
	}
}

func Test_NewPlaybackConfig_Empty_Addr(t *testing.T) {
	addr := ":8001"
	file := ""
	config, err := NewPlaybackConfig(&addr, &file)

	if config != nil {
		t.Error("NewPlaybackConfig with missing data produced a struct. Expected nil")
	}

	if err != ErrMissingServerParam {
		t.Errorf("Expected %v; got %v", ErrMissingServerParam, err)
	}
}
