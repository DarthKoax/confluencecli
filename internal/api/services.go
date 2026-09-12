package api

import (
	"github.com/darthkoax/confluencecli/internal/client"
)

type Services struct {
	Content       *ContentService
	Spaces        *SpaceService
	Search        *SearchService
	Users         *UserService
	Groups        *GroupService
	Settings      *SettingsService
	Audit         *AuditService
	Templates     *TemplateService
	ContentStates *ContentStateService
	InlineTasks   *InlineTaskService
	Relations     *RelationService
	LongTasks     *LongTaskService
	System        *SystemService
	Blueprints    *BlueprintService
	HealthCheck   *HealthCheckService
}

func NewServices(c *client.Client) *Services {
	return &Services{
		Content:       NewContentService(c),
		Spaces:        NewSpaceService(c),
		Search:        NewSearchService(c),
		Users:         NewUserService(c),
		Groups:        NewGroupService(c),
		Settings:      NewSettingsService(c),
		Audit:         NewAuditService(c),
		Templates:     NewTemplateService(c),
		ContentStates: NewContentStateService(c),
		InlineTasks:   NewInlineTaskService(c),
		Relations:     NewRelationService(c),
		LongTasks:     NewLongTaskService(c),
		System:        NewSystemService(c),
		Blueprints:    NewBlueprintService(c),
		HealthCheck:   NewHealthCheckService(c),
	}
}
