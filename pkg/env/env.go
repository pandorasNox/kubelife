package env

type Cfg struct {
	HcloudToken string `envconfig:"KUBELIFE_HCLOUD_TOKEN" desc:"Foo Bar, what do you want more?"`
}
