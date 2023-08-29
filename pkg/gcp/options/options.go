package options

type Options struct {
	Project                         string
	ImpersonatedServiceAccountEmail string
}
type Option interface {
	Apply(*Options)
}

type project string

func (o project) Apply(i *Options) {
	if o != "" {
		i.Project = string(o)
	}
}

func WithProject(d string) Option {
	return project(d)
}

type impersonatedServiceAccountEmail string

func (o impersonatedServiceAccountEmail) Apply(i *Options) {
	if o != "" {
		i.ImpersonatedServiceAccountEmail = string(o)
	}
}

func WithImpersonatedServiceAccountEmail(d string) Option {
	return impersonatedServiceAccountEmail(d)
}
