package export

var defaultConfig = Config{
	Host:       "localhost",
	Port:       "9000",
	Database:   "default",
	Username:   "",
	Password:   "",
	Format:     "CSVWithNames",
	Query:      "",
	OutputFile: "",
	BatchSize:  0,
}

type Config struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Format     string `yaml:"format"`
	Query      string `yaml:"query"`
	OutputFile string `yaml:"outputFile"`
	BatchSize  int    `yaml:"batchSize"`
}

func GetConfig() Config {
	conf := defaultConfig
	return conf
}
