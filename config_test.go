package pazuzu

import (
	"strings"
	"testing"
)

// TestConfigSet : tests on its mutability.
func TestConfigSet(t *testing.T) {
	errCnf := NewConfig()
	if errCnf != nil {
		t.Errorf("%v", errCnf)
	}

	config := GetConfig()
	if len(config.Base) <= 0 {
		t.Error("Please fill 'Base' property of initial config.")
	}

	beforeBase := config.Base
	const someAnotherBase = "foo-bar-zoo-spam-eggs"
	config.SetBase(someAnotherBase)

	if strings.Compare(config.Base, someAnotherBase) != 0 {
		t.Errorf("Unexpected value [%v]", config.Base)
	}

	if strings.Compare(config.Base, beforeBase) == 0 {
		t.Errorf("Not changed 'Base' value!")
	}
}

// TODO: save
func TestConfigSave(ctx *testing.T) {
}

// TODO: load
func TestConfigLoad(ctx *testing.T) {
}
