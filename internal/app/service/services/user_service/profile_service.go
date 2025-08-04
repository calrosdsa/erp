package userservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/connection"
	"erp/internal/domain"
	"erp/internal/app/event-bus/event"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"fmt"
	"strings"
	"time"
)

type ProfileService struct {
	conn              *connection.Connection
	timeout           time.Duration
	emitLog           logger.EmitLog
	permissionService permission.PermissionService
	userService       *UserService
	eventHelper       *helpers.EventHelper
}

func NewProfileService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	userService *UserService,
	logger logger.Logger,
) *ProfileService {
	return &ProfileService{
		conn:              conn,
		timeout:           timeout,
		permissionService: permissionService,
		userService:       userService,
		eventHelper:       helpers.Event,
		emitLog: logger.EmitLog("profile-service"),
	}
}

func (s *ProfileService) CreateUser(req *common.RequestContext, i *dto.CreateUserRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var err error
	tx := s.conn.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateUser"))
		}
	}()
	allow := s.permissionService.CheckPermission(ctx, req, domain.USER, domain.CREATE)
	if !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	var (
		user         model.User
		userRelation model.UserRelation
	)
	user,err = s.userService.InsertUser(ctx, tx, i.Body.Email)
	if err != nil {
		return err
	}
	role,err := s.conn.Q.Role.WithContext(ctx).Where(
		s.conn.Q.Role.UUID.Eq(i.Body.RoleUUID),
	).First()
	userRelation.RoleID = role.ID
	userRelation.UserID = int64(user.ID)
	userRelation.User = user
	profile, err := s.createProfile(tx, i)
	if err != nil {
		return err
	}
	userRelation.ProfileID = int64(profile.ID)
	userRelation.Profile = profile
	// if req.Profile.PartyTypeCode == domain.PARTY_ADMIN {
	// 	for _, companyId := range i.Body.CompanyIds {
	// 		fmt.Println("COMPANY ID",companyId)
	// 		if allow := s.permissionService.CheckIfUserIsCompanyAdmin(ctx, req, int64(companyId), domain.CREATE_USER); !allow {
	// 			fmt.Println("NO ALLOW CUSTOM PERMISSION")
	// 			return domain.ACTION_NOT_ALLOWED
	// 		}
	// 		userRelation.CompanyID = int64(companyId)
	// 		err = s.createUserRelation(ctx,tx, &userRelation)
	// 		if err != nil && err != domain.USER_ALREDY_EXIST {
	// 			return err
	// 		}
	// 	}
	// } else {
	// 	userRelation.CompanyID = req.ActiveCompany.ID
	// 	userRelation.Company = req.ActiveCompany
	// 	if allow := s.permissionService.CheckIfUserIsCompanyAdmin(ctx, req, req.ActiveCompany.ID, domain.CREATE_USER); !allow {
	// 		return domain.ACTION_NOT_ALLOWED
	// 	}
	// 	err = s.createUserRelation(ctx,tx, &userRelation)
	// 	if err != nil && err != domain.USER_ALREDY_EXIST {
	// 		return err
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		return err
	}
	s.eventHelper.Publish(event.NOTIFICATION_EVENT, event.NotificationData{
		NotificationEventType: event.NOTIFY_NEW_CLIENT,
		Data: event.NotificationPayload{
			Payload:        userRelation,
			RequestContext: *req,
		},
	},
	)
	return nil
}

func (s *ProfileService) createUserRelation(ctx context.Context,tx *query.QueryTx, userRelation *model.UserRelation) error {
	_,err := tx.UserRelation.WithContext(ctx).Where(
		s.conn.Q.UserRelation.RoleID.Eq(userRelation.RoleID),
		s.conn.Q.UserRelation.CompanyID.Eq(userRelation.CompanyID),
		s.conn.Q.UserRelation.ProfileID.Eq(userRelation.ProfileID),
		s.conn.Q.UserRelation.UserID.Eq(userRelation.UserID),
	).First()
	if err == nil {
		return domain.USER_ALREDY_EXIST
	}

	err = tx.UserRelation.WithContext(ctx).Create(&model.UserRelation{
		UserID: userRelation.UserID,
		ProfileID: userRelation.ProfileID,
		CompanyID: userRelation.CompanyID,
		RoleID: userRelation.RoleID,
	})
	if err != nil {
		return err
	}

	return err
}

