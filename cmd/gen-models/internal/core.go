package internal

import (
	"gorm.io/gen"
)

func CoreModels(g *gen.Generator) []interface{} {
	entity := g.GenerateModel("entities")
	activities := g.GenerateModel("activities")
	activityDeadline := g.GenerateModel("activity_deadlines")
	activityComment := g.GenerateModel("activity_comments")
	activityMention := g.GenerateModel("activity_mentions")
	currencyExchange := g.GenerateModel("currency_exchanges")
	workSpace := g.GenerateModel("workspaces")
	workSpaceModule := g.GenerateModel("workspace_modules")
	module := g.GenerateModel("modules")
	moduleSections := g.GenerateModel("module_sections")
	stage := g.GenerateModel("stages")
	connection := g.GenerateModel("connections")
	partyReferences := g.GenerateModel("party_references")

	notification := g.GenerateModel("notifications")
	mention := g.GenerateModel("mentions")

	address := g.GenerateModel("addresses")
	return []interface{}{
		entity,
		activities,
		stage,
		currencyExchange,
		module,
		moduleSections,
		activityDeadline,
		activityComment,
		activityMention,
		notification,
		mention,

		workSpace,
		workSpaceModule,

		address,
		connection,
		partyReferences,
	}
}
