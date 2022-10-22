package earlyconfig

import (
	"github.com/hashicorp/terraform/internal/tfdiags"
	"github.com/titoluyo/terraform-config-inspect/tfconfig"
)

// LoadModule loads some top-level metadata for the module in the given
// directory.
func LoadModule(dir string) (*tfconfig.Module, tfdiags.Diagnostics) {
	mod, diags := tfconfig.LoadModule(dir)
	return mod, wrapDiagnostics(diags)
}
