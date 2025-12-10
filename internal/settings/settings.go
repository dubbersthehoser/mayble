package settings

type Settings struct {
	ConfigDir  string `json:"config_dir"` 
	ConfigFile string `json:"config_file"`

	DBDriver   string `json:"db_driver"`
	DBFile     string `json:"db_file"`
}

