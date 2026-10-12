package auth

import common "github.com/quangbach27/golang-common"

type Resource struct {
	common.Enum[ResourceType]
}

type ResourceType string

func (rt ResourceType) Values() []string {
	return []string{"hour", "training"}
}

var (
	ResourceHour     = common.MustEnum[Resource]("hour")
	ResourceTraining = common.MustEnum[Resource]("training")
)
