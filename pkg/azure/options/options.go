package options

type Options struct {
	//target SubscriptionId
	SubscriptionId string
	//target ResourceGroupName
	ResourceGroupName string
	//target TenantId
	TenantId string
}

type Option interface {
	Apply(*Options)
}

type tenantId string

func (o tenantId) Apply(i *Options) {
	if o != "" {
		i.TenantId = string(o)
	}
}

func WithTenantId(d string) Option {
	return tenantId(d)
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

type resourcegroupName string

func (o resourcegroupName) Apply(i *Options) {
	if o != "" {
		i.ResourceGroupName = string(o)
	}
}

func WithResourceGroupName(d string) Option {
	return resourcegroupName(d)
}
