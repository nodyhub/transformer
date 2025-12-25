package sarif

import (
	// Re-export common types for convenience
	"github.com/nodyhub/transformer/builtin/sarif/common"

	// Import sub-packages to register them
	_ "github.com/nodyhub/transformer/builtin/sarif/convert"
	_ "github.com/nodyhub/transformer/builtin/sarif/merge"
	_ "github.com/nodyhub/transformer/builtin/sarif/suppress"
)

type (
	Report           = common.Report
	Run              = common.Run
	Tool             = common.Tool
	Driver           = common.Driver
	Result           = common.Result
	Message          = common.Message
	Location         = common.Location
	PhysicalLocation = common.PhysicalLocation
	ArtifactLocation = common.ArtifactLocation
	Region           = common.Region
)
