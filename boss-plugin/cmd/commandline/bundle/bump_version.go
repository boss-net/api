package bundle

import (
	"github.com/boss-net/api/boss-plugin/internal/utils/log"
	"github.com/boss-net/api/boss-plugin/pkg/entities/manifest_entities"
)

func BumpVersion(bundlePath string, targetVersion string) {
	packager, err := loadBundlePackager(bundlePath)
	if err != nil {
		log.Error("Failed to load bundle packager: %v", err)
		return
	}

	targetVersionObject, err := manifest_entities.NewVersion(targetVersion)
	if err != nil {
		log.Error("Failed to parse target version: %v", err)
		return
	}

	packager.BumpVersion(targetVersionObject)
	if err := packager.Save(); err != nil {
		log.Error("Failed to save bundle packager: %v", err)
		return
	}
}
