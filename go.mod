module build_kubernetes_v0.1.0

go 1.20

require (
	bitbucket.org/kardianos/osext v0.0.0-20181027061946-15c52d0993e9
	github.com/GoogleCloudPlatform/kubernetes v0.0.0-00010101000000-000000000000
	github.com/coreos/go-etcd v2.0.0+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/fsouza/go-dockerclient v0.0.0-20140601215550-3b6f84ca70be
	github.com/gorilla/mux v1.8.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
)

require (
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.1.0 // indirect
)

replace github.com/GoogleCloudPlatform/kubernetes => ./