func (s *ProfileService) createProfile(tx *query.QueryTx, i *dto.CreateUserRequest) (model.Profile, error) {
	var (
		profile model.Profile
	)

	partyID,err := tx.Profile.InsertParty(i.Body.PartyCode)
	if err != nil {
		return profile, err
	}
	profile.GivenName = i.Body.GivenName
	profile.FamilyName = i.Body.FamilyName
	profile.EmailAddress = i.Body.Email
	profile.ID = partyID
	profile.PhoneNumber = &i.Body.PhoneNumber
	err = tx.Profile.Save(&profile)
	if err != nil {
		return profile, err
	}
	keyValues := make([]*model.KeyValue, len(i.Body.KeyValueData))
	if len(keyValues) == 0 {
		return profile, nil
	}
	for i, keyValue := range i.Body.KeyValueData {
		kV := &model.KeyValue{}
		kV.Key = keyValue.Key
		kV.Value = keyValue.Value
		kV.PartyID = profile.ID
		keyValues[i] = kV
	}
	err = tx.KeyValue.CreateInBatches(keyValues, len(keyValues))
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (s *ProfileService) GetUserProfileDetail(req *common.RequestContext, i *dto.RequestEntity) (
	dto.ResultEntity[dto.ProfileDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.ResultEntity[dto.ProfileDto]
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserProfileDetail"))
		}
	}()
	allow := s.permissionService.CheckPermission(ctx, req, domain.USER, domain.VIEW)
	if !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	profileUuid := i.ID
	profile, err := s.conn.Q.Profile.WithContext(ctx).Where(
		s.conn.Q.Profile.UUID.Eq(profileUuid),
	).First()
	res.Entity = dto.ProfileDTOFromModel(profile)
	return res, err
}

func (s *ProfileService) GetProfileSession(req *common.RequestContext) (
	dto.ResultEntity[dto.ProfileDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.ResultEntity[dto.ProfileDto]
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetProfileSession"))
		}
	}()
	profile, err := s.conn.Q.Profile.WithContext(ctx).Where(
		s.conn.Q.Profile.ID.Eq(req.Profile.ID),
	).First()
	if err != nil {
		return res, err
	}
	res.Entity = dto.ProfileDTOFromModel(profile)
	return res, err
}

func (s *ProfileService) UpdateProfileSession(req *common.RequestContext,i *dto.UpdateProfileRequest) (error){
	ctx,cancel := context.WithTimeout(req.Ctx,s.timeout)
	defer cancel()
	var (
		err error
	)
	defer func ()  {
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("UpdateProfileSession"))
		}
	}()
	err = s.updateProfile(ctx,req.Profile.ID,&i.Body.ProfileFields)
	return err
}

func (s *ProfileService) updateProfile(ctx context.Context,profileID int64,i *dto.EditableProfileFields)(error){
	_,err := s.conn.Q.Profile.WithContext(ctx).Where(
		s.conn.Q.Profile.ID.Eq(profileID),
	).Updates(model.Profile{
		GivenName: i.GivenName,
		FamilyName: i.FamilyName,
		PhoneNumber: i.PhoneNumber,
	})
	return err
}

func (s *ProfileService) GetCompanyUserProfiles(req *common.RequestContext, d *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.ProfileL], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var result dto.PaginationResult[[]dto.ProfileL]
	allow := s.permissionService.CheckPermission(ctx, req, domain.USER, domain.VIEW)
	if !allow {
		return result, domain.ACTION_NOT_ALLOWED
	}

	var query strings.Builder
	query.WriteString(fmt.Sprintf(`
	      select p.id,p.uuid,p.given_name,p.family_name,p.email_address,p.phone_number,
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

	s.conn.Db.Raw(query.String()).Scan(&result.Results)

	return result, nil

}

