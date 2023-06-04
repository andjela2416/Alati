package store

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	//posts       = "posts/%s/%s"
	//postsLabels = "posts/%s/%s/%s"
	//all         = "posts"
	config        = "configs/%s/%s"
	groups        = "groups/%s/%s/%s/%s" // groups/idg/version/labels/idc
	configs       = "configs/%s/%s"
	configsLabels = "configs/%s/%s/%s"
	all           = "configs"
	allGroups     = "groups"
	postId        = "key/%s"
)

func generateKey(version string, labels string) (string, string) {
	id := uuid.New().String()
	if labels != "" {
		return fmt.Sprintf(configsLabels, id, version, labels), id
	} else {
		return fmt.Sprintf(configs, id, version), id
	}

}

func constructKey(id string, version string, labels string) string {
	if labels != "" {
		return fmt.Sprintf(configsLabels, id, version, labels)
	} else {
		return fmt.Sprintf(configs, id, version)
	}

}

func constructKey2(id string) string {

	return fmt.Sprintf(postId, id)
}
