package builtin

import (
	// Import built-in transformers to register them
	_ "github.com/nodyhub/transformer/builtin/echo"
	_ "github.com/nodyhub/transformer/builtin/file"
	_ "github.com/nodyhub/transformer/builtin/sarif"
	_ "github.com/nodyhub/transformer/builtin/shell"
	_ "github.com/nodyhub/transformer/builtin/tools"
)
