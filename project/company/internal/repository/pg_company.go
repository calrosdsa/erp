package company_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	// "erp/gen/proto"
	"erp/internal/app/entity"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	logger "erp/pkg/logger"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	CreateCompany(tx *query.QueryTx, req *common.RequestContext, d *dto.CreateCompanyRequest) (
		res model.Company, err error)
	GetAllUserCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CompanyDto], err error)
	// GetCompanyByUuid(ctx context.Context, uuid string) (res model.Company, err error)
	GetCompanyDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.CompanyDto], err error)
	GetAccountSetting(req *common.RequestContext, i *dto.RequestData) (
		res dto.AccountSettingsDto, err error)
	EditAccountSetting(req *common.RequestContext, d dto.AccountSettingData) (err error)
}

type companyRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	emitLog   logger.EmitLog
	generator helpers.Generator
}

func NewCompanyRepository(
	db db.Connection,
	helpers *helpers.Helpers,
	logger logger.Logger,
) *companyRepository {
	return &companyRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		emitLog:   logger.EmitLog("company-service"),
		generator: helpers.Generator,
	}
}

func (r *companyRepository) EditAccountSetting(req *common.RequestContext, d dto.AccountSettingData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.AccountSetting.UnderlyingDB().WithContext(req.Ctx).Where(
		tx.AccountSetting.CompanyID.Eq(req.ActiveCompany.ID),
	).Updates(data).Error
	if err != nil {
		return
	}
	// err = tx.AccountSetting.InsertActivity(req.ActiveCompany.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	// if err != nil {
	// 	return
	// }
	err = tx.Commit()
	return
}

func (r *companyRepository) GetAccountSetting(req *common.RequestContext, i *dto.RequestData) (
	res dto.AccountSettingsDto, err error) {
	companyID := r.convertor.StrtoInt(i.ID)
	cashAcctQ := r.Q.Ledger.As("cash")
	bankAcctQ := r.Q.Ledger.As("bank")
	payableAcctQ := r.Q.Ledger.As("payable")
	cOfGoodsSoldAcctQ := r.Q.Ledger.As("cost_of_goods")
	receivableAcctQ := r.Q.Ledger.As("receivable")
	incomeAcctQ := r.Q.Ledger.As("income")
	e := r.Q.AccountSetting

	err = e.WithContext(req.Ctx).Select(
		cashAcctQ.ID.As("cash_acct_id"), cashAcctQ.Name.As("cash_acct"),
		bankAcctQ.ID.As("bank_acct_id"), bankAcctQ.Name.As("bank_acct"),
		payableAcctQ.ID.As("payable_acct_id"), payableAcctQ.Name.As("payable_acct"),
		cOfGoodsSoldAcctQ.ID.As("cost_of_goods_sold_acct_id"), cOfGoodsSoldAcctQ.Name.As("cost_of_goods_sold_acct"),
		receivableAcctQ.ID.As("receivable_acct_id"),
		receivableAcctQ.Name.As("receivable_acct"),
		incomeAcctQ.ID.As("income_acct_id"), incomeAcctQ.Name.As("income_acct"),
	).Where(e.CompanyID.Eq(companyID)).
		Join(cashAcctQ, cashAcctQ.ID.EqCol(e.CashAccunt)).
		Join(bankAcctQ, bankAcctQ.ID.EqCol(e.BankAccount)).
		Join(payableAcctQ, payableAcctQ.ID.EqCol(e.PayableAccount)).
		Join(cOfGoodsSoldAcctQ, cOfGoodsSoldAcctQ.ID.EqCol(e.CostOfGoodSoldAccount)).
		Join(receivableAcctQ, receivableAcctQ.ID.EqCol(e.ReceivableAccount)).
		Join(incomeAcctQ, incomeAcctQ.ID.EqCol(e.IncomeAccount)).Scan(&res)
	if err != nil {
		return
	}
	return
}

