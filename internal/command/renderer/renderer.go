package renderer

import (
	"github.com/hashicorp/terraform/internal/command/arguments"
)

type Hook string

const (
	Interrupted Hook = "Interrupted"
	FatalInterrupt Hook = "FatalInterrupt"
)

type Renderer interface {
	Plan(plan Plan)

	Log(format string) (n int, err error)
	Logf(format string, args ...any) (n int, err error)

	ErrorLog(format string) (n int, err error)
	ErrorLogf(format string, args ...any) (n int, err error)

	Colorize(v string) string
}

func New(viewType arguments.ViewType, opts Opts) Renderer {
	switch viewType {
	case arguments.ViewHuman:
		return &human{opts}
	case arguments.ViewJSON:
		return &json{opts}
	default:
		panic("unrecognized view type " + viewType.String())
	}
}
