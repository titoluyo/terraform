package terraform

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/internal/providers"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func TestManagedDataValidate(t *testing.T) {
	cfg := map[string]cty.Value{
		"input":   cty.NullVal(cty.DynamicPseudoType),
		"output":  cty.NullVal(cty.DynamicPseudoType),
		"trigger": cty.NullVal(cty.DynamicPseudoType),
		"uuid":    cty.NullVal(cty.String),
	}

	// empty
	req := providers.ValidateResourceConfigRequest{
		TypeName: "terraform_box",
		Config:   cty.ObjectVal(cfg),
	}

	resp := validateDataResourceConfig(req)
	if resp.Diagnostics.HasErrors() {
		t.Error("empty config error:", resp.Diagnostics.ErrWithWarnings())
	}

	// invalid computed values
	cfg["output"] = cty.StringVal("oops")
	req.Config = cty.ObjectVal(cfg)

	resp = validateDataResourceConfig(req)
	if !resp.Diagnostics.HasErrors() {
		t.Error("expected error")
	}

	msg := resp.Diagnostics.Err().Error()
	if !strings.Contains(msg, "attribute is read-only") {
		t.Error("unexpected error", msg)
	}
}

func TestManagedDataUpgradeState(t *testing.T) {
	schema := dataResourceSchema()
	ty := schema.Block.ImpliedType()

	state := cty.ObjectVal(map[string]cty.Value{
		"input":  cty.StringVal("input"),
		"output": cty.StringVal("input"),
		"trigger": cty.ListVal([]cty.Value{
			cty.StringVal("a"), cty.StringVal("b"),
		}),
		"uuid": cty.StringVal("not-quite-a-uuid"),
	})

	jsState, err := ctyjson.Marshal(state, ty)
	if err != nil {
		t.Fatal(err)
	}

	// empty
	req := providers.UpgradeResourceStateRequest{
		TypeName:     "terraform_box",
		RawStateJSON: jsState,
	}

	resp := upgradeDataResourceState(req)
	if resp.Diagnostics.HasErrors() {
		t.Error("upgrade state error:", resp.Diagnostics.ErrWithWarnings())
	}

	if !resp.UpgradedState.RawEquals(state) {
		t.Errorf("prior state was:\n%#v\nupgraded state is:\n%#v\n", state, resp.UpgradedState)
	}
}

func TestManagedDataRead(t *testing.T) {
	req := providers.ReadResourceRequest{
		TypeName: "terraform_box",
		PriorState: cty.ObjectVal(map[string]cty.Value{
			"input":  cty.StringVal("input"),
			"output": cty.StringVal("input"),
			"trigger": cty.ListVal([]cty.Value{
				cty.StringVal("a"), cty.StringVal("b"),
			}),
			"uuid": cty.StringVal("not-quite-a-uuid"),
		}),
	}

	resp := readDataResourceState(req)
	if resp.Diagnostics.HasErrors() {
		t.Fatal("unexpected error", resp.Diagnostics.ErrWithWarnings())
	}

	if !resp.NewState.RawEquals(req.PriorState) {
		t.Errorf("prior state was:\n%#v\nnew state is:\n%#v\n", req.PriorState, resp.NewState)
	}
}

