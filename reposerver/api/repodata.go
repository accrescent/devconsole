package api

import (
	"strings"

	"golang.org/x/exp/slices"
)

var abis = []string{
	"arm64_v8a",
	"armeabi_v7a",
	"x86_64",
	"x86",
}

func getSplitInfo(s string) (name string, t splitType, typeName string) {
	typeName = strings.TrimSuffix(strings.TrimPrefix(s, "base-"), ".apk")
	switch {
	case typeName == "master":
		t = master
	case strings.HasSuffix(typeName, "dpi"):
		t = density
	case slices.Contains(abis, typeName):
		t = abi
		typeName = strings.Replace(typeName, "_", "-", 1)
	default:
		t = lang
	}

	name = "split." + typeName + ".apk"
	name = strings.Replace(name, "split.master.apk", "base.apk", 1)

	return
}

type splitType int

const (
	abi splitType = iota
	density
	lang

	master
)

type repoData struct {
	Version       string   `json:"verison"`
	VersionCode   int      `json:"version_code"`
	ABISplits     []string `json:"abi_splits"`
	DensitySplits []string `json:"density_splits"`
	LangSplits    []string `json:"lang_splits"`
}
