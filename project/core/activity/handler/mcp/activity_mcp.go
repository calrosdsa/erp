package activity_mcp

import (
	"context"
	"encoding/json"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	activity_ucase "erp/project/core/activity/usecase"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ActivityMcp struct {
	mcpServer       *server.MCPServer
	sessionHelper   helpers.SessionHelper
	locale          helpers.Locale
	errorHelper     helpers.ErrorHelper
	permission      repository.PermissionService
	activityUseCase activity_ucase.ActivityUseCase
}

func NewActivityMcp(
	mcpServer *server.MCPServer,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	activityUseCase activity_ucase.ActivityUseCase,
) {
	activityMcp := &ActivityMcp{
		mcpServer:       mcpServer,
		sessionHelper:   helpers.Session,
		locale:          helpers.Locale,  
		errorHelper:     helpers.Error,
		permission:      permission,
		activityUseCase: activityUseCase,
	}

	activityMcp.registerTools()
}

func (a *ActivityMcp) registerTools() {
	// Create Activity Tool
	createActivityToolspec := mcp.NewTool("create_activity",
		mcp.WithDescription("Create a new activity record with the provided activity data"),
		mcp.WithNumber("party_id",
			mcp.Required(),
			mcp.Description("Party ID associated with the activity"),
		),
		mcp.WithString("party_name",
			mcp.Description("Name of the party"),
		),
		mcp.WithNumber("entity_id",
			mcp.Description("Entity ID associated with the activity"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of activity (e.g., 'comment', 'deadline', 'meeting')"),
		),
		mcp.WithBoolean("is_pinned",
			mcp.Description("Whether the activity is pinned (default: false)"),
		),
		// Activity deadline fields
		mcp.WithString("deadline",
			mcp.Description("Activity deadline in ISO format (YYYY-MM-DDTHH:MM:SSZ)"),
		),
		mcp.WithString("deadline_link",
			mcp.Description("Link associated with the deadline"),
		),
		mcp.WithNumber("deadline_party_id",
			mcp.Description("Party ID for the deadline"),
		),
		mcp.WithString("deadline_address",
			mcp.Description("Address for the deadline activity"),
		),
		mcp.WithString("deadline_title",
			mcp.Description("Title for the deadline activity"),
		),
		mcp.WithString("deadline_content",
			mcp.Description("Content for the deadline activity"),
		),
		mcp.WithString("deadline_color",
			mcp.Description("Color for the deadline activity"),
		),
		mcp.WithBoolean("deadline_is_completed",
			mcp.Description("Whether the deadline is completed (default: false)"),
		),
		mcp.WithNumber("deadline_profile_id",
			mcp.Description("Profile ID for the deadline"),
		),
		// Activity comment fields
		mcp.WithString("comment",
			mcp.Description("Comment text for the activity"),
		),
	)

	a.mcpServer.AddTool(createActivityToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return a.handleCreateActivity(ctx, request)
	})

	// Edit Activity Tool
	editActivityToolspec := mcp.NewTool("edit_activity",
		mcp.WithDescription("Update an existing activity record with new data"),
		mcp.WithNumber("party_id",
			mcp.Required(),
			mcp.Description("Party ID associated with the activity"),
		),
		mcp.WithString("party_name",
			mcp.Description("Name of the party"),
		),
		mcp.WithNumber("entity_id",
			mcp.Description("Entity ID associated with the activity"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of activity (e.g., 'comment', 'deadline', 'meeting')"),
		),
		mcp.WithBoolean("is_pinned",
			mcp.Description("Whether the activity is pinned"),
		),
		// Activity deadline fields
		mcp.WithString("deadline",
			mcp.Description("Activity deadline in ISO format (YYYY-MM-DDTHH:MM:SSZ)"),
		),
		mcp.WithString("deadline_link",
			mcp.Description("Link associated with the deadline"),
		),
		mcp.WithNumber("deadline_party_id",
			mcp.Description("Party ID for the deadline"),
		),
		mcp.WithString("deadline_address",
			mcp.Description("Address for the deadline activity"),
		),
		mcp.WithString("deadline_title",
			mcp.Description("Title for the deadline activity"),
		),
		mcp.WithString("deadline_content",
			mcp.Description("Content for the deadline activity"),
		),
		mcp.WithString("deadline_color",
			mcp.Description("Color for the deadline activity"),
		),
		mcp.WithBoolean("deadline_is_completed",
			mcp.Description("Whether the deadline is completed"),
		),
		mcp.WithNumber("deadline_profile_id",
			mcp.Description("Profile ID for the deadline"),
		),
		// Activity comment fields
		mcp.WithString("comment",
			mcp.Description("Comment text for the activity"),
		),
	)

	a.mcpServer.AddTool(editActivityToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return a.handleEditActivity(ctx, request)
	})

	// Delete Activity Tool
	deleteActivityToolspec := mcp.NewTool("delete_activity",
		mcp.WithDescription("Delete an activity record by ID"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("Activity ID to delete"),
		),
	)

	a.mcpServer.AddTool(deleteActivityToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return a.handleDeleteActivity(ctx, request)
	})
}

