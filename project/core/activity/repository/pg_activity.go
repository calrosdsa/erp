package activity_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"
	"fmt"
)

type ActivityRepository interface {
	CreateActivity(tx *query.QueryTx, req *common.RequestContext, i dto.ActivityData) (
		res dto.ActivityDto, err error)
	EditActivity(tx *query.QueryTx, req *common.RequestContext, i dto.ActivityData) error
	DeleteActivity(tx *query.QueryTx, req *common.RequestContext, i dto.DeleteRequest) error
	CreateActivityStatus(tx *query.QueryTx, from string, to string, activity *model.Activity) (err error)
}

type activityRepository struct {
	conn      db.Connection
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewActivityRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) ActivityRepository {
	return &activityRepository{
		conn:      conn,
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}
func (r *activityRepository) DeleteActivity(tx *query.QueryTx, req *common.RequestContext, i dto.DeleteRequest) (err error) {
	id := r.convertor.StrtoInt(i.ID)
	if i.Type == proto.ActivityType_ACTIVITY.String() {
		err = r.deleteActivityDeadline(tx, req.Ctx, int32(id))
	}
	if i.Type == proto.ActivityType_COMMENT.String() {
		err = r.deleteActivityComment(tx, req.Ctx, int32(id))
	}
	return err
}

func (r *activityRepository) EditActivity(tx *query.QueryTx, req *common.RequestContext, i dto.ActivityData) (err error) {
	if i.Type == proto.ActivityType_ACTIVITY.String() {
		err = r.editActivityDeadline(tx, req.Ctx, i.ActivityDeadLine)
	}
	if i.Type == proto.ActivityType_COMMENT.String() {
		err = r.editActivityComment(tx, req.Ctx, i.ActivityComment)
	}
	return err
}

func (r *activityRepository) CreateActivity(tx *query.QueryTx, req *common.RequestContext, i dto.ActivityData) (
	res dto.ActivityDto, err error) {

	activity := model.Activity{
		PartyID:   i.PartyID,
		ProfileID: req.Profile.ID,
		Type:      i.Type,
	}

	fmt.Println("INSERTING ACTIVITY", activity)
	if err = tx.Activity.Save(&activity); err != nil {
		return
	}

	// Map base activity fields
	res = dto.ActivityDto{
		ID:        activity.ID,
		Type:      activity.Type,
		CreatedAt: activity.CreatedAt,
		IsPinned:  activity.IsPinned,
		Data:      activity.Data,
		ProfileID: activity.ProfileID,
	}

	// Handle activity type-specific logic
	if i.Type == proto.ActivityType_ACTIVITY.String() {
		i.ActivityDeadLine.Fields.ActivityID = activity.ID
		var activityDeadline model.ActivityDeadline
		if activityDeadline, err = r.createActivityDeadline(tx, req.Ctx, i.ActivityDeadLine); err != nil {
			return res, err
		}

		res.Link = activityDeadline.Link
		res.PartyID = activityDeadline.PartyID
		res.Deadline = &activityDeadline.Deadline
		res.Address = activityDeadline.Address
		res.Title = activityDeadline.Title
		res.Content = activityDeadline.Content
		res.Color = &activityDeadline.Color
		res.IsCompleted = &activityDeadline.IsCompleted
	}

	if i.Type == proto.ActivityType_COMMENT.String() {
		i.ActivityComment.Fields.ActivityID = activity.ID
		var activityComment model.ActivityComment
		if activityComment, err = r.createActivityComment(tx, req.Ctx, i.ActivityComment, activity.ID); err != nil {
			return res, err
		}
		res.Comment = &activityComment.Comment
	}

	// Fetch profile details
	profileQ := tx.Profile
	profile, err := profileQ.WithContext(req.Ctx).Select(
		profileQ.GivenName,
		profileQ.FamilyName,
		profileQ.Avatar,
	).Where(
		profileQ.ID.Eq(activity.ProfileID),
	).First()
	if err != nil {
		return
	}
	res.ProfileFamilyName = profile.FamilyName
	res.ProfileGivenName = profile.GivenName
	res.ProfileAvatar = profile.Avatar

	// Fetch mentions
	var mentions []dto.ActivityMentionDto
	if err = tx.ActivityMention.WithContext(req.Ctx).Select(
		tx.ActivityMention.ID,
		tx.ActivityMention.ActivityID,
		tx.ActivityMention.StartIndex,
		tx.ActivityMention.EndIndex,
		tx.Profile.GivenName,
		tx.Profile.FamilyName, // Fixed from duplicate GivenName
		tx.Profile.UUID.As("profile_uuid"),
	).Join(tx.Profile,
		tx.Profile.ID.EqCol(tx.ActivityMention.ProfileID),
	).Where(
		tx.ActivityMention.ActivityID.Eq(activity.ID), // Changed from In() to Eq()
	).Scan(&mentions); err != nil {
		return
	}
	res.Mentions = mentions
	return
}

func (r *activityRepository) createActivityComment(tx *query.QueryTx, ctx context.Context,
	d dto.ActivityCommentData, activityID int32) (res model.ActivityComment, err error) {
	fields := d.Fields
	fields.ActivityID = activityID
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	err = tx.WithContext(ctx).ActivityComment.Save(&res)
	if err != nil {
		return
	}
	err = r.processMentions(tx, ctx, d.Mentions, activityID)
	return
}

func (r *activityRepository) processMentions(tx *query.QueryTx, ctx context.Context, mentions []dto.ActivityMentionData,
	activityID int32) (err error) {
	for _, mention := range mentions {
		switch mention.Action {
		case string(domain.CREATE):
			err = r.saveMention(tx, ctx, mention, activityID)
		case string(domain.DELETE):
			err = r.deleteMention(tx, ctx, mention)
		}
		if err != nil {
			return
		}
	}
	return
}

func (r *activityRepository) deleteMention(tx *query.QueryTx, ctx context.Context, mention dto.ActivityMentionData) (err error) {
	_, err = tx.Mention.WithContext(ctx).Where(
		tx.ActivityMention.ID.Eq(mention.ID),
	).Delete()
	return
}

func (r *activityRepository) saveMention(tx *query.QueryTx, ctx context.Context, mention dto.ActivityMentionData,
	activityID int32) (err error) {
	var res model.ActivityMention
	fields := mention.Fields
	fields.ActivityID = activityID
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	// Update indexes before saving the deal to prevent modifying indexes for newly created entries.
	err = tx.WithContext(ctx).ActivityMention.Save(&res)
	if err != nil {
		return
	}
	return
}

func (r *activityRepository) editActivityComment(tx *query.QueryTx, ctx context.Context,
	d dto.ActivityCommentData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.ActivityComment.UnderlyingDB().WithContext(ctx).Model(
		&model.ActivityComment{ActivityID: d.Fields.ActivityID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = r.processMentions(tx, ctx, d.Mentions, d.Fields.ActivityID)
	return
}

func (r *activityRepository) deleteActivityComment(tx *query.QueryTx, ctx context.Context,
	id int32) (err error) {
	_, err = tx.WithContext(ctx).ActivityComment.Where(
		tx.ActivityComment.ActivityID.Eq(id),
	).Delete()
	if err != nil {
		return
	}
	_, err = tx.Mention.WithContext(ctx).Where(
		tx.ActivityMention.ActivityID.Eq(id),
	).Delete()
	return
}

func (r *activityRepository) deleteActivityDeadline(tx *query.QueryTx, ctx context.Context,
	id int32) (err error) {
	_, err = tx.WithContext(ctx).ActivityDeadline.Where(
		tx.ActivityDeadline.ActivityID.Eq(id),
	).Delete()
	return
}
func (r *activityRepository) editActivityDeadline(tx *query.QueryTx, ctx context.Context,
	d dto.ActivityDeadlineData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.ActivityDeadline.UnderlyingDB().WithContext(ctx).Model(
		&model.ActivityDeadline{ActivityID: d.Fields.ActivityID},
	).Updates(data).Error
	if err != nil {
		return
	}
	return
}

func (r *activityRepository) createActivityDeadline(tx *query.QueryTx, ctx context.Context,
	d dto.ActivityDeadlineData) (res model.ActivityDeadline, err error) {
	fields := d.Fields
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	// Update indexes before saving the deal to prevent modifying indexes for newly created entries.
	err = tx.WithContext(ctx).ActivityDeadline.Save(&res)
	if err != nil {
		return
	}
	return
}

func (r *activityRepository) CreateActivityStatus(tx *query.QueryTx, from string, to string, activity *model.Activity) (err error) {
	if from == to {
		return 
	}
	activityStageData := map[string]interface{}{
		"source":      from,
		"destination": to,
	}
	jsonString := domain.JsonStringify(activityStageData)
	activity.Data = &jsonString
	err = tx.Activity.Save(activity)
	return
}
