package encryptor

import (
	"github.com/spf13/viper"

	"github.com/KurosawaAngel/gobackup/config"
	"github.com/KurosawaAngel/gobackup/logger"
)

// Base encryptor
type Base struct {
	model       config.ModelConfig
	viper       *viper.Viper
	archivePath string
}

// Encryptor interface
type Encryptor interface {
	perform() (encryptPath string, err error)
}

func newBase(archivePath string, model config.ModelConfig) (base *Base) {
	base = &Base{
		archivePath: archivePath,
		model:       model,
		viper:       model.EncryptWith.Viper,
	}
	return
}

// Run encryptor
func Run(archivePath string, model config.ModelConfig) (encryptPath string, err error) {
	l := logger.Tag("Encryptor")

	base := newBase(archivePath, model)
	var enc Encryptor
	switch model.EncryptWith.Type {
	case "openssl":
		enc = NewOpenSSL(base)
	default:
		encryptPath = archivePath
		return
	}

	l.Info("encrypt | " + model.EncryptWith.Type)
	encryptPath, err = enc.perform()
	if err != nil {
		return
	}
	l.Info("encrypted:", encryptPath)

	// save Extension
	model.Viper.Set("Ext", model.Viper.GetString("Ext")+".enc")

	return
}