func (a *ActivityMcp) handleCreateActivity(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := a.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Extract required fields
	partyID, err := request.RequireFloat("party_id")
	if err != nil {
		return mcp.NewToolResultError("Party ID is required: " + err.Error()), nil
	}

	activityType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError("Activity type is required: " + err.Error()), nil
	}

	// Build activity data
	activityData := dto.ActivityData{
		PartyID: int64(partyID),
		Type:    activityType,
	}

	// Handle optional fields
	if partyName := request.GetString("party_name", ""); partyName != "" {
		activityData.PartyName = partyName
	}

	if entityID := request.GetFloat("entity_id", 0); entityID > 0 {
		activityData.EntityID = int(entityID)
	}

	if isPinned := request.GetBool("is_pinned", false); isPinned {
		activityData.IsPinned = &isPinned
	}

	// Handle deadline data if provided
	if deadlineStr := request.GetString("deadline", ""); deadlineStr != "" {
		deadline, parseErr := time.Parse(time.RFC3339, deadlineStr)
		if parseErr != nil {
			return mcp.NewToolResultError("Invalid deadline format (expected ISO format): " + parseErr.Error()), nil
		}

		activityData.ActivityDeadLine = dto.ActivityDeadlineData{
			Fields: dto.ActivityDeadlineFields{
				Deadline:    deadline,
				IsCompleted: request.GetBool("deadline_is_completed", false),
				Color:       request.GetString("deadline_color", ""),
			},
		}

		// Handle optional deadline fields
		if link := request.GetString("deadline_link", ""); link != "" {
			activityData.ActivityDeadLine.Fields.Link = &link
		}

		if deadlinePartyID := request.GetFloat("deadline_party_id", 0); deadlinePartyID > 0 {
			id := int64(deadlinePartyID)
			activityData.ActivityDeadLine.Fields.PartyID = &id
		}

		if address := request.GetString("deadline_address", ""); address != "" {
			activityData.ActivityDeadLine.Fields.Address = &address
		}

		if title := request.GetString("deadline_title", ""); title != "" {
			activityData.ActivityDeadLine.Fields.Title = &title
		}

		if content := request.GetString("deadline_content", ""); content != "" {
			activityData.ActivityDeadLine.Fields.Content = &content
		}

		if profileID := request.GetFloat("deadline_profile_id", 0); profileID > 0 {
			id := int64(profileID)
			activityData.ActivityDeadLine.Fields.ProfileID = &id
		}
	}

	// Handle comment data if provided
	if comment := request.GetString("comment", ""); comment != "" {
		activityData.ActivityComment = dto.ActivityCommentData{
			Fields: dto.ActivityCommentFields{
				Comment: comment,
			},
		}
	}

	// Call the usecase
	result, err := a.activityUseCase.CreateActivity(req, activityData)
	if err != nil {
		return mcp.NewToolResultError("Failed to create activity: " + err.Error()), nil
	}

	// Return success message with result
	message := a.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CommentAddedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)

	response := map[string]interface{}{
		"message":  message,
		"activity": result,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (a *ActivityMcp) handleEditActivity(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := a.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Extract required fields
	partyID, err := request.RequireFloat("party_id")
	if err != nil {
		return mcp.NewToolResultError("Party ID is required: " + err.Error()), nil
	}

	activityType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError("Activity type is required: " + err.Error()), nil
	}

	// Build activity data
	activityData := dto.ActivityData{
		PartyID: int64(partyID),
		Type:    activityType,
	}

	// Handle optional fields
	if partyName := request.GetString("party_name", ""); partyName != "" {
		activityData.PartyName = partyName
	}

	if entityID := request.GetFloat("entity_id", 0); entityID > 0 {
		activityData.EntityID = int(entityID)
	}

	if isPinned := request.GetBool("is_pinned", false); isPinned {
		activityData.IsPinned = &isPinned
	}

	// Handle deadline data if provided
	if deadlineStr := request.GetString("deadline", ""); deadlineStr != "" {
		deadline, parseErr := time.Parse(time.RFC3339, deadlineStr)
		if parseErr != nil {
			return mcp.NewToolResultError("Invalid deadline format (expected ISO format): " + parseErr.Error()), nil
		}

		activityData.ActivityDeadLine = dto.ActivityDeadlineData{
			Fields: dto.ActivityDeadlineFields{
				Deadline:    deadline,
				IsCompleted: request.GetBool("deadline_is_completed", false),
				Color:       request.GetString("deadline_color", ""),
			},
		}

		// Handle optional deadline fields
		if link := request.GetString("deadline_link", ""); link != "" {
			activityData.ActivityDeadLine.Fields.Link = &link
		}

		if deadlinePartyID := request.GetFloat("deadline_party_id", 0); deadlinePartyID > 0 {
			id := int64(deadlinePartyID)
			activityData.ActivityDeadLine.Fields.PartyID = &id
		}

		if address := request.GetString("deadline_address", ""); address != "" {
			activityData.ActivityDeadLine.Fields.Address = &address
		}

		if title := request.GetString("deadline_title", ""); title != "" {
			activityData.ActivityDeadLine.Fields.Title = &title
		}

		if content := request.GetString("deadline_content", ""); content != "" {
			activityData.ActivityDeadLine.Fields.Content = &content
		}

		if profileID := request.GetFloat("deadline_profile_id", 0); profileID > 0 {
			id := int64(profileID)
			activityData.ActivityDeadLine.Fields.ProfileID = &id
		}
	}

	// Handle comment data if provided
	if comment := request.GetString("comment", ""); comment != "" {
		activityData.ActivityComment = dto.ActivityCommentData{
			Fields: dto.ActivityCommentFields{
				Comment: comment,
			},
		}
	}

	// Call the usecase
	err = a.activityUseCase.EditActivity(req, activityData)
	if err != nil {
		return mcp.NewToolResultError("Failed to edit activity: " + err.Error()), nil
	}

	// Return success message
	message := a.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.EditedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)

	response := map[string]string{
		"message": message,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (a *ActivityMcp) handleDeleteActivity(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := a.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get activity ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build delete request
	deleteRequest := dto.DeleteRequest{
		ID: strconv.Itoa(int(id)),
	}

	// Call the usecase
	err = a.activityUseCase.DeleteActivity(req, deleteRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to delete activity: " + err.Error()), nil
	}

	// Return success message
	message := a.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.DeletedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)

	response := map[string]string{
		"message": message,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}