package internal 

import (
	"gorm.io/gen"
)

func ChatModels(g *gen.Generator) []interface{} {
	chat := g.GenerateModel("chats")
	chatMember := g.GenerateModel("chat_members")
	chatMessage := g.GenerateModel("chat_messages")

	return []interface{}{
		chat,
		chatMember,
		chatMessage,
	}
}
