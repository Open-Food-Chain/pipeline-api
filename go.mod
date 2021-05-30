module github.com/The-New-Fork/api-pipeline

go 1.15

require (
	github.com/go-chi/render v1.0.1
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/stretchr/testify v1.6.1
	github.com/unchain/pipeline v0.0.0-20201221180813-38ecf41bce98
	github.com/unchainio/interfaces v0.2.1
	github.com/unchainio/pkg v0.22.1
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/spf13/viper v1.2.2 => github.com/unchain/viper v1.2.2-0.20190712174521-9bf201c29832
	github.com/unchain/pipeline v0.0.0-20201221180813-38ecf41bce98 => github.com/The-New-Fork/pipeline v0.0.0-20201221180813-38ecf41bce98
)
