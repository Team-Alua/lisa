package client

type ContentTarget struct {
	TitleId string `yaml:"titleId"`
	DirectoryName string `yaml:"directoryName"`
	Blocks uint64 `yaml:"blocks"`
}

type Content struct {
	Type string `yaml:"Type"`
	AccountId uint64 `yaml:"AccountId"`
	Target *ContentTarget `yaml:"Target"`
}

