module github.com/Xiangzhuzhihui/apisix-go-plugin-runner

go 1.18

require (
	github.com/ReneKroon/ttlcache/v2 v2.4.0
	github.com/api7/ext-plugin-proto v0.6.0
	github.com/google/flatbuffers v2.0.0+incompatible
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/thediveo/enumflag v0.10.1
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/miekg/dns v1.0.14 => github.com/miekg/dns v1.1.25
	// github.com/thediveo/enumflag@v0.10.1 depends on github.com/spf13/cobra@v0.0.7
	github.com/spf13/cobra v0.0.7 => github.com/spf13/cobra v1.2.1
)
