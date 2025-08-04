package deal_mcp

import (
	"context"
	"encoding/json"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	deal_ucase "erp/project/crm/deal/usecase"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DealMcp struct {
	mcpServer     *server.MCPServer
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       deal_ucase.DealUseCase
}

func NewDealMcp(
	mcpServer *server.MCPServer,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	usecase deal_ucase.DealUseCase,
) {
	dealMcp := &DealMcp{
		mcpServer:     mcpServer,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		permission:    permission,
		usecase:       usecase,
	}

	dealMcp.registerTools()
}

func (d *DealMcp) registerTools() {
	// Get Deals Tool
	getDealsToolspec := mcp.NewTool("get_deals",
		mcp.WithDescription("Retrieve a paginated list of deals with optional filtering and sorting capabilities"),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Number of deals per page (default: 10)"),
		),
		mcp.WithString("created_at",
			mcp.Description("Filter by creation date (format: YYYY-MM-DD)"),
		),
		mcp.WithString("updated_at",
			mcp.Description("Filter by update date (format: YYYY-MM-DD)"),
		),
		mcp.WithString("sort",
			mcp.Description("Sort field"),
		),
		mcp.WithString("order",
			mcp.Description("Sort order (asc/desc)"),
		),
	)

	d.mcpServer.AddTool(getDealsToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return d.handleGetDeals(ctx, request)
	})

	// Get Deal Detail Tool
	getDealToolspec := mcp.NewTool("get_deal_detail",
		mcp.WithDescription("Retrieve comprehensive details for a specific deal by ID"),
		mcp.WithNumber("id", 
			mcp.Required(),
			mcp.Description("The deal ID to retrieve"),
		),
	)

	d.mcpServer.AddTool(getDealToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return d.handleGetDeal(ctx, request)
	})

	// Create Deal Tool
	createDealToolspec := mcp.NewTool("create_deal",
		mcp.WithDescription("Create a new deal record with the provided deal data"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Deal name"),
		),
		mcp.WithNumber("stage_id",
			mcp.Required(),
			mcp.Description("Stage ID for the deal"),
		),
		mcp.WithNumber("amount",
			mcp.Required(),
			mcp.Description("Deal amount"),
		),
		mcp.WithString("currency",
			mcp.Required(),
			mcp.Description("Currency code (e.g., USD, EUR)"),
		),
		mcp.WithNumber("responsible_id",
			mcp.Required(),
			mcp.Description("ID of the responsible person"),
		),
		mcp.WithString("start_date",
			mcp.Required(),
			mcp.Description("Deal start date (format: YYYY-MM-DD)"),
		),
		mcp.WithString("deal_type",
			mcp.Description("Type of the deal (e.g., Ventas, Servicios, Postventa, Ventas Integradas,Ventas de Mercancías)"),
		),
		mcp.WithString("source",
			mcp.Description("Deal source"),
		),
		mcp.WithNumber("customer_id",
			mcp.Description("Customer ID associated with the deal"),
		),
		mcp.WithBoolean("available_for_everyone",
			mcp.Description("Whether the deal is available for everyone (default: false)"),
		),
	)

	d.mcpServer.AddTool(createDealToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return d.handleCreateDeal(ctx, request)
	})

	// Edit Deal Tool
	editDealToolspec := mcp.NewTool("edit_deal",
		mcp.WithDescription("Update an existing deal record with new data"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("Deal ID to update"),
		),
		mcp.WithString("name",
			mcp.Description("Deal name"),
		),
		mcp.WithNumber("stage_id",
			mcp.Description("Stage ID for the deal"),
		),
		mcp.WithNumber("amount",
			mcp.Description("Deal amount"),
		),
		mcp.WithString("currency",
			mcp.Description("Currency code (e.g., USD, EUR)"),
		),
		mcp.WithNumber("responsible_id",
			mcp.Description("ID of the responsible person"),
		),
		mcp.WithString("start_date",
			mcp.Description("Deal start date (format: YYYY-MM-DD)"),
		),
		mcp.WithString("deal_type",
			mcp.Description("Type of the deal"),
		),
		mcp.WithString("source",
			mcp.Description("Deal source"),
		),
		mcp.WithNumber("customer_id",
			mcp.Description("Customer ID associated with the deal"),
		),
		mcp.WithBoolean("available_for_everyone",
			mcp.Description("Whether the deal is available for everyone"),
		),
	)

	d.mcpServer.AddTool(editDealToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return d.handleEditDeal(ctx, request)
	})

	// Deal Transition Tool
	transitionDealToolspec := mcp.NewTool("transition_deal",
		mcp.WithDescription("Transition a deal between stages"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("Deal ID to transition"),
		),
		mcp.WithNumber("source_stage_id",
			mcp.Required(),
			mcp.Description("Source stage ID"),
		),
		mcp.WithNumber("source_index",
			mcp.Required(),
			mcp.Description("Source index position"),
		),
		mcp.WithNumber("destination_stage_id",
			mcp.Required(), 
			mcp.Description("Destination stage ID"),
		),
		mcp.WithNumber("destination_index",
			mcp.Required(),
			mcp.Description("Destination index position"),
		),
		mcp.WithString("source_name",
			mcp.Description("Source stage name"),
		),
		mcp.WithString("destination_name",
			mcp.Description("Destination stage name"),
		),
	)

	d.mcpServer.AddTool(transitionDealToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return d.handleDealTransition(ctx, request)
	})
}

