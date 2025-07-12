package types

type AppConfig struct {
	Port       string `yaml:"port"`
	CydresUrl  string `yaml:"cydresUrl"`
	StorageUrl string `yaml:"storageUrl"`
}
