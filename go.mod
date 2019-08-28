module smartqn/common

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/aiwuTech/fileLogger v0.0.0-20150625084638-fdc0c5b08dd6
	github.com/garyburd/redigo v1.6.0
	github.com/golang/protobuf v1.3.2
	github.com/linnv/logx v0.0.0-20190825041807-16e58b3a5351
	github.com/satori/go.uuid v1.2.0
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/garyburd/redigo v1.6.0 => ./vendor/github.com/garyburd/redigo