func TestManagedDataPlan(t *testing.T) {
	schema := dataResourceSchema().Block
	ty := schema.ImpliedType()

	for name, tc := range map[string]struct {
		prior    cty.Value
		proposed cty.Value
		planned  cty.Value
	}{
		"create": {
			prior: cty.NullVal(ty),
			proposed: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.NullVal(cty.DynamicPseudoType),
				"output":  cty.NullVal(cty.DynamicPseudoType),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.NullVal(cty.String),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.NullVal(cty.DynamicPseudoType),
				"output":  cty.NullVal(cty.DynamicPseudoType),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.UnknownVal(cty.String),
			}),
		},

		"create-output": {
			prior: cty.NullVal(ty),
			proposed: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.NullVal(cty.DynamicPseudoType),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.NullVal(cty.String),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.UnknownVal(cty.String),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.UnknownVal(cty.String),
			}),
		},

		"update-input": {
			prior: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			proposed: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.UnknownVal(cty.List(cty.String)),
				"output":  cty.StringVal("input"),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.UnknownVal(cty.List(cty.String)),
				"output":  cty.UnknownVal(cty.List(cty.String)),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
		},

		"update-trigger": {
			prior: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			proposed: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.StringVal("new-value"),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.UnknownVal(cty.String),
				"trigger": cty.StringVal("new-value"),
				"uuid":    cty.UnknownVal(cty.String),
			}),
		},

		"update-input-trigger": {
			prior: cty.ObjectVal(map[string]cty.Value{
				"input":  cty.StringVal("input"),
				"output": cty.StringVal("input"),
				"trigger": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("value"),
				}),
				"uuid": cty.StringVal("not-quite-a-uuid"),
			}),
			proposed: cty.ObjectVal(map[string]cty.Value{
				"input":  cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"output": cty.StringVal("input"),
				"trigger": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("new value"),
				}),
				"uuid": cty.StringVal("not-quite-a-uuid"),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":  cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"output": cty.UnknownVal(cty.List(cty.String)),
				"trigger": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("new value"),
				}),
				"uuid": cty.UnknownVal(cty.String),
			}),
		},
	} {
		t.Run("plan-"+name, func(t *testing.T) {
			req := providers.PlanResourceChangeRequest{
				TypeName:         "terraform_box",
				PriorState:       tc.prior,
				ProposedNewState: tc.proposed,
			}

			resp := planDataResourceChange(req)
			if resp.Diagnostics.HasErrors() {
				t.Fatal(resp.Diagnostics.ErrWithWarnings())
			}

			if !resp.PlannedState.RawEquals(tc.planned) {
				t.Errorf("expected:\n%#v\ngot:\n%#v\n", tc.planned, resp.PlannedState)
			}
		})
	}
}

func TestManagedDataApply(t *testing.T) {
	testUUIDHook = func() string {
		return "not-quite-a-uuid"
	}
	defer func() {
		testUUIDHook = nil
	}()

	schema := dataResourceSchema().Block
	ty := schema.ImpliedType()

	for name, tc := range map[string]struct {
		prior   cty.Value
		planned cty.Value
		state   cty.Value
	}{
		"create": {
			prior: cty.NullVal(ty),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.NullVal(cty.DynamicPseudoType),
				"output":  cty.NullVal(cty.DynamicPseudoType),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.UnknownVal(cty.String),
			}),
			state: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.NullVal(cty.DynamicPseudoType),
				"output":  cty.NullVal(cty.DynamicPseudoType),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
		},

		"create-output": {
			prior: cty.NullVal(ty),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.UnknownVal(cty.String),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.UnknownVal(cty.String),
			}),
			state: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
		},

		"update-input": {
			prior: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"output":  cty.UnknownVal(cty.List(cty.String)),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			state: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"output":  cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
		},

		"update-trigger": {
			prior: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.NullVal(cty.DynamicPseudoType),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.UnknownVal(cty.String),
				"trigger": cty.StringVal("new-value"),
				"uuid":    cty.UnknownVal(cty.String),
			}),
			state: cty.ObjectVal(map[string]cty.Value{
				"input":   cty.StringVal("input"),
				"output":  cty.StringVal("input"),
				"trigger": cty.StringVal("new-value"),
				"uuid":    cty.StringVal("not-quite-a-uuid"),
			}),
		},

		"update-input-trigger": {
			prior: cty.ObjectVal(map[string]cty.Value{
				"input":  cty.StringVal("input"),
				"output": cty.StringVal("input"),
				"trigger": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("value"),
				}),
				"uuid": cty.StringVal("not-quite-a-uuid"),
			}),
			planned: cty.ObjectVal(map[string]cty.Value{
				"input":  cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"output": cty.UnknownVal(cty.List(cty.String)),
				"trigger": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("new value"),
				}),
				"uuid": cty.UnknownVal(cty.String),
			}),
			state: cty.ObjectVal(map[string]cty.Value{
				"input":  cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"output": cty.ListVal([]cty.Value{cty.StringVal("new-input")}),
				"trigger": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("new value"),
				}),
				"uuid": cty.StringVal("not-quite-a-uuid"),
			}),
		},
	} {
		t.Run("apply-"+name, func(t *testing.T) {
			req := providers.ApplyResourceChangeRequest{
				TypeName:     "terraform_box",
				PriorState:   tc.prior,
				PlannedState: tc.planned,
			}

			resp := applyDataResourceChange(req)
			if resp.Diagnostics.HasErrors() {
				t.Fatal(resp.Diagnostics.ErrWithWarnings())
			}

			if !resp.NewState.RawEquals(tc.state) {
				t.Errorf("expected:\n%#v\ngot:\n%#v\n", tc.state, resp.NewState)
			}
		})
	}
}
