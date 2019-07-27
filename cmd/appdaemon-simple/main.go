package main

import (
	"github.com/kimkit/appdaemon/pkg/cmdsvr"
	"github.com/kimkit/appdaemon/pkg/common"
)

func main() {
	if common.Config.UI.Run {
		common.Logger.LogError("main.main", "`-ui` argument not support for simple version")
	} else {
		cmdsvr.Run()
	}
}
