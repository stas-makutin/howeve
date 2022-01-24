package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stas-makutin/howeve/utils"
	"gopkg.in/yaml.v3"
)

type HttpAssetFlag byte

const (
	HAFShowHidden = HttpAssetFlag(1 << iota)
	HAFDirListing
)

var httpAssetFlagMap = map[string]HttpAssetFlag{
	"show-hidden": HAFShowHidden,
	"dir-listing": HAFDirListing,
}

func (flags *HttpAssetFlag) UnmarshalYAML(node *yaml.Node) (err error) {
	*flags = 0
	utils.ParseOptions(node.Value, func(flag string) bool {
		if fl, ok := httpAssetFlagMap[flag]; ok {
			*flags |= fl
			return true
		}
		err = fmt.Errorf("line %d, column %d: unknown flag '%v'", node.Line, node.Column, flag)
		return false
	})
	return err
}

func (flags HttpAssetFlag) String() string {
	var res strings.Builder
	for s, mask := range httpAssetFlagMap {
		if (flags & mask) != 0 {
			if res.Len() > 0 {
				res.WriteString(",")
			}
			res.WriteString(s)
		}
	}
	return res.String()
}

func (flags HttpAssetFlag) MarshalYAML() (interface{}, error) {
	return flags.String(), nil
}

func (flags HttpAssetFlag) MarshalJSON() ([]byte, error) {
	return json.Marshal(flags.String())
}

type HTTPAsset struct {
	Route     string        `yaml:"route,omitempty" json:"route,omitempty"`
	Path      string        `yaml:"path,omitempty" json:"path,omitempty"`
	IndexFile string        `yaml:"indexFile,omitempty" json:"indexFile,omitempty"`
	Includes  []string      `yaml:"includes,omitempty" json:"includes,omitempty"`
	Excludes  []string      `yaml:"excludes,omitempty" json:"excludes,omitempty"`
	Flags     HttpAssetFlag `yaml:"flags,omitempty" json:"flags,omitempty"`
}
