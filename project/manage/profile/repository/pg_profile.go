package profile_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
)

type ProfileRepository interface {
	UpdateProfileSession(req *common.RequestContext, i *dto.UpdateProfileRequest) error

	GetUserProfileDetail(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ProfileDto, error)
	GetProfileSession(req *common.RequestContext) (
		res dto.ProfileDto, err error)
	GetProfiles(req *common.RequestContext, d dto.ProfilesRequest) (res []dto.ProfileDto, err error)
	GetCompanyUserProfiles(req *common.RequestContext, d *dto.RequestPaginationData) (
		res []dto.ProfileDto, err error)
}

type profileRepo struct {
	Q         *query.Query
	query     helpers.QueryHelper
	convertor helpers.ConvertorHelper
	session repository.SessionService
}

func NewProfileRepo(
	db db.Connection,
	helpers *helpers.Helpers,
	session repository.SessionService,
) ProfileRepository {
	return &profileRepo{
		Q:         db.GetQ(),
		query:     helpers.Query,
		convertor: helpers.Convertor,
		session: session,
	}
}


func (s *profileRepo) GetUserProfileDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ProfileDto, err error) {
	
	profileID :=s.convertor.StrtoInt(i.ID)
	err = s.Q.Profile.WithContext(req.Ctx).Where(
		s.Q.Profile.ID.Eq(profileID),
	).Scan(&res)

	return res, err
}

func (s *profileRepo) GetProfileSession(req *common.RequestContext) (
	res dto.ProfileDto, err error) {
	err = s.Q.Profile.WithContext(req.Ctx).Where(
		s.Q.Profile.ID.Eq(req.Profile.ID),
	).Scan(&res)
	if err != nil {
		return res, err
	}
	return
}

func (s *profileRepo) UpdateProfileSession(req *common.RequestContext, i *dto.UpdateProfileRequest) (err error) {

	err = s.updateProfile(req.Ctx, req.Profile.ID, &i.Body.ProfileFields)
	return err
}

func (s *profileRepo) updateProfile(ctx context.Context, profileID int64, i *dto.EditableProfileFields) error {
	_, err := s.Q.Profile.WithContext(ctx).Where(
		s.Q.Profile.ID.Eq(profileID),
	).Updates(model.Profile{
		GivenName:   i.GivenName,
		FamilyName:  i.FamilyName,
		PhoneNumber: i.PhoneNumber,
	})
	return err
}

func (s *profileRepo) GetCompanyUserProfiles(req *common.RequestContext, d *dto.RequestPaginationData) (
	res []dto.ProfileDto, err error) {
	var query strings.Builder
	query.WriteString(fmt.Sprintf(`
		select p.id,p.uuid,p.given_name,p.family_name,
		CONCAT(given_name, ' ', family_name) as full_name,
		p.email_address,p.phone_number,
		pt.code as party_code,pt.name as party_name
		from user_relations as ur
		inner join profiles as p on p.id = ur.profile_id
		inner join parties as parties on parties.id = p.id
		inner join party_types as pt on pt.code = parties.party_type_code
		where ur.company_id = %d
	`, req.ActiveCompany.ID))

	if d.Query != "" {
		query.WriteString(fmt.Sprintf(`CONCAT(given_name, ' ', family_name) ILIKE %s`, "%"+d.Query+"%"))
	}
	err = s.Q.Profile.UnderlyingDB().Raw(query.String()).Scan(&res).Error
	return
}

func (r *profileRepo) GetProfiles(req *common.RequestContext, d dto.ProfilesRequest) (res []dto.ProfileDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Deal
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.profilesQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *profileRepo) profilesQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`select e.id,e.uuid,e.given_name,e.family_name,
		CONCAT(e.given_name, ' ', e.family_name) as full_name,
		e.email_address,e.phone_number
		from profiles as e
		`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}
