package datafactory

import (
	defsecTypes "github.com/aquasecurity/defsec/pkg/types"
)

type DataFactory struct {
	DataFactories []Factory
}

type Factory struct {
	defsecTypes.Metadata
	EnablePublicNetwork defsecTypes.BoolValue
}
