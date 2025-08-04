package a_company_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/config"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type AdminCompanyRepository interface {
	GetParentCompanies(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.CompanyDto], error)
	GetCompany(req *common.AdminRequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.CompanyDto], error)

	GetCompanyModules(req *common.AdminRequestContext, d *dto.RequestData) ([]dto.CompanyEntityDto, error)
	AddCompanyModules(req *common.AdminRequestContext, d *dto.AddCompanyModules) error
	GetCompanyUsers(req *common.AdminRequestContext, d *dto.RequestData) ([]dto.UserDto, error)
	AddCompanyUser(req *common.AdminRequestContext, d *dto.CreateUserAdminRequest) (model.UserRelation, error)
	CreateCompany(tx *query.QueryTx, req *common.AdminRequestContext, d *dto.CreateCompanyAdminRequest) (
		model.Company,model.CompanyDefault, error)
}

type adminCompanyRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	generator helpers.Generator
	pgConfig  config.PGConfig
}

func NewAdminCompanyRepository(
	db db.Connection,
	helpers *helpers.Helpers,
	pgConfig config.PGConfig,
) AdminCompanyRepository {
	return &adminCompanyRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		generator: helpers.Generator,
		pgConfig:  pgConfig,
	}
}

func (r *adminCompanyRepository) AddCompanyUser(req *common.AdminRequestContext, d *dto.CreateUserAdminRequest) (
	res model.UserRelation, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	res.User, err = r.createUser(req.Ctx, tx, d.Body.Identifier)
	if err != nil {
		return
	}
	res.Profile, err = r.createProfile(req.Ctx, tx, d)
	if err != nil {
		return
	}
	company, err := tx.Company.WithContext(req.Ctx).Where(
		tx.Company.ID.Eq(d.Body.CompanyID),
	).First()
	res.Company = *company
	role, err := tx.Role.WithContext(req.Ctx).Where(
		tx.Role.CompanyID.Eq(d.Body.CompanyID),
	).First()
	if err != nil {
		return
	}
	res.Role = *role
	err = tx.UserRelation.UnderlyingDB().WithContext(req.Ctx).Where(
		model.UserRelation{
			CompanyID: res.Company.ID,
			RoleID:    res.Role.ID,
			ProfileID: res.Profile.ID,
			UserID:    res.User.ID,
		},
	).FirstOrCreate(&res).Error
	if err != nil {
		return
	}
	err = r.createRoleActions(req.Ctx, tx, res.Role, d)
	err = tx.Commit()
	return
}

func (r *adminCompanyRepository) createRoleActions(ctx context.Context, tx *query.QueryTx, role model.Role,
	d *dto.CreateUserAdminRequest) (err error) {
	companyEntityQ := tx.CompanyEntity
	entities, err := companyEntityQ.WithContext(ctx).Where(
		companyEntityQ.CompanyID.Eq(d.Body.CompanyID),
	).Find()
	if err != nil {
		return
	}
	var roleActions []*model.RoleAction
	for _, entity := range entities {
		actions, err := tx.Action.WithContext(ctx).Where(
			tx.Action.EntityID.Eq(entity.EntityID),
		).Find()
		if err != nil {
			return err
		}
		for _, action := range actions {
			roleAction := &model.RoleAction{}
			roleAction.ActionID = action.ID
			roleAction.RoleID = role.ID
			roleActions = append(roleActions, roleAction)
		}
	}
	err = tx.RoleAction.WithContext(ctx).CreateInBatches(roleActions, len(roleActions))
	return
}

func (r *adminCompanyRepository) createProfile(ctx context.Context, tx *query.QueryTx, d *dto.CreateUserAdminRequest) (
	res model.Profile, err error) {
	partyID, err := tx.Profile.InsertParty(proto.PartyType_admin.String())
	if err != nil {
		return
	}
	res.ID = partyID
	res.EmailAddress = d.Body.Identifier
	res.GivenName = d.Body.FirstName
	res.FamilyName = d.Body.LastName
	err = tx.WithContext(ctx).Profile.Save(&res)
	return
}

