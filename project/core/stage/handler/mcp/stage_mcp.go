package stage_mcp

import (
	"context"
	"encoding/json"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	stage_ucase "erp/project/core/stage/usecase"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type StageMcp struct {
	mcpServer     *server.MCPServer
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       stage_ucase.StageUseCase
}

func NewStageMcp(
	mcpServer *server.MCPServer,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	usecase stage_ucase.StageUseCase,
) {
	stageMcp := &StageMcp{
		mcpServer:     mcpServer,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		permission:    permission,
		usecase:       usecase,
	}

	stageMcp.registerTools()
}

func (s *StageMcp) registerTools() {
	// Get Stages Tool
	getStagesToolspec := mcp.NewTool("get_stages",
		mcp.WithDescription("Retrieve a list of stages with optional filtering capabilities"),
		mcp.WithString("entity_id",
			mcp.Required(),
			mcp.Description("Entity ID to filter stages"),
		),
		mcp.WithString("name",
			mcp.Description("Filter by stage name"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Number of stages per page (default: 10)"),
		),
		mcp.WithString("sort",
			mcp.Description("Sort field"),
		),
		mcp.WithString("order",
			mcp.Description("Sort order (asc/desc)"),
		),
	)

	s.mcpServer.AddTool(getStagesToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleGetStages(ctx, request)
	})

	// Create Stage Tool
	createStageToolspec := mcp.NewTool("create_stage",
		mcp.WithDescription("Create a new stage with the provided stage data"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Stage name"),
		),
		mcp.WithNumber("entity_id",
			mcp.Required(),
			mcp.Description("Entity ID for the stage"),
		),
		mcp.WithString("color",
			mcp.Required(),
			mcp.Description("Stage color (hex color code)"),
		),
		mcp.WithNumber("index",
			mcp.Required(),
			mcp.Description("Stage index position"),
		),
	)

	s.mcpServer.AddTool(createStageToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleCreateStage(ctx, request)
	})

	// Edit Stage Tool
	editStageToolspec := mcp.NewTool("edit_stage",
		mcp.WithDescription("Update an existing stage with new data"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("Stage ID to update"),
		),
		mcp.WithString("name",
			mcp.Description("Stage name"),
		),
		mcp.WithNumber("entity_id",
			mcp.Description("Entity ID for the stage"),
		),
		mcp.WithString("color",
			mcp.Description("Stage color (hex color code)"),
		),
		mcp.WithNumber("index",
			mcp.Description("Stage index position"),
		),
	)

	s.mcpServer.AddTool(editStageToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleEditStage(ctx, request)
	})

	// Stage Transition Tool
	stageTransitionToolspec := mcp.NewTool("stage_transition",
		mcp.WithDescription("Perform a stage transition by reordering stages"),
		mcp.WithNumber("source_id",
			mcp.Required(),
			mcp.Description("Source stage ID"),
		),
		mcp.WithNumber("source_index",
			mcp.Required(),
			mcp.Description("Source index position"),
		),
		mcp.WithNumber("destination_id",
			mcp.Required(),
			mcp.Description("Destination stage ID"),
		),
		mcp.WithNumber("destination_index",
			mcp.Required(),
			mcp.Description("Destination index position"),
		),
	)

	s.mcpServer.AddTool(stageTransitionToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleStageTransition(ctx, request)
	})

	// Delete Stage Tool
	deleteStageToolspec := mcp.NewTool("delete_stage",
		mcp.WithDescription("Delete a stage by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Stage ID to delete"),
		),
	)

	s.mcpServer.AddTool(deleteStageToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleDeleteStage(ctx, request)
	})
}

