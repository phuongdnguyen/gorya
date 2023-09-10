package options

type Options struct {
	//target SubscriptionId
	SubscriptionId string
	//target TargetResourceGroups
	TargetResourceGroups []string
}

type Option interface {
	Apply(*Options)
}

type subscriptionId string

func (o subscriptionId) Apply(i *Options) {
	if o != "" {
		i.SubscriptionId = string(o)
	}
}

func WithSubscriptionId(d string) Option {
	return subscriptionId(d)
}

type targetResourceGroups []string

func (o targetResourceGroups) Apply(i *Options) {
	if o != nil {
		i.TargetResourceGroups = []string(o)
	}
}

func WithTargetResourceGroups(d []string) Option {
	return targetResourceGroups(d)
}
