package task_mcp

import (
	"context"
	"encoding/json"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	task_ucase "erp/project/projet_module/task/usecase"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type TaskMcp struct {
	mcpServer     *server.MCPServer
	sessionHelper helpers.SessionHelper
	locale        helpers.Locale
	errorHelper   helpers.ErrorHelper
	permission    repository.PermissionService
	usecase       task_ucase.TaskUseCase
}

func NewTaskMcp(
	mcpServer *server.MCPServer,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
	usecase task_ucase.TaskUseCase,
) {
	taskMcp := &TaskMcp{
		mcpServer:     mcpServer,
		sessionHelper: helpers.Session,
		locale:        helpers.Locale,
		errorHelper:   helpers.Error,
		permission:    permission,
		usecase:       usecase,
	}

	taskMcp.registerTools()
}

func (t *TaskMcp) registerTools() {
	// Get Tasks Tool
	getTasksToolspec := mcp.NewTool("get_tasks",
		mcp.WithDescription("Retrieve a paginated list of tasks with optional filtering and sorting capabilities"),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Number of tasks per page (default: 10)"),
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
		mcp.WithString("status",
			mcp.Description("Filter by task status (e.g., todo, in_progress, done)"),
		),
		mcp.WithNumber("assignee",
			mcp.Description("Filter by assignee ID"),
		),
		mcp.WithString("priority",
			mcp.Description("Filter by priority (e.g., low, medium, high, urgent)"),
		),
		mcp.WithString("search",
			mcp.Description("Search in title and description"),
		),
		mcp.WithString("due_date_from",
			mcp.Description("Due date range start (format: YYYY-MM-DD)"),
		),
		mcp.WithString("due_date_to",
			mcp.Description("Due date range end (format: YYYY-MM-DD)"),
		),
		mcp.WithString("created_from",
			mcp.Description("Created date range start (format: YYYY-MM-DD)"),  
		),
		mcp.WithString("created_to",
			mcp.Description("Created date range end (format: YYYY-MM-DD)"),
		),
	)

	t.mcpServer.AddTool(getTasksToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return t.handleGetTasks(ctx, request)
	})

	// Get Task Detail Tool
	getTaskToolspec := mcp.NewTool("get_task_detail",
		mcp.WithDescription("Retrieve comprehensive details for a specific task by ID"),
		mcp.WithNumber("id", 
			mcp.Required(),
			mcp.Description("The task ID to retrieve"),
		),
	)

	t.mcpServer.AddTool(getTaskToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return t.handleGetTask(ctx, request)
	})

	// Create Task Tool
	createTaskToolspec := mcp.NewTool("create_task",
		mcp.WithDescription("Create a new task record with the provided task data"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("Project ID that the task belongs to"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Task title"),
		),
		mcp.WithString("status",
			mcp.Description("Task status name (e.g., todo, in_progress, done)"),
		),
		mcp.WithNumber("stage_id",
			mcp.Description("Stage ID for the task (alternative to status)"),
		),
		mcp.WithNumber("assignee",
			mcp.Description("ID of the person assigned to the task"),
		),
		mcp.WithString("description",
			mcp.Description("Task description"),
		),
		mcp.WithString("priority",
			mcp.Description("Task priority (e.g., low, medium, high, urgent)"),
		),
		mcp.WithString("due_date",
			mcp.Description("Task due date (format: YYYY-MM-DD)"),
		),
		mcp.WithNumber("index",
			mcp.Description("Task index position within the stage (default: 1)"),
		),
	)

	t.mcpServer.AddTool(createTaskToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return t.handleCreateTask(ctx, request)
	})

	// Edit Task Tool
	editTaskToolspec := mcp.NewTool("edit_task",
		mcp.WithDescription("Update an existing task record with new data"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("Task ID to update"),
		),
		mcp.WithNumber("project_id",
			mcp.Description("Project ID that the task belongs to"),
		),
		mcp.WithString("title",
			mcp.Description("Task title"),
		),
		mcp.WithString("status",
			mcp.Description("Task status name (e.g., todo, in_progress, done)"),
		),
		mcp.WithNumber("stage_id",
			mcp.Description("Stage ID for the task"),
		),
		mcp.WithNumber("assignee",
			mcp.Description("ID of the person assigned to the task"),
		),
		mcp.WithString("description",
			mcp.Description("Task description"),
		),
		mcp.WithString("priority",
			mcp.Description("Task priority (e.g., low, medium, high, urgent)"),
		),
		mcp.WithString("due_date",
			mcp.Description("Task due date (format: YYYY-MM-DD)"),
		),
		mcp.WithNumber("index",
			mcp.Description("Task index position within the stage"),
		),
	)

	t.mcpServer.AddTool(editTaskToolspec, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return t.handleEditTask(ctx, request)
	})
}

