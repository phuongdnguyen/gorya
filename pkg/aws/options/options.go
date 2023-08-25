package options

type Options struct {
	AwsRegion   string
	AwsEndpoint string
	AwsRoleArn  string
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

type awsRoleArn string

func (o awsRoleArn) Apply(i *Options) {
	if o != "" {
		i.AwsRoleArn = string(o)
	}
}

func WithRoleArn(d string) Option {
	return awsRoleArn(d)
}
