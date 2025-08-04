package company_service

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/domain"
	"erp/internal/app/entity"
	"erp/internal/app/service/helpers"
	userservice "erp/internal/app/service/services/user_service"
	logger "erp/pkg/logger"
	"erp/pkg/permission"
	"time"

	"gorm.io/gorm"
)

type CompanyService struct {
	timeout           *time.Duration
	conn              *connection.Connection
	convertor         helpers.ConvertorHelper
	emitLog           logger.EmitLog
	generator         helpers.Generator
	userService       *userservice.UserService
	permissionService permission.PermissionService
}

func NewCompanyService(conn *connection.Connection, timeout *time.Duration,
	helpers *helpers.Helpers,
	userService *userservice.UserService,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *CompanyService {
	return &CompanyService{
		conn:              conn,
		timeout:           timeout,
		convertor:         helpers.Convertor,
		emitLog:           logger.EmitLog("company-service"),
		generator:         helpers.Generator,
		userService:       userService,
		permissionService: permissionService,
	}
}

func (s *CompanyService) CreateCompany(req *common.RequestContext, d *dto.CreateCompanyRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	tx := s.conn.Db.Begin()
	var err error
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCompany"))
			tx.Rollback()
		}
	}()
	var company model.Company
	company.Name = d.Body.Name
	if d.Body.ParentID != nil {
		company.ParentID = d.Body.ParentID
		company.Ordinal = int32(entity.DEPTH_TWO)
	} else {
		parentCompany, err := s.permissionService.GetParentCompany(req, tx)
		if err != nil {
			return err
		}
		company.ParentID = &req.ActiveCompany.ID
		company.Ordinal = parentCompany.Ordinal + 1
	}

	// code := s.generateCompanyCode(tx)
	company.Code = s.generateCompanyCode(tx)
	err = tx.Save(&company).Error
	if err != nil {
		return err
	}
	// err = s.addCompanyDepartment(ctx, tx, *company.ParentID, uint(company.ID))
	// if err != nil {
	// 	return err
	// }
	err = s.userService.InsertUserCompany(req, ctx, tx, company.ID)
	if err != nil {
		return err
	}
	return tx.Commit().Error
}

func (s *CompanyService) generateCompanyCode(db *gorm.DB) (code string) {
	for {
		code = s.generator.GenerateCode()
		// Check if the code already exists in the database
		var count int64
		err := db.Model(&model.Company{}).
			Where("code = ?", code).
			Count(&count).Error
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("generateCompanyCode"))
		}
		// If the code is unique, break the loop
		if count == 0 {
			break
		}
	}
	return
}

func (s *CompanyService) GetValidaParentCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (dto.PaginationResult[[]dto.CompanyDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var result dto.PaginationResult[[]dto.CompanyDto]
	var companies []*model.Company
	queryBuilder := s.conn.Db.WithContext(ctx).Table("companies").
		Joins("JOIN user_relations ON user_relations.company_id = companies.id").
		Where("user_relations.user_id = ? and ordinal = ?", req.User.ID, entity.SUB_PARENT_ORDINAL)

	err := queryBuilder.
		Count(&result.Total).Error

	if d.Query != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	}

	err = queryBuilder.
		Scopes(s.conn.Paginate(req.Params)).
		Find(&companies).Error
	if err != nil {
		return result, err
	}

	companiesDto := make([]dto.CompanyDto, len(companies))
	for i, company := range companies {
		companiesDto[i] = dto.CompanyDTOFromModel(company)
	}
	result.Results = companiesDto

	return result, err
}

func (s *CompanyService) GetAllUserCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (dto.PaginationResult[[]dto.CompanyDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var result dto.PaginationResult[[]dto.CompanyDto]
	if allow := s.permissionService.CheckPermission(ctx, req, domain.COMPANY, domain.VIEW); !allow {
		return result, domain.ACTION_NOT_ALLOWED
	}
	var companies []*model.Company
	queryBuilder := s.conn.Db.WithContext(ctx).Table("companies").
		Joins("JOIN user_relations ON user_relations.company_id = companies.id").
		Where("user_relations.user_id = ?", req.User.ID)

	err := queryBuilder.
		Count(&result.Total).Error

	if d.Query != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	}

	err = queryBuilder.
		Scopes(s.conn.Paginate(req.Params)).
		Find(&companies).Error
	if err != nil {
		return result, err
	}
	companiesDto := make([]dto.CompanyDto, len(companies))
	for i, company := range companies {
		companiesDto[i] = dto.CompanyDTOFromModelWithID(company)
	}
	result.Results = companiesDto

	return result, err
}

func (s *CompanyService) GetCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (dto.PaginationResult[[]dto.CompanyDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var result dto.PaginationResult[[]dto.CompanyDto]
	if allow := s.permissionService.CheckPermission(ctx, req, domain.COMPANY, domain.VIEW); !allow {
		return result, domain.ACTION_NOT_ALLOWED
	}
	var companies []*model.Company
	queryBuilder := s.conn.Db.WithContext(ctx).Table("companies").
		Joins("JOIN user_relations ON user_relations.company_id = companies.id").
		Where("user_relations.user_id = ?", req.User.ID)

	err := queryBuilder.
		Count(&result.Total).Error

	if d.Query != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	}

	err = queryBuilder.
		Scopes(s.conn.Paginate(req.Params)).
		// Preload("CompanyPlugins").Preload("CompanyDepartments").
		Find(&companies).Error
	if err != nil {
		return result, err
	}

	companiesDto := make([]dto.CompanyDto, len(companies))
	for i, company := range companies {
		companiesDto[i] = dto.CompanyDTOFromModel(company)
	}
	result.Results = companiesDto
	return result, err
}

func (s *CompanyService) GetCompanyByUuid(ctx context.Context, uuid string) (model.Company, error) {
	var company model.Company
	err := s.conn.Db.WithContext(ctx).First(&company, "uuid = ?", uuid).Error
	if err != nil {
		return company, err
	}
	return company, err
}

func (s *CompanyService) GetCompanyDetail(req *common.RequestContext, i *dto.RequestEntity) (
	dto.EntityResponse[dto.CompanyDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var res dto.EntityResponse[dto.CompanyDto]
	var company model.Company
	err := s.conn.Db.WithContext(ctx).Where(&model.Company{UUID: i.ID}).First(&company).Error
	if err != nil {
		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyDetail"))
	}
	res.Body.Result = dto.CompanyDTOFromModel(&company)
	return res, err
}

// func (s *CompanyService) GetCompanyActions(req *common.RequestContext) []entity.Action {
// 	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
// 	defer cancel()
// 	var actions []entity.Action
// 	err := s.conn.Db.WithContext(ctx).Where(&entity.Action{EntityID: uint(model.COMPANY_ENTITY_ID)}).Find(&actions).Error
// 	if err != nil {
// 		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyActions"))
// 	}
// 	return actions
// }
