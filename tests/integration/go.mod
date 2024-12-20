module go.nhat.io/requestid/tests/integration

go 1.22
toolchain go1.22.5

require (
	github.com/go-chi/chi/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
	go.nhat.io/requestid v0.1.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0
)

require (
	github.com/bool64/ctxd v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.nhat.io/requestid => ../../