func (r *adminCompanyRepository) createUser(ctx context.Context, tx *query.QueryTx, identifier string) (model.User, error) {
	userPassword := r.generator.GeneratePassword()
	pass := r.pgConfig.CryptoPass
	var (
		res model.User
		err error
	)
	user, err := tx.User.WithContext(ctx).Where(
		r.Q.User.Identifier.Eq(identifier),
	).First()
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	// .(`SELECT id FROM users WHERE identifier = ?`, user.Identifier).Scan(&existingUserID).Error
	if err != nil {
		fmt.Println("FAIL TO CHECK USER EXISTENCE", err)
		return res, err
	}
	if user != nil {
		return *user, err
	}
	partyID, err := tx.User.WithContext(ctx).InsertParty(proto.PartyType_admin.String())
	if err != nil {
		return res, err
	}
	err = tx.User.UnderlyingDB().WithContext(ctx).Raw(`
			INSERT INTO users (id,identifier, password_hash)
			VALUES (?,?, pgp_sym_encrypt(?, ?))
			RETURNING id
		`, partyID, identifier, userPassword, pass).Scan(&res.ID).Error
	if err != nil {
		return res, err
	}
	res.Identifier = identifier
	return res, err
}

func (r *adminCompanyRepository) GetCompanyUsers(req *common.AdminRequestContext, d *dto.RequestData) (
	res []dto.UserDto, err error) {
	companyID := r.convertor.StrtoInt(d.ID)
	userQ := r.Q.User
	userRelationQ := r.Q.UserRelation
	err = userRelationQ.WithContext(req.Ctx).Select(
		userQ.ID, userQ.UUID, userQ.Identifier, userQ.CreatedAt,
	).Where(
		userRelationQ.CompanyID.Eq(companyID),
	).Join(userQ, userQ.ID.EqCol(userRelationQ.UserID)).Limit(domain.DEFAULT_LIMIT).Scan(&res)
	return
}

func (r *adminCompanyRepository) AddCompanyModules(req *common.AdminRequestContext, d *dto.AddCompanyModules) (err error) {
	companyEntityQ := r.Q.CompanyEntity
	fmt.Println("MODULES", d.Body.Modules)
	for _, module := range d.Body.Modules {
		fmt.Println("MODULE", module)
		companyEntity := &model.CompanyEntity{}
		companyEntity.CompanyID = *module.CompanyID
		companyEntity.EntityID = module.EntityID
		companyEntity.Enabled = module.Enabled
		err = companyEntityQ.WithContext(req.Ctx).UnderlyingDB().Save(companyEntity).Error
		if err != nil {
			return
		}
	}

	return
}

func (r *adminCompanyRepository) GetCompanyModules(req *common.AdminRequestContext, d *dto.RequestData) (
	res []dto.CompanyEntityDto, err error) {
	companyID := r.convertor.StrtoInt(d.ID)
	entityQ := r.Q.Entity
	companyEntityQ := r.Q.CompanyEntity
	err = r.Q.Entity.WithContext(req.Ctx).Select(
		entityQ.ID.As("entity_id"), entityQ.Name.As("entity_name"), companyEntityQ.CompanyID,
		companyEntityQ.Enabled,
	).LeftJoin(companyEntityQ, companyEntityQ.EntityID.EqCol(entityQ.ID),
		companyEntityQ.CompanyID.Eq(companyID)).Scan(&res)
	return
}

func (r *adminCompanyRepository) GetCompany(req *common.AdminRequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.CompanyDto], err error) {
	companyQ := r.Q.Company
	r.Q.Company.WithContext(req.Ctx).
		Select(
			companyQ.Name, companyQ.UUID, companyQ.ID, companyQ.CreatedAt,
		).
		Where(
			companyQ.ID.Eq(r.convertor.StrtoInt(d.ID)),
		).Scan(&res.Entity)

	return
}

func (r *adminCompanyRepository) GetParentCompanies(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CompanyDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	companyQ := r.Q.Company
	builder := r.Q.WithContext(req.Ctx).Company

	//ADDING CONDITIONS
	conds = append(conds, companyQ.ParentID.IsNull())
	if d.Query != "" {
		conds = append(conds, companyQ.Name.Like("%"+d.Query+"%"))
	}

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.Company.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}

	builder = builder.Select(
		companyQ.ID, companyQ.UUID, companyQ.Name, companyQ.Code, companyQ.CreatedAt,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (r *adminCompanyRepository) CreateCompany(tx *query.QueryTx, req *common.AdminRequestContext, d *dto.CreateCompanyAdminRequest) (
	company model.Company,companyDefaults model.CompanyDefault, err error) {
	companyQ := tx.Company
	companyID, err := companyQ.WithContext(req.Ctx).InsertParty(proto.PartyType_company.String())
	if err != nil {
		return
	}
	company.ID = companyID

	company.Name = d.Body.Name
	err = companyQ.WithContext(req.Ctx).Save(&company)
	if err != nil {
		return
	}

	companyDefaults.Currency = string(common.CurrencyCodeBOB)
	companyDefaults.CompanyID = company.ID
	err = tx.CompanyDefault.WithContext(req.Ctx).Save(&companyDefaults)
	if err != nil {
		return
	}
	return 
}