func (s *StageMcp) handleGetStages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := s.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get required entity_id
	entityID, err := request.RequireString("entity_id")
	if err != nil {
		return mcp.NewToolResultError("Entity ID is required: " + err.Error()), nil
	}

	// Build request parameters
	stagesRequest := dto.StagesRequest{
		EntityID: entityID,
	}

	// Handle pagination
	if page := request.GetFloat("page", 1); page > 0 {
		stagesRequest.Size = strconv.Itoa(int(page) * 10) // Convert page to size format
	}
	
	if limit := request.GetFloat("limit", 10); limit > 0 {  
		stagesRequest.Size = strconv.Itoa(int(limit))
	} else if stagesRequest.Size == "" {
		stagesRequest.Size = "10" // Default size
	}

	// Handle optional filters
	if name := request.GetString("name", ""); name != "" {
		stagesRequest.Name = name
	}
	
	if sort := request.GetString("sort", ""); sort != "" {
		stagesRequest.OrderColumn = sort
	}
	
	if order := request.GetString("order", ""); order != "" {
		stagesRequest.Order = order
	}

	// Call the usecase
	result, err := s.usecase.GetStages(req, stagesRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get stages: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (s *StageMcp) handleCreateStage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := s.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Extract required fields
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("Name is required: " + err.Error()), nil
	}

	entityID, err := request.RequireFloat("entity_id")
	if err != nil {
		return mcp.NewToolResultError("Entity ID is required: " + err.Error()), nil
	}

	color, err := request.RequireString("color")
	if err != nil {
		return mcp.NewToolResultError("Color is required: " + err.Error()), nil
	}

	index, err := request.RequireFloat("index")
	if err != nil {
		return mcp.NewToolResultError("Index is required: " + err.Error()), nil
	}

	// Build stage data
	stageData := dto.StageData{
		Fields: dto.StageFields{
			Name:     name,
			EntityID: int32(entityID),
			Color:    color,
			Index:    int32(index),
		},
	}

	// Call the usecase
	result, err := s.usecase.CreateStage(req, stageData)
	if err != nil {
		return mcp.NewToolResultError("Failed to create stage: " + err.Error()), nil
	}

	// Return success message with result
	message := s.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreatedSuccessfully"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)

	response := map[string]interface{}{
		"message": message,
		"stage":   result,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (s *StageMcp) handleEditStage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := s.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get stage ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build stage data with ID
	stageData := dto.StageData{
		ID: int32(id),
		Fields: dto.StageFields{
			Index: 1, // Default index if not provided
		},
	}

	// Handle optional fields - only update fields that are provided
	if name := request.GetString("name", ""); name != "" {
		stageData.Fields.Name = name
	}

	if entityID := request.GetFloat("entity_id", 0); entityID > 0 {
		stageData.Fields.EntityID = int32(entityID)
	}

	if color := request.GetString("color", ""); color != "" {
		stageData.Fields.Color = color
	}

	if index := request.GetFloat("index", 0); index > 0 {
		stageData.Fields.Index = int32(index)
	}

	// Call the usecase
	err = s.usecase.EditStage(req, stageData)
	if err != nil {
		return mcp.NewToolResultError("Failed to edit stage: " + err.Error()), nil
	}

	// Return success message
	message := s.locale.MustLocalize(
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

func (s *StageMcp) handleStageTransition(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := s.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get required fields
	sourceID, err := request.RequireFloat("source_id")
	if err != nil {
		return mcp.NewToolResultError("Source ID is required: " + err.Error()), nil
	}

	sourceIndex, err := request.RequireFloat("source_index")
	if err != nil {
		return mcp.NewToolResultError("Source index is required: " + err.Error()), nil
	}

	destinationID, err := request.RequireFloat("destination_id")
	if err != nil {
		return mcp.NewToolResultError("Destination ID is required: " + err.Error()), nil
	}

	destinationIndex, err := request.RequireFloat("destination_index")
	if err != nil {
		return mcp.NewToolResultError("Destination index is required: " + err.Error()), nil
	}

	// Build transition request
	transitionRequest := dto.StageTransitionData{
		SourceID:         int32(sourceID),
		SourceIndex:      int32(sourceIndex),
		DestinationID:    int32(destinationID),
		DestinationIndex: int32(destinationIndex),
	}

	// Call the usecase
	err = s.usecase.StageTransition(req, transitionRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to transition stage: " + err.Error()), nil
	}

	// Return success message
	message := s.locale.MustLocalize(
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

func (s *StageMcp) handleDeleteStage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := s.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get stage ID
	idStr, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build delete request
	deleteRequest := dto.DeleteRequest{
		ID: idStr,
	}

	// Call the usecase
	err = s.usecase.DeleteStage(req, &deleteRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to delete stage: " + err.Error()), nil
	}

	// Return success message
	message := s.locale.MustLocalize(
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