package command

import "testing"

func TestDeployCommandExists(t *testing.T) {
	if _, ok := Commands["deploy"]; !ok {
		t.Fatal("Command 'deploy' should exist in Commands map")
	}
}
