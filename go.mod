module github.com/ddkwork/unison

go 1.23

require (
	github.com/OpenPrinting/goipp v1.1.0
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/ebitengine/purego v0.7.1
	github.com/google/uuid v1.6.0
	github.com/grandcat/zeroconf v1.0.0
	github.com/richardwilkes/json v0.3.0
	github.com/yuin/goldmark v1.7.4
	golang.org/x/image v0.18.0
	golang.org/x/sys v0.22.0
	golang.org/x/text v0.16.0
)

require (
	github.com/alecthomas/chroma/v2 v2.14.0
	github.com/ddkwork/golibrary v0.0.0-20240728131742-3497a6ed9010
	github.com/richardwilkes/toolbox v1.99.0
)

replace github.com/richardwilkes/toolbox => ./internal/toolbox-1.99.0

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dc0d/caseconv v0.5.0 // indirect
	github.com/dlclark/regexp2 v1.11.2 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.61 // indirect
	github.com/pkg/term v1.1.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/gofumpt v0.6.0 // indirect
)
