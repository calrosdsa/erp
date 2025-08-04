package core

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"erp/pkg/logger"
	"fmt"
)

type coreService struct {
	conn    db.Connection
	Q       *query.Query
	emitLog logger.EmitLog
	locale  helpers.Locale
}

func NewCoreService(
	conn db.Connection,
	logger logger.Logger,
	helpers *helpers.Helpers,
) repository.CoreService {
	return &coreService{
		conn:    conn,
		Q:       conn.GetQ(),
		emitLog: logger.EmitLog("core-service"),
		locale:  helpers.Locale,
	}
}
func (r *coreService) GetPartyIDByType(ctx context.Context, partyID int64, partyType string) (res *int64, err error) {
	partyRefQ := r.Q.PartyReference
	partyQ := r.Q.Party
	party, err := partyQ.WithContext(ctx).
		Join(partyRefQ, partyRefQ.PartyID.Eq(partyID)).
		Where(
			partyQ.PartyTypeCode.Eq(partyType),
		).First()
	if err != nil {
		return
	}
	res = &party.ID
	return
}

func (r *coreService) GetPartyAccountingDimension(req *common.RequestContext, partyID int64) (
	res dto.AccountingDimensionDto, err error) {
	refQ := r.Q.PartyReference
	partyRefQ := r.Q.Party
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	costCenter := dto.AccountingDimensionDto{}
	err = refQ.WithContext(req.Ctx).
		Select(
			costCenterQ.Name.As("cost_center"), costCenterQ.UUID.As("cost_center_uuid"), costCenterQ.ID.As("cost_center_id"),
		).
		Join(partyRefQ, partyRefQ.ID.EqCol(refQ.ReferenceID)).
		LeftJoin(costCenterQ, refQ.ReferenceID.EqCol(costCenterQ.ID)).
		Where(
			refQ.PartyID.Eq(partyID),
			partyRefQ.PartyTypeCode.Eq(proto.PartyType_costCenter.String()),
		).Scan(&costCenter)
	if err != nil {
		return
	}
	project := dto.AccountingDimensionDto{}
	err = refQ.WithContext(req.Ctx).
		Select(
			projectQ.Name.As("project"), projectQ.UUID.As("project_uuid"), projectQ.ID.As("project_id"),
		).
		Join(partyRefQ, partyRefQ.ID.EqCol(refQ.ReferenceID)).
		LeftJoin(projectQ, refQ.ReferenceID.EqCol(projectQ.ID)).
		Where(
			refQ.PartyID.Eq(partyID),
			partyRefQ.PartyTypeCode.Eq(proto.PartyType_project.String()),
		).Scan(&project)
	res.CostCenter = costCenter.CostCenter
	res.CostCenterID = costCenter.CostCenterID
	res.CostCenterUUID = costCenter.CostCenterUUID
	res.ProjectID = project.ProjectID
	res.Project = project.Project
	res.ProjectUUID = project.ProjectUUID
	return
}

func (r *coreService) GetAddress(req *common.RequestContext, id int64) (res dto.AddressDto, err error) {
	addressQ := r.Q.Address
	err = r.Q.Address.WithContext(req.Ctx).Select(
		addressQ.UUID, addressQ.Title, addressQ.City, addressQ.StreetLine1, addressQ.StreetLine2,
		addressQ.Province, addressQ.PostalCode, addressQ.PhoneNumber, addressQ.PhoneNumber, addressQ.Company,
	).Where(
		addressQ.ID.Eq(id),
	).Scan(&res)
	if err != nil {
		r.emitLog.Err(err, logger.OptionsLog.WithMethod("GerPartyAddresses"))
	}
	return
}

func (r *coreService) GetContact(req *common.RequestContext, id int64) (res dto.ContactDto, err error) {
	contactQ := r.Q.Contact
	err = r.Q.Contact.WithContext(req.Ctx).Select(
		contactQ.UUID, contactQ.Name, contactQ.Email, contactQ.PhoneNumber,
		contactQ.Gender,
	).Where(
		contactQ.ID.Eq(id),
	).Scan(&res)
	if err != nil {
		r.emitLog.Err(err, logger.OptionsLog.WithMethod("GerPartyAddresses"))
	}
	return
}

func (r *coreService) GetPartyAddresses(req *common.RequestContext, partyID int64) []dto.AddressDto {
	var addresses []dto.AddressDto
	partyReferenceQ := r.Q.PartyReference
	addressQ := r.Q.Address
	err := r.Q.PartyReference.WithContext(req.Ctx).Select(
		addressQ.UUID, addressQ.Title, addressQ.City, addressQ.StreetLine1, addressQ.StreetLine2,
	).Join(addressQ, partyReferenceQ.PartyID.EqCol(addressQ.ID)).Where(
		partyReferenceQ.ReferenceID.Eq(partyID),
	).Scan(&addresses)
	if err != nil {
		r.emitLog.Err(err, logger.OptionsLog.WithMethod("GerPartyAddresses"))
	}
	return addresses
}

