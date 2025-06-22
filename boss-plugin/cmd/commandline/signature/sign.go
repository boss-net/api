package signature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boss-net/api/boss-plugin/internal/utils/encryption"
	"github.com/boss-net/api/boss-plugin/internal/utils/log"
	"github.com/boss-net/api/boss-plugin/pkg/plugin_packager/decoder"
	"github.com/boss-net/api/boss-plugin/pkg/plugin_packager/signer/withkey"
)

func Sign(bosspkgPath string, privateKeyPath string, verification *decoder.Verification) error {
	// read the plugin and private key
	plugin, err := os.ReadFile(bosspkgPath)
	if err != nil {
		log.Error("Failed to read plugin file: %v", err)
		return err
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Error("Failed to read private key file: %v", err)
		return err
	}

	privateKey, err := encryption.LoadPrivateKey(privateKeyBytes)
	if err != nil {
		log.Error("Failed to load private key: %v", err)
		return err
	}

	// sign the plugin
	pluginFile, err := withkey.SignPluginWithPrivateKey(plugin, verification, privateKey)
	if err != nil {
		log.Error("Failed to sign plugin: %v", err)
		return err
	}

	// write the signed plugin to a file
	dir := filepath.Dir(bosspkgPath)
	base := filepath.Base(bosspkgPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	outputPath := filepath.Join(dir, fmt.Sprintf("%s.signed%s", name, ext))

	err = os.WriteFile(outputPath, pluginFile, 0644)
	if err != nil {
		log.Error("Failed to write signed plugin file: %v", err)
		return err
	}

	log.Info("Plugin signed successfully, output path: %s", outputPath)

	return nil
}
