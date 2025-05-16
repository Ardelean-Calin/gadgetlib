package functions

type FunctionECM struct {
	InstanceName string
	DevAddr      string
	HostAddr     string
}

func (f FunctionECM) Apply(configPath, funcPath string) error {
	panic("unimplemented")
}

func (f FunctionECM) Name() string {
	panic("unimplemented")
}
