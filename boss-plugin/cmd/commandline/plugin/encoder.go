package plugin

import (
	"bytes"

	"github.com/boss-net/api/boss-plugin/internal/utils/log"
	"gopkg.in/yaml.v3"
)

func marshalYamlBytes(v any) []byte {
	buf := bytes.NewBuffer([]byte{})
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err := encoder.Encode(v)
	if err != nil {
		log.Error("failed to marshal yaml: %s", err)
		return nil
	}
	return buf.Bytes()
}
