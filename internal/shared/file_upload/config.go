package file_upload

type Config struct {
	MaxSize           int64
	AllowedTypes      []string
	AllowedExtensions []string
	UploadDir         string
}

func ProductFileConfig() *Config {
	return &Config{
		MaxSize:      2 * 1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/jpg", "image/png", "image/webp"},
		UploadDir:    "./uploads/products",
	}
}
