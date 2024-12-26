package shared

var FhirbaseOptions FhirbaseConfig

type DatabaseOptions struct {
	SSLMode  string `mapstructure:"sslmode"`
	Host     string `mapstructure:"host"`
	Port     uint64 `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Database string `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

type FhirbaseConfig struct {
	Database    DatabaseOptions `mapstructure:"database"`
	Display     DisplayOptions  `mapstructure:"display"`
	Debug       DebugOptions    `mapstructure:"debug"`
	Keys        KeyOptions      `mapstructure:"keys"`
	FhirVersion string          `mapstructure:"fhir"`
}

type KeyOptions struct {
	Quit   string `mapstructure:"quit"`
	Up     string `mapstructure:"up"`
	Down   string `mapstructure:"down"`
	Select string `mapstructure:"select"`
	Toggle string `mapstructure:"toggle"`
	Back   string `mapstructure:"back"`
	Help  string `mapstructure:"help"`
}

type DebugOptions struct {
	FilePath    string `mapstructure:"file_path"`
	LogMessages bool   `mapstructure:"log_messages"`
}

type DisplayOptions struct {
	Cursor string `mapstructure:"cursor"`
}

const (
	RootView      string = "root"
	ConfigView    string = "config"
	InitDbView    string = "init"
	LoadDbView    string = "load"
	BulkGetView   string = "bulk"
	TransformView string = "transform"
	WebServerView string = "web"
	UpdateView    string = "update"
)
