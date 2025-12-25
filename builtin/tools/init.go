package tools

import (
	_ "github.com/nodyhub/transformer/builtin/tools/codeql"
	_ "github.com/nodyhub/transformer/builtin/tools/gosec"
	_ "github.com/nodyhub/transformer/builtin/tools/nmap"
	_ "github.com/nodyhub/transformer/builtin/tools/nuclei"
	_ "github.com/nodyhub/transformer/builtin/tools/osv-scanner"
	_ "github.com/nodyhub/transformer/builtin/tools/semgrep"
	_ "github.com/nodyhub/transformer/builtin/tools/trivy"
	_ "github.com/nodyhub/transformer/builtin/tools/trufflehog"
)
