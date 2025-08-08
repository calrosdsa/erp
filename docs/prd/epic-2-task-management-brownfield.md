# Task Management for Project Module - Brownfield Enhancement

## Epic Title
**Task Management for Project Module - Brownfield Enhancement**

## Epic Goal
Add essential task management capabilities to the existing project module, enabling teams to create, assign, and track tasks within projects, providing immediate value with minimal disruption to the current system.

## Epic Description

### Existing System Context
- **Current relevant functionality:** Basic project CRUD operations (create, read, list) with name and status
- **Technology stack:** Go + Echo v4 + GORM + PostgreSQL + NATS event system
- **Integration points:** Existing project module at `project/projet_module/project/`, party system, and user management

### Enhancement Details
- **What's being added/changed:** Task management system with task creation, assignment, status tracking, and basic progress monitoring
- **How it integrates:** New `tasks` table linked to existing projects, leverages existing user/party system for assignments, follows established module patterns
- **Success criteria:** Users can create tasks within projects, assign to team members, track progress, and view task lists

## Stories

### Story 1: Create Task Management Database Schema and Models
- Add `tasks` table with project_id FK, assignee, title, description, status, priority, due_date
- Generate GORM models and queries using existing code generation patterns
- Create task DTOs following existing patterns

### Story 2: Implement Task CRUD Operations and REST API
- Create task handler, repository, and use case layers following project module structure
- Add REST endpoints: create task, get task, list project tasks, update task status
- Implement task assignment to existing users/parties

### Story 3: Add Task Status Tracking and Basic Progress View
- Implement task status transitions (To Do → In Progress → Done)
- Use kanban for task progress tracking
- Add task filtering and search capabilities within projects

### Story 4: Refactor Task Module Following CRM Deal Module Architecture
- **Repository Layer Refactoring:** Restructure task repository to match Deal repository patterns (`project/crm/deal/repository/pg_deal.go`)
- **Handler Layer Enhancement:** Refactor task handler to follow Deal handler structure with proper operation naming and documentation
- **Use Case Layer Alignment:** Update task use case to match Deal use case patterns with enhanced business logic
- **MCP Integration:** Implement comprehensive MCP handler following Deal MCP pattern with all CRUD operations
- **Module Integration:** Update module.go to follow Deal module's DI pattern and service registration
- **DTO Alignment:** Ensure task DTOs follow Deal DTO patterns for consistency across the system

## Compatibility Requirements
- [x] Existing APIs remain unchanged
- [x] Database schema changes are backward compatible
- [x] UI changes follow existing patterns
- [x] Performance impact is minimal

## Risk Mitigation
- **Primary Risk:** Database schema changes affecting existing project functionality
- **Mitigation:** Only adding new tables with FK constraints, no modifications to existing `projects` table
- **Rollback Plan:** Drop `tasks` table and remove task-related endpoints; existing project functionality unaffected

## Definition of Done
- [x] All stories completed with acceptance criteria met
- [x] Existing functionality verified through testing
- [x] Integration points working correctly
- [x] Documentation updated appropriately
- [x] No regression in existing features

## Story Manager Handoff

**Key Considerations for Story Development:**

- This is an enhancement to an existing Go-based ERP system running **Go + Echo v4 + GORM + PostgreSQL + NATS**
- Integration points: **Existing project module, user/party system for assignments, authentication/authorization middleware**
- Existing patterns to follow: **Handler→UseCase→Repository layers, GORM code generation, REST API with Huma v2, DTO patterns, event-driven architecture**
- Critical compatibility requirements: **No modifications to existing `projects` table, maintain existing project API contracts, follow established module structure in `project/projet_module/`**
- Each story must include verification that existing project functionality (list projects, create project, get project details) remains intact

The epic should maintain system integrity while delivering **essential task management capabilities that integrate seamlessly with the existing project module**.

## Validation Checklist

### Scope Validation
- [x] Epic can be completed in 3 stories maximum
- [x] No architectural documentation is required (follows existing patterns)
- [x] Enhancement follows existing patterns (handler→usecase→repository)
- [x] Integration complexity is manageable (FK relationships only)

### Risk Assessment
- [x] Risk to existing system is low (no existing table modifications)
- [x] Rollback plan is feasible (drop new table)
- [x] Testing approach covers existing functionality
- [x] Team has sufficient knowledge of integration points (standard GORM/REST patterns)

### Completeness Check
- [x] Epic goal is clear and achievable
- [x] Stories are properly scoped and sequential
- [x] Success criteria are measurable
- [x] Dependencies are identified (existing project/user systems)

## Why This Prioritization?

**Task Management** was chosen as the first epic because:

1. **Highest Business Value:** Tasks are the foundation of project management - without them, other features like milestones and time tracking are meaningless
2. **Lowest Risk:** Purely additive (new table + endpoints), zero impact on existing functionality
3. **Enables Future Features:** Creates the foundation for time tracking, milestones, and progress reporting
4. **Follows Existing Patterns:** Leverages established architecture patterns in the codebase
5. **Immediate User Value:** Teams can start organizing work within projects immediately

### Deferred Features (for future epics)
- **Milestones:** Requires tasks to be meaningful
- **Time Tracking:** Builds on task management
- **Calendar View:** Complex UI, better as separate epic
- **Team Collaboration:** Can be added incrementally
- **Progress Reporting:** Needs task data to report on

This epic delivers maximum value with minimal risk and sets the foundation for more advanced project management features.