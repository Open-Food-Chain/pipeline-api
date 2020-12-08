module github.com/The-New-Fork/api-pipeline

go 1.15

replace github.com/spf13/viper v1.2.2 => github.com/unchainio/viper v1.2.2-0.20190712174521-9bf201c29832

require (
	github.com/go-chi/render v1.0.1
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/stretchr/testify v1.6.1
	github.com/unchain/pipeline v0.0.0-20201207235952-afa156021d9f
	github.com/unchainio/interfaces v0.2.1
	github.com/unchainio/pkg v0.22.1
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135
)