func (d *DealMcp) handleGetDeals(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := d.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Build request parameters
	dealsRequest := dto.DealsRequest{}
	
	// Handle pagination via Size field
	if page := request.GetFloat("page", 1); page > 0 {
		dealsRequest.Size = strconv.Itoa(int(page) * 10) // Convert page to size format
	}
	
	if limit := request.GetFloat("limit", 10); limit > 0 {  
		dealsRequest.Size = strconv.Itoa(int(limit))
	} else if dealsRequest.Size == "" {
		dealsRequest.Size = "10" // Default size
	}

	// Handle optional filters
	if createdAt := request.GetString("created_at", ""); createdAt != "" {
		dealsRequest.CreatedAt = createdAt
	}
	
	if updatedAt := request.GetString("updated_at", ""); updatedAt != "" {
		dealsRequest.UpdatedAt = updatedAt
	}
	
	if sort := request.GetString("sort", ""); sort != "" {
		dealsRequest.OrderColumn = sort
	}
	
	if order := request.GetString("order", ""); order != "" {
		dealsRequest.Order = order
	}

	// Call the usecase
	result, err := d.usecase.GetDeals(req, dealsRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get deals: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (d *DealMcp) handleGetDeal(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := d.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get deal ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build request
	dealRequest := dto.RequestEntity{
		ID: strconv.Itoa(int(id)),
	}

	// Call the usecase
	result, err := d.usecase.GetDeal(req, dealRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get deal: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (d *DealMcp) handleCreateDeal(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := d.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Extract required fields
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("Name is required: " + err.Error()), nil
	}

	stageID, err := request.RequireFloat("stage_id")
	if err != nil {
		return mcp.NewToolResultError("Stage ID is required: " + err.Error()), nil
	}

	amount, err := request.RequireFloat("amount")
	if err != nil {
		return mcp.NewToolResultError("Amount is required: " + err.Error()), nil
	}

	currency, err := request.RequireString("currency")
	if err != nil {
		return mcp.NewToolResultError("Currency is required: " + err.Error()), nil
	}

	responsibleID, err := request.RequireFloat("responsible_id")
	if err != nil {
		return mcp.NewToolResultError("Responsible ID is required: " + err.Error()), nil
	}

	startDateStr, err := request.RequireString("start_date")
	if err != nil {
		return mcp.NewToolResultError("Start date is required: " + err.Error()), nil
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return mcp.NewToolResultError("Invalid start date format (expected YYYY-MM-DD): " + err.Error()), nil
	}

	// Build deal data
	dealData := dto.DealData{
		Fields: dto.DealFields{
			Name:          name,
			StageID:       int64(stageID),
			Amount:        int64(amount),
			Currency:      currency,
			ResponsibleID: int64(responsibleID),
			StartDate:     startDate,
			Index:         1, // Default index
		},
	}

	// Handle optional fields
	if dealType := request.GetString("deal_type", ""); dealType != "" {
		dealData.Fields.DealType = &dealType
	}

	if source := request.GetString("source", ""); source != "" {
		dealData.Fields.Source = &source
	}

	if customerID := request.GetFloat("customer_id", 0); customerID > 0 {
		id := int64(customerID)
		dealData.Fields.CustomerID = &id
	}

	if availableForEveryone := request.GetBool("available_for_everyone", false); availableForEveryone {
		dealData.Fields.AvailableForEveryone = availableForEveryone
	}

	// Call the usecase
	result, err := d.usecase.CreateDeal(req, dealData)
	if err != nil {
		return mcp.NewToolResultError("Failed to create deal: " + err.Error()), nil
	}

	// Return success message with result
	response := map[string]interface{}{
		"message": "Deal created successfully",
		"deal":    result,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (d *DealMcp) handleEditDeal(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := d.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get deal ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build deal data with ID
	dealData := dto.DealData{
		ID: int64(id),
		Fields: dto.DealFields{
			Index: 1, // Default index if not provided
		},
	}

	// Handle optional fields - only update fields that are provided
	if name := request.GetString("name", ""); name != "" {
		dealData.Fields.Name = name
	}

	if stageID := request.GetFloat("stage_id", 0); stageID > 0 {
		dealData.Fields.StageID = int64(stageID)
	}

	if amount := request.GetFloat("amount", 0); amount > 0 {
		dealData.Fields.Amount = int64(amount)
	}

	if currency := request.GetString("currency", ""); currency != "" {
		dealData.Fields.Currency = currency
	}

	if responsibleID := request.GetFloat("responsible_id", 0); responsibleID > 0 {
		dealData.Fields.ResponsibleID = int64(responsibleID)
	}

	if startDateStr := request.GetString("start_date", ""); startDateStr != "" {
		startDate, parseErr := time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			return mcp.NewToolResultError("Invalid start date format (expected YYYY-MM-DD): " + parseErr.Error()), nil
		}
		dealData.Fields.StartDate = startDate
	}

	if dealType := request.GetString("deal_type", ""); dealType != "" {
		dealData.Fields.DealType = &dealType
	}

	if source := request.GetString("source", ""); source != "" {
		dealData.Fields.Source = &source
	}

	if customerID := request.GetFloat("customer_id", 0); customerID > 0 {
		id := int64(customerID)
		dealData.Fields.CustomerID = &id
	}

	if availableForEveryone := request.GetBool("available_for_everyone", false); availableForEveryone {
		dealData.Fields.AvailableForEveryone = availableForEveryone
	}

	// Call the usecase
	err = d.usecase.EditDeal(req, dealData)
	if err != nil {
		return mcp.NewToolResultError("Failed to edit deal: " + err.Error()), nil
	}

	// Return success message
	message := d.locale.MustLocalize(
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

func (d *DealMcp) handleDealTransition(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := d.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get required fields
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	sourceStageID, err := request.RequireFloat("source_stage_id")
	if err != nil {
		return mcp.NewToolResultError("Source stage ID is required: " + err.Error()), nil
	}

	sourceIndex, err := request.RequireFloat("source_index")
	if err != nil {
		return mcp.NewToolResultError("Source index is required: " + err.Error()), nil
	}

	destinationStageID, err := request.RequireFloat("destination_stage_id")
	if err != nil {
		return mcp.NewToolResultError("Destination stage ID is required: " + err.Error()), nil
	}

	destinationIndex, err := request.RequireFloat("destination_index")
	if err != nil {
		return mcp.NewToolResultError("Destination index is required: " + err.Error()), nil
	}

	// Build transition request
	transitionRequest := dto.EntityTransitionData{
		ID:                 int64(id),
		SourceStageID:      int64(sourceStageID),
		SourceIndex:        int32(sourceIndex),
		DestinationStageID: int64(destinationStageID),
		DestinationIndex:   int32(destinationIndex),
	}

	// Handle optional fields
	if sourceName := request.GetString("source_name", ""); sourceName != "" {
		transitionRequest.SourceName = sourceName
	}

	if destinationName := request.GetString("destination_name", ""); destinationName != "" {
		transitionRequest.DestionationName = destinationName // Note: typo in DTO field name
	}

	// Call the usecase
	err = d.usecase.DealTransition(req, transitionRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to transition deal: " + err.Error()), nil
	}

	// Return success message
	message := d.locale.MustLocalize(
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