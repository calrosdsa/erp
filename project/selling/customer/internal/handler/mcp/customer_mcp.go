package customer_mcp

import (
	"context"
	"encoding/json"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	customer_ucase "erp/project/selling/customer/internal/usecase"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CustomerMcp struct {
	mcpServer     *server.MCPServer
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       customer_ucase.CustomerUseCase
}

func NewCustomerMcp(
	mcpServer *server.MCPServer,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	usecase customer_ucase.CustomerUseCase,
) {
	customerMcp := &CustomerMcp{
		mcpServer:     mcpServer,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		permission:    permission,
		usecase:       usecase,
	}

	customerMcp.registerTools()
}

func (c *CustomerMcp) registerTools() {
	// Get Customers Tool
	getCustomersToolspec := mcp.NewTool("get_customers",
		mcp.WithDescription("Retrieve a paginated list of customers with optional filtering and sorting capabilities"),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Number of customers per page (default: 10)"),
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
		mcp.WithString("search",
			mcp.Description("Search term for customer name or other fields"),
		),
	)

	c.mcpServer.AddTool(getCustomersToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return c.handleGetCustomers(ctx, request)
	})

	// Get Customer Detail Tool
	getCustomerToolspec := mcp.NewTool("get_customer_detail",
		mcp.WithDescription("Retrieve comprehensive details for a specific customer by ID"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The customer ID to retrieve"),
		),
	)

	c.mcpServer.AddTool(getCustomerToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return c.handleGetCustomer(ctx, request)
	})

	// Create Customer Tool
	createCustomerToolspec := mcp.NewTool("create_customer",
		mcp.WithDescription("Create a new customer record with the provided customer data"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Customer name (1-50 characters)"),
		),
		mcp.WithString("customer_type",
			mcp.Required(),
			mcp.Description("Customer type (1-50 characters)"),
		),
		mcp.WithNumber("group_id",
			mcp.Description("Group ID for the customer (optional)"),
		),
		mcp.WithArray("contacts",
			mcp.Description("Array of contact data for the customer (optional)"),
		),
	)

	c.mcpServer.AddTool(createCustomerToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return c.handleCreateCustomer(ctx, request)
	})

	// Edit Customer Tool
	editCustomerToolspec := mcp.NewTool("edit_customer",
		mcp.WithDescription("Update an existing customer record with new data"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("Customer ID to update"),
		),
		mcp.WithString("name",
			mcp.Description("Customer name (1-50 characters)"),
		),
		mcp.WithString("customer_type",
			mcp.Description("Customer type (1-50 characters)"),
		),
		mcp.WithNumber("group_id",
			mcp.Description("Group ID for the customer"),
		),
		mcp.WithArray("contacts",
			mcp.Description("Array of contact data for the customer"),
		),
	)

	c.mcpServer.AddTool(editCustomerToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return c.handleEditCustomer(ctx, request)
	})

	// Get Customer Types Tool
	getCustomerTypesToolspec := mcp.NewTool("get_customer_types",
		mcp.WithDescription("Retrieve available customer types for creating or updating customers"),
	)

	c.mcpServer.AddTool(getCustomerTypesToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return c.handleGetCustomerTypes(ctx, request)
	})
}

func (c *CustomerMcp) handleGetCustomers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := c.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Build request parameters
	customersRequest := dto.RequestPaginationData{}

	// Handle pagination
	if page := request.GetFloat("page", 1); page > 0 {
		customersRequest.Page = strconv.Itoa(int(page))
	} else {
		customersRequest.Page = "1" // Default page
	}

	if limit := request.GetFloat("limit", 10); limit > 0 {
		customersRequest.Size = strconv.Itoa(int(limit))
	} else {
		customersRequest.Size = "10" // Default size
	}

	// Handle optional filters
	// if createdAt := request.GetString("created_at", ""); createdAt != "" {
	// 	customersRequest.CreatedAt = createdAt
	// }

	// if updatedAt := request.GetString("updated_at", ""); updatedAt != "" {
	// 	customersRequest.UpdatedAt = updatedAt
	// }

	if sort := request.GetString("sort", ""); sort != "" {
		customersRequest.OrderColumn = sort
	}

	if order := request.GetString("order", ""); order != "" {
		customersRequest.Order = order
	}

	if query := request.GetString("query", ""); query != "" {
		customersRequest.Query = query
	}

	// Call the usecase
	result, err := c.usecase.GetCustomers(req, &customersRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get customers: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (c *CustomerMcp) handleGetCustomer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := c.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get customer ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build request
	customerRequest := dto.RequestEntity{
		ID: strconv.Itoa(int(id)),
	}

	// Call the usecase
	result, err := c.usecase.GetCustomerDetail(req, &customerRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get customer: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (c *CustomerMcp) handleCreateCustomer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := c.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Extract required fields
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("Name is required: " + err.Error()), nil
	}

	customerType, err := request.RequireString("customer_type")
	if err != nil {
		return mcp.NewToolResultError("Customer type is required: " + err.Error()), nil
	}

	// Build customer data
	customerData := dto.CustomerData{
		Fields: dto.CustomerFields{
			Name:         name,
			CustomerType: customerType,
		},
		Contacts: []dto.ContactData{}, // Default empty contacts
	}

	// Handle optional fields
	if groupID := request.GetFloat("group_id", 0); groupID > 0 {
		id := int64(groupID)
		customerData.Fields.GroupID = &id
	}

	// Handle contacts array if provided - for now keep empty as structure needs proper parsing
	customerData.Contacts = []dto.ContactData{}

	// Call the usecase
	result, err := c.usecase.CreateCustomer(req, customerData)
	if err != nil {
		return mcp.NewToolResultError("Failed to create customer: " + err.Error()), nil
	}

	// Return success message with result
	message := c.locale.MustLocalize(
		helpers.OptionsLocale.WithID("Message.CreateCustomerSuccess"),
		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
	)

	response := map[string]interface{}{
		"message":  message,
		"customer": result,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (c *CustomerMcp) handleEditCustomer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := c.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get customer ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build customer data with ID
	customerData := dto.CustomerData{
		ID: int64(id),
		Fields: dto.CustomerFields{
			// Only update fields that are provided
		},
		Contacts: []dto.ContactData{}, // Default empty contacts
	}

	// Handle optional fields - only update fields that are provided
	if name := request.GetString("name", ""); name != "" {
		customerData.Fields.Name = name
	}

	if customerType := request.GetString("customer_type", ""); customerType != "" {
		customerData.Fields.CustomerType = customerType
	}

	if groupID := request.GetFloat("group_id", 0); groupID > 0 {
		id := int64(groupID)
		customerData.Fields.GroupID = &id
	}

	// Handle contacts array if provided - for now keep empty as structure needs proper parsing
	customerData.Contacts = []dto.ContactData{}

	// Call the usecase
	err = c.usecase.EditCustomer(req, customerData)
	if err != nil {
		return mcp.NewToolResultError("Failed to edit customer: " + err.Error()), nil
	}

	// Return success message
	message := c.locale.MustLocalize(
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

func (c *CustomerMcp) handleGetCustomerTypes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := c.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Call the usecase
	result := c.usecase.GetCustomerTypes(req)

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}