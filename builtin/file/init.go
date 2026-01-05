package file

import (
	// Import file transformers to register them
	_ "github.com/nodyhub/transformer/builtin/file/read"
	_ "github.com/nodyhub/transformer/builtin/file/write"
)
