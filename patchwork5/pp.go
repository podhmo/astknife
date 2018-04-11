package patchwork5

import (
	"encoding/json"
	"os"
)

func dumpRegions(regions []*Region) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(regions)
}