func (s *companyRepository) CreateCompany(tx *query.QueryTx, req *common.RequestContext, d *dto.CreateCompanyRequest) (
	res model.Company, err error) {

	res.Name = d.Body.Name
	if d.Body.ParentID != nil {
		res.ParentID = d.Body.ParentID
		res.Ordinal = int32(entity.DEPTH_TWO)
	} else {
		res.ParentID = &req.ActiveCompany.ID
	}

	// code := s.generateCompanyCode(tx)
	res.Code = s.generateCompanyCode(tx.Company.UnderlyingDB())
	err = tx.Company.Save(&res)
	if err != nil {
		return
	}
	// err = s.addCompanyDepartment(ctx, tx, *res.ParentID, uint(res.ID))
	// if err != nil {
	// 	return err
	// }
	err = s.insertUserCompany(req, tx, res.ID)
	if err != nil {
		return
	}
	return
}
func (s *companyRepository) insertUserCompany(req *common.RequestContext, tx *query.QueryTx, companyId int64) (err error) {
	err = tx.Company.UnderlyingDB().WithContext(req.Ctx).Exec(`insert into user_relations(company_id,user_id,role_id,profile_id) values($1,$2,$3,$4)
	ON CONFLICT (company_id, user_id,role_id,profile_id) DO NOTHING`, companyId, req.User.ID, req.Role.ID, req.Profile.ID).Error
	return
}

func (s *companyRepository) generateCompanyCode(db *gorm.DB) (code string) {
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

func (s *companyRepository) GetValidaParentCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CompanyDto], err error) {
	return
}

func (r *companyRepository) GetAllUserCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CompanyDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	companyQ := r.Q.Company
	builder := r.Q.WithContext(req.Ctx).Company

	//ADDING CONDITIONS
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
	}
	userRelationQ := r.Q.UserRelation
	builder = builder.Select(
		companyQ.ID, companyQ.UUID, companyQ.Name, companyQ.Code, companyQ.CreatedAt,
	).
		Join(userRelationQ, userRelationQ.CompanyID.EqCol(companyQ.ID),
			userRelationQ.UserID.Eq(req.User.ID)).
		Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return

	// var result dto.PaginationResult[[]dto.CompanyDto]
	// var companies []*model.Company
	// queryBuilder := s.Q.Company.UnderlyingDB().WithContext(req.Ctx).Table("companies").
	// 	Joins("JOIN user_relations ON user_relations.company_id = companies.id").
	// 	Where("user_relations.user_id = ?", req.User.ID)

	// err := queryBuilder.
	// 	Count(&result.Total).Error

	// if d.Query != "" {
	// 	queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	// }

	// err = queryBuilder.
	// 	Scopes(s.conn.Paginate(req.Params)).
	// 	Find(&companies).Error
	// if err != nil {
	// 	return result, err
	// }
	// companiesDto := make([]dto.CompanyDto, len(companies))
	// for i, company := range companies {
	// 	companiesDto[i] = dto.CompanyDTOFromModelWithID(company)
	// }
	// result.Results = companiesDto

	// return result, err
}

// func (r *companyRepository) GetCompanyByUuid(ctx context.Context, uuid string) (res model.Company, err error) {
// 	company, err := r.Q.Company.WithContext(ctx).Where(
// 		r.Q.Company.UUID.Eq(uuid),
// 	).First()
// 	if err != nil {
// 		return
// 	}
// 	return *company, err
// }

func (s *companyRepository) GetCompanyDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.CompanyDto], err error) {
	id := s.convertor.StrtoInt(i.ID)
	companyQ := s.Q.Company
	err = s.Q.Company.WithContext(req.Ctx).
		Select(companyQ.ID, companyQ.UUID, companyQ.Name, companyQ.CreatedAt).
		Where(companyQ.ID.Eq(id)).
		Scan(&res.Entity)
	return res, err
}
