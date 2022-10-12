package renderer

import (
	"github.com/hashicorp/terraform/internal/command/jsonplan"
	"github.com/hashicorp/terraform/internal/command/jsonprovider"
)

type Plan struct {
	OutputChanges   map[string]jsonplan.Change       `json:"output_changes"`
	ResourceChanges []jsonplan.ResourceChange        `json:"resource_changes"`
	ResourceDrift   []jsonplan.ResourceChange        `json:"resource_drift"`
	ProviderSchema  map[string]jsonprovider.Provider `json:"provider_schema"`
}

func FromJson(plan jsonplan.Plan, providers map[string]jsonprovider.Provider) Plan {
	return Plan{
		OutputChanges:   plan.OutputChanges,
		ResourceChanges: plan.ResourceChanges,
		ResourceDrift:   plan.ResourceDrift,
		ProviderSchema:  providers,
	}
}
