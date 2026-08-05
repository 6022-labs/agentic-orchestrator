package enums

import "strings"

// BridgePlatform identifies an external chat platform bridged to the broker.
type BridgePlatform string
type BridgePlatformList []BridgePlatform

const (
	BridgePlatformSlack BridgePlatform = "slack"
)

func GetAllBridgePlatforms() BridgePlatformList {
	return BridgePlatformList{
		BridgePlatformSlack,
	}
}

func (p BridgePlatform) String() string {
	return string(p)
}

func (p BridgePlatformList) String() string {
	var keys []string
	for _, v := range p {
		keys = append(keys, string(v))
	}
	return strings.Join(keys, ", ")
}

func (p BridgePlatformList) Contains(platform BridgePlatform) bool {
	for _, v := range p {
		if v == platform {
			return true
		}
	}
	return false
}
