package generator

import "github.com/potato/simple-restful-api/infra/command"

type EntityGenerator interface {
	CreateEntityAnsSave(events []*command.Event) error
}
