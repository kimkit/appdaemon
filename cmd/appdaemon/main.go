package main

import (
	"github.com/kimkit/appdaemon/pkg/apisvr"
	"github.com/kimkit/appdaemon/pkg/cmdsvr"
	"github.com/kimkit/appdaemon/pkg/common"
)

func main() {
	if common.Config.UI.Run {
		apisvr.Run()
	} else {
		cmdsvr.Run()
	}
}
