package store

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	//posts       = "posts/%s/%s"
	//postsLabels = "posts/%s/%s/%s"
	//all         = "posts"
	//config        = "configs/%s/%s"
	//groups        = "groups/%s/%s/%s/%s" // groups/idg/version/labels/idc
	groups        = "groups/%s/%s/"
	groups2       = "groups/%s/"
	configs       = "configs/%s/%s"
	configs2      = "configs/%s"
	configsLabels = "configs/%s/%s/%s"
	//groupsLabels  = "groups/%s/%s/%s"
	all       = "configs"
	allGroups = "groups"
)

func generateKey(version string, labels string) (string, string) {
	id := uuid.New().String()
	if labels != "" {
		return fmt.Sprintf(configsLabels, id, version, labels), id
	} else {
		return fmt.Sprintf(configs, id, version), id
	}

}

func generateGroupKey(version string) (string, string) {
	id := uuid.New().String()
	/*if labels != "" {
		return fmt.Sprintf(groupsLabels, id, version, labels), id
	} else {
	*/return fmt.Sprintf(groups, id, version), id
	//}

}

func constructKey(id string, version string, labels string) string {
	if labels != "" {
		return fmt.Sprintf(configsLabels, id, version, labels)
	} else {
		return fmt.Sprintf(configs, id, version)
	}

}
func constructKey2(id string) string {
	return fmt.Sprintf(configs2, id)
}
func constructGroupKey(id string, version string) string {
	/*if labels != "" {
		return fmt.Sprintf(groupsLabels, id, version, labels)
	} else {
	*/return fmt.Sprintf(groups, id, version)
	//}

}
func constructGroupKey2(id string) string {
	/*if labels != "" {
		return fmt.Sprintf(groupsLabels, id, version, labels)
	} else {
	*/return fmt.Sprintf(groups2, id)
	//}

}
