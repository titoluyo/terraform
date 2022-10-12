package renderer

import (
	"github.com/hashicorp/terraform/internal/command/arguments"
)

type Renderer interface {
	Plan(plan Plan)

	Printf() (n int, err error)
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
