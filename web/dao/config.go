package dao

var config *Options

type Options struct {
	Type int
}

func NewOptions(type_ int) {
	config = &Options{Type: type_}
}

func GetConfig() *Options {
	return config
}
