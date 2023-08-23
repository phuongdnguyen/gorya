package options

type Options struct {
	AwsRegion   string
	AwsEndpoint string
}

type Option interface {
	Apply(*Options)
}

type awsRegion string

func (o awsRegion) Apply(i *Options) {
	if o != "" {
		i.AwsRegion = string(o)
	}
}

func WithRegion(d string) Option {
	return awsRegion(d)
}

type awsEndpoint string

func (o awsEndpoint) Apply(i *Options) {
	if o != "" {
		i.AwsEndpoint = string(o)
	}
}

func WithEndpoint(d string) Option {
	return awsEndpoint(d)
}
