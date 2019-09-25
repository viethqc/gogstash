package cmder

import (
	"os"

	"github.com/viethqc/gogstash/KDGoLib/errutil"
)

// Main entry point
func Main(module Module, commandModules ...Module) {
	app := NewApp(module, commandModules...)
	errutil.Trace(app.Run(os.Args))
}
