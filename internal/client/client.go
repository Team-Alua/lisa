package client

type ContentTarget struct {
	TitleId string `yaml:"titleId"`
	DirectoryName string `yaml:"directoryName"`
	Blocks uint16 `yaml:"blocks"`
}

type Content struct {
	Type string `yaml:"Type"`
	AccountId string `yaml:"AccountId"`
	Target *ContentTarget `yaml:"Target"`
}

