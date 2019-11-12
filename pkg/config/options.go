package config

type StructureTestOptions struct {
	Host        string
	Username    string
	Password    string
	Driver      string
	TestReport  string
	ConfigFiles []string

	JSON    bool
	Quiet   bool
	Force   bool
	NoColor bool
}
