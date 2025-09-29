package main

import (
	"github.com/Hahn814/magos/cmd/cli/cmd"
	_ "github.com/Hahn814/magos/cmd/cli/cmd/get"
	_ "github.com/Hahn814/magos/cmd/cli/cmd/get/agent"
	_ "github.com/Hahn814/magos/cmd/cli/cmd/get/agents"
)

func main() {
	cmd.Execute()
}
