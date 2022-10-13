package renderer

import (
	"github.com/hashicorp/terraform/internal/command/jsonplan"
	"github.com/hashicorp/terraform/internal/command/jsonprovider"
)

type human struct {
	opts Opts
}

func (h *human) Plan(plan Plan) {
	//TODO implement me
	panic("implement me")
}

func (h *human) plan(plan Plan) {

}

func (h *human) refresh(plan Plan) (rendered bool) {
	var changes []jsonplan.ResourceChange
	for _, drift := range plan.ResourceDrift {
		schema := plan.ProviderSchema[drift.ProviderName]
	}
}

func (h *human) resourceChange(change jsonplan.ResourceChange, schema jsonprovider.Block) string {

	switch change.Change.Actions {
	case []string{"create"}:
		
	}

}
