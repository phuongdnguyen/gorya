package options

type GCEOptions struct {
	ProjectID  string
	Zone       string
	InstanceID string
}

func NewGCEOptions(projectID, zone, instanceID string) *GCEOptions {
	return &GCEOptions{
		ProjectID:  projectID,
		Zone:       zone,
		InstanceID: instanceID,
	}
}