func (r *coreService) GetPartyContacts(req *common.RequestContext, partyID int64) []dto.ContactDto {
	var contacts []dto.ContactDto
	partyReferenceQ := r.Q.PartyReference
	contactQ := r.Q.Contact
	err := r.Q.PartyReference.WithContext(req.Ctx).Select(
		contactQ.ID, contactQ.UUID, contactQ.Name, contactQ.Email, contactQ.PhoneNumber,
		contactQ.Gender,
	).Join(contactQ, partyReferenceQ.PartyID.EqCol(contactQ.ID)).Where(
		partyReferenceQ.ReferenceID.Eq(partyID),
		contactQ.CompanyID.Eq(req.ActiveCompany.ID),
	).Scan(&contacts)
	if err != nil {
		r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPartyContacts"))
	}
	return contacts
}

func (s *coreService) GerActivitiesByPartyID(req *common.RequestContext, partyID int64) []dto.ActivityDto {
	var err error
	res := []dto.ActivityDto{}
	mentions := []dto.ActivityMentionDto{}
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GerActivitiesByPartyID"))
		}
	}()
	activityQ := s.Q.Activity
	profileQ := s.Q.Profile
	aDeadlineQ := s.Q.ActivityDeadline
	aCommentQ := s.Q.ActivityComment
	err = s.Q.Activity.WithContext(req.Ctx).Select(
		activityQ.ProfileID, activityQ.ID, activityQ.IsPinned, activityQ.CreatedAt, activityQ.Type,
		activityQ.Data,
		profileQ.GivenName.As("profile_given_name"), profileQ.FamilyName.As("profile_family_name"),
		profileQ.Avatar.As("profile_avatar"),
		aDeadlineQ.Content, aDeadlineQ.Title, aDeadlineQ.Color, aDeadlineQ.Deadline,
		aDeadlineQ.Content, aDeadlineQ.IsCompleted, aDeadlineQ.Address, aDeadlineQ.Link,
		aCommentQ.Comment,
	).
		Join(profileQ, profileQ.ID.EqCol(activityQ.ProfileID)).
		LeftJoin(
			aDeadlineQ, aDeadlineQ.ActivityID.EqCol(activityQ.ID),
		).
		LeftJoin(
			aCommentQ, aCommentQ.ActivityID.EqCol(activityQ.ID),
		).
		Where(
			activityQ.PartyID.Eq(partyID),
		).Order(activityQ.CreatedAt.Desc()).Limit(domain.DEFAULT_ACTIVITY_LIMIT).Scan(&res)
	if err != nil {
		return res
	}
	ids := make([]int32, len(res))
	for i, activity := range res {
		ids[i] = activity.ID
	}
	err = s.Q.ActivityMention.WithContext(req.Ctx).Select(
		s.Q.ActivityMention.ID, s.Q.ActivityMention.ActivityID, s.Q.ActivityMention.StartIndex, s.Q.ActivityMention.EndIndex,
		s.Q.Profile.GivenName, s.Q.Profile.GivenName, s.Q.Profile.UUID.As("profile_uuid"),
	).Join(
		s.Q.Profile, s.Q.Profile.ID.EqCol(s.Q.ActivityMention.ProfileID),
	).Where(
		s.Q.ActivityMention.ActivityID.In(ids...),
	).Scan(&mentions)
	// Mapea las menciones por ActivityID
	mentionMap := make(map[int32][]dto.ActivityMentionDto)
	for _, mention := range mentions {
		mentionMap[mention.ActivityID] = append(mentionMap[mention.ActivityID], mention)
	}
	fmt.Println("MENTIONS", mentions, "Acitivity ids", ids)

	// Asigna las menciones a cada actividad correspondiente
	for i := range res {
		activity := res[i]
		res[i].Mentions = mentionMap[activity.ID]
		// activity.Mentions = mentionMap[activity.ID]
	}

	for _, activity := range res {
		switch proto.ActivityType_value[activity.Type] {
		// case int32(proto.ActivityType_CREATE):
		// 	res[i].Comment = s.locale.MustLocalize(
		// 		helpers.OptionsLocale.WithID("Activity.Created"),
		// 		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		// 		helpers.OptionsLocale.WithTemplate(map[string]string{
		// 			"User": fmt.Sprintf("%s %s", activity.ProfileGivenName, activity.ProfileFamilyName),
		// 		}),
		// 	)
		// case int32(proto.ActivityType_EDIT):
		// 	res[i].Comment = s.locale.MustLocalize(
		// 		helpers.OptionsLocale.WithID("Activity.Edited"),
		// 		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		// 		helpers.OptionsLocale.WithTemplate(map[string]string{
		// 			"User": fmt.Sprintf("%s %s", activity.ProfileGivenName, activity.ProfileFamilyName),
		// 		}),
		// 	)
		// case int32(proto.ActivityType_EDIT_PAID_AMOUT):
		// 	res[i].Comment = s.locale.MustLocalize(
		// 		helpers.OptionsLocale.WithID("Activity.EditPaidAmount"),
		// 		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		// 		helpers.OptionsLocale.WithTemplate(map[string]string{
		// 			"User": fmt.Sprintf("%s %s", activity.ProfileGivenName, activity.ProfileFamilyName),
		// 		}),
		// 	)
		// case int32(proto.ActivityType_UPDATE_STATUS):
		// 	res[i].Comment = s.locale.MustLocalize(
		// 		helpers.OptionsLocale.WithID("Activity.UpdateStatus"),
		// 		helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		// 		helpers.OptionsLocale.WithTemplate(map[string]string{
		// 			"User":   fmt.Sprintf("%s %s", activity.ProfileGivenName, activity.ProfileFamilyName),
		// 			"Status": *activity.Arg1,
		// 		}),
		// 	)
		}
	}

	return res
}
