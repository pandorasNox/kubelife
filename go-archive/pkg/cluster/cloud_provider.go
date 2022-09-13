package cluster

type cloudProvider interface {
	provision() error
}

type HetznerCloudProvider struct {
}

func (hcp HetznerCloudProvider) provision() error {
	return nil
}