func (t *TaskMcp) handleGetTasks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := t.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Build request parameters
	tasksRequest := dto.TasksRequest{}
	
	// Handle pagination via Size field
	if page := request.GetFloat("page", 1); page > 0 {
		tasksRequest.Size = strconv.Itoa(int(page) * 10) // Convert page to size format
	}
	
	if limit := request.GetFloat("limit", 10); limit > 0 {  
		tasksRequest.Size = strconv.Itoa(int(limit))
	} else if tasksRequest.Size == "" {
		tasksRequest.Size = "10" // Default size
	}

	// Handle optional filters
	if createdAt := request.GetString("created_at", ""); createdAt != "" {
		tasksRequest.CreatedAt = createdAt
	}
	
	if updatedAt := request.GetString("updated_at", ""); updatedAt != "" {
		tasksRequest.UpdatedAt = updatedAt
	}
	
	if sort := request.GetString("sort", ""); sort != "" {
		tasksRequest.OrderColumn = sort
	}
	
	if order := request.GetString("order", ""); order != "" {
		tasksRequest.Order = order
	}

	// Task-specific filters
	if status := request.GetString("status", ""); status != "" {
		tasksRequest.TaskFilterParams.Status = status
	}

	if assignee := request.GetFloat("assignee", 0); assignee > 0 {
		tasksRequest.TaskFilterParams.Assignee = int64(assignee)
	}

	if priority := request.GetString("priority", ""); priority != "" {
		tasksRequest.TaskFilterParams.Priority = priority
	}

	if search := request.GetString("search", ""); search != "" {
		tasksRequest.TaskFilterParams.Search = search
	}

	if dueDateFrom := request.GetString("due_date_from", ""); dueDateFrom != "" {
		tasksRequest.TaskFilterParams.DueDateFrom = dueDateFrom
	}

	if dueDateTo := request.GetString("due_date_to", ""); dueDateTo != "" {
		tasksRequest.TaskFilterParams.DueDateTo = dueDateTo
	}

	if createdFrom := request.GetString("created_from", ""); createdFrom != "" {
		tasksRequest.TaskFilterParams.CreatedFrom = createdFrom
	}

	if createdTo := request.GetString("created_to", ""); createdTo != "" {
		tasksRequest.TaskFilterParams.CreatedTo = createdTo
	}

	// Call the usecase
	result, err := t.usecase.GetTasks(req, tasksRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get tasks: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (t *TaskMcp) handleGetTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := t.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get task ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build request
	taskRequest := dto.RequestEntity{
		ID: strconv.Itoa(int(id)),
	}

	// Call the usecase
	result, err := t.usecase.GetTask(req, taskRequest)
	if err != nil {
		return mcp.NewToolResultError("Failed to get task: " + err.Error()), nil
	}

	// Return the result as JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (t *TaskMcp) handleCreateTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := t.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Extract required fields
	projectID, err := request.RequireFloat("project_id")
	if err != nil {
		return mcp.NewToolResultError("Project ID is required: " + err.Error()), nil
	}

	title, err := request.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError("Title is required: " + err.Error()), nil
	}

	// Build task data
	taskData := dto.TaskData{
		Fields: dto.TaskFields{
			ProjectID: int64(projectID),
			Title:     title,
			Index:     1, // Default index
		},
	}

	// Handle status or stage_id - one of them should be provided
	// if status := request.GetString("status", ""); status != "" {
	// 	taskData.Fields.Status = status
	// }
	
	// if stageID := request.GetFloat("stage_id", 0); stageID > 0 {
	// 	taskData.Fields.StageID = int32(stageID)
	// }

	// // Validate that either status or stage_id is provided
	// if taskData.Fields.Status == "" && taskData.Fields.StageID == 0 {
	// 	return mcp.NewToolResultError("Either status or stage_id must be provided"), nil
	// }

	// Handle optional fields
	if assignee := request.GetFloat("assignee", 0); assignee > 0 {
		id := int64(assignee)
		taskData.Fields.Assignee = &id
	}

	if description := request.GetString("description", ""); description != "" {
		taskData.Fields.Description = &description
	}

	if priority := request.GetString("priority", ""); priority != "" {
		taskData.Fields.Priority = &priority
	}

	if dueDateStr := request.GetString("due_date", ""); dueDateStr != "" {
		dueDate, parseErr := time.Parse("2006-01-02", dueDateStr)
		if parseErr != nil {
			return mcp.NewToolResultError("Invalid due date format (expected YYYY-MM-DD): " + parseErr.Error()), nil
		}
		taskData.Fields.DueDate = &dueDate
	}

	if index := request.GetFloat("index", 0); index > 0 {
		taskData.Fields.Index = int32(index)
	}

	// Call the usecase
	result, err := t.usecase.CreateTask(req, taskData)
	if err != nil {
		return mcp.NewToolResultError("Failed to create task: " + err.Error()), nil
	}

	// Return success message with result
	response := map[string]interface{}{
		"message": "Task created successfully",
		"task":    result,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (t *TaskMcp) handleEditTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get session for authentication
	req, err := t.sessionHelper.GetSession(ctx)
	if err != nil {
		return mcp.NewToolResultError("Not Authorized: " + err.Error()), nil
	}

	// Get task ID
	id, err := request.RequireFloat("id")
	if err != nil {
		return mcp.NewToolResultError("ID is required: " + err.Error()), nil
	}

	// Build task data with ID
	taskData := dto.TaskData{
		ID:     int64(id),
		Fields: dto.TaskFields{},
	}

	// Handle optional fields - only update fields that are provided
	if projectID := request.GetFloat("project_id", 0); projectID > 0 {
		taskData.Fields.ProjectID = int64(projectID)
	}

	if title := request.GetString("title", ""); title != "" {
		taskData.Fields.Title = title
	}

	// if status := request.GetString("status", ""); status != "" {
	// 	taskData.Fields.Status = status
	// }

	if stageID := request.GetFloat("stage_id", 0); stageID > 0 {
		taskData.Fields.StageID = int32(stageID)
	}

	if assignee := request.GetFloat("assignee", 0); assignee > 0 {
		id := int64(assignee)
		taskData.Fields.Assignee = &id
	}

	if description := request.GetString("description", ""); description != "" {
		taskData.Fields.Description = &description
	}

	if priority := request.GetString("priority", ""); priority != "" {
		taskData.Fields.Priority = &priority
	}

	if dueDateStr := request.GetString("due_date", ""); dueDateStr != "" {
		dueDate, parseErr := time.Parse("2006-01-02", dueDateStr)
		if parseErr != nil {
			return mcp.NewToolResultError("Invalid due date format (expected YYYY-MM-DD): " + parseErr.Error()), nil
		}
		taskData.Fields.DueDate = &dueDate
	}

	if index := request.GetFloat("index", 0); index > 0 {
		taskData.Fields.Index = int32(index)
	}

	// Call the usecase
	err = t.usecase.EditTask(req, taskData)
	if err != nil {
		return mcp.NewToolResultError("Failed to edit task: " + err.Error()), nil
	}

	// Return success message
	message := t.locale.MustLocalize(
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