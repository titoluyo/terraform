package terraform

import (
	"errors"

	"github.com/hashicorp/terraform/internal/configs/configschema"
	"github.com/hashicorp/terraform/internal/providers"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func nullResourceSchema() providers.Schema {
	return providers.Schema{
		Block: &configschema.Block{
			Attributes: map[string]*configschema.Attribute{
				"input":   {Type: cty.DynamicPseudoType, Optional: true},
				"output":  {Type: cty.DynamicPseudoType, Computed: true},
				"replace": {Type: cty.DynamicPseudoType, Optional: true},
			},
		},
	}
}

func validateNullResourceConfig(req providers.ValidateResourceConfigRequest) (resp providers.ValidateResourceConfigResponse) {
	if req.Config.IsNull() {
		return resp
	}

	input := req.Config.GetAttr("output")
	if !input.IsNull() {
		resp.Diagnostics = resp.Diagnostics.Append(errors.New(`"output" attribute is read-only`))
	}
	return resp
}

func upgradeNullResourceState(req providers.UpgradeResourceStateRequest) (resp providers.UpgradeResourceStateResponse) {
	ty := nullResourceSchema().Block.ImpliedType()
	val, err := ctyjson.Unmarshal(req.RawStateJSON, ty)
	if err != nil {
		resp.Diagnostics = resp.Diagnostics.Append(err)
		return resp
	}

	resp.UpgradedState = val
	return resp
}

func readNullResourceState(req providers.ReadResourceRequest) (resp providers.ReadResourceResponse) {
	resp.NewState = req.PriorState
	return resp
}

func planNullResourceChange(req providers.PlanResourceChangeRequest) (resp providers.PlanResourceChangeResponse) {
	if req.ProposedNewState.IsNull() {
		// destroy op
		resp.PlannedState = req.ProposedNewState
		return resp
	}

	planned := req.ProposedNewState.AsValueMap()

	input := req.ProposedNewState.GetAttr("input")
	replace := req.ProposedNewState.GetAttr("replace")

	switch {
	case req.PriorState.IsNull():
		planned["output"] = cty.UnknownVal(input.Type())
		resp.PlannedState = cty.ObjectVal(planned)
		return resp

	case !req.PriorState.GetAttr("input").RawEquals(input):
		planned["output"] = cty.UnknownVal(input.Type())

	case !req.PriorState.GetAttr("replace").RawEquals(replace):
		resp.RequiresReplace = append(resp.RequiresReplace, cty.GetAttrPath("replace"))
		planned["output"] = cty.UnknownVal(input.Type())
	}

	resp.PlannedState = cty.ObjectVal(planned)
	return resp
}

func applyNullResourceChange(req providers.ApplyResourceChangeRequest) (resp providers.ApplyResourceChangeResponse) {
	if req.PlannedState.IsNull() {
		resp.NewState = req.PlannedState
		return resp
	}

	newState := req.PlannedState.AsValueMap()
	output := req.PlannedState.GetAttr("output")
	if !output.IsNull() && !output.IsKnown() {
		newState["output"] = req.PlannedState.GetAttr("input")
	}

	resp.NewState = cty.ObjectVal(newState)

	return resp
}
