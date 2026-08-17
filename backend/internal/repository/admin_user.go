package repository

import (
	"infinite-canvas/backend/internal/model"

	"gorm.io/gorm"
)

// DeleteUserAccount 删除用户主体及其业务数据。系统渠道、审计记录、兑换码批次保留。
func (r *Repository) DeleteUserAccount(userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var account model.User
		if err := tx.Select("id", "email").Where("id = ?", userID).First(&account).Error; err != nil {
			return err
		}
		projectIDs, err := pluckIDs(tx, &model.Project{}, "user_id = ?", userID)
		if err != nil {
			return err
		}
		assetIDs, err := pluckIDs(tx, &model.Asset{}, "user_id = ?", userID)
		if err != nil {
			return err
		}
		versionIDs, err := pluckIDs(tx, &model.AssetVersion{}, "asset_id IN ?", assetIDs)
		if err != nil {
			return err
		}
		voiceIDs, err := pluckIDs(tx, &model.VoiceProfile{}, "user_id = ?", userID)
		if err != nil {
			return err
		}
		// 用户渠道含密钥，软删会把密钥留在库里；连同已软删记录一并物理删除。
		channelIDs, err := pluckIDs(tx.Unscoped(), &model.ModelChannel{}, "user_id = ? AND scope = ?", userID, model.ChannelScopeUser)
		if err != nil {
			return err
		}
		instanceIDs, err := pluckIDs(tx, &model.WorkflowInstance{}, "project_id IN ?", projectIDs)
		if err != nil {
			return err
		}
		stepIDs, err := pluckIDs(tx, &model.WorkflowStepInstance{}, "workflow_instance_id IN ?", instanceIDs)
		if err != nil {
			return err
		}
		shotIDs, err := pluckIDs(tx, &model.Shot{}, "project_id IN ?", projectIDs)
		if err != nil {
			return err
		}

		type accountDelete struct {
			query    string
			args     []any
			model    any
			unscoped bool
		}
		deletes := []accountDelete{
			{"user_id = ?", []any{userID}, &model.AuthSession{}, false},
			{"user_id = ?", []any{userID}, &model.UserIdentity{}, false},
			{"user_id = ?", []any{userID}, &model.TaskTextDelta{}, false},
			{"user_id = ?", []any{userID}, &model.TaskLog{}, false},
			{"user_id = ?", []any{userID}, &model.Message{}, false},
			{"user_id = ?", []any{userID}, &model.SessionFile{}, false},
			{"user_id = ?", []any{userID}, &model.Result{}, false},
			{"user_id = ?", []any{userID}, &model.Session{}, false},
			{"user_id = ?", []any{userID}, &model.Task{}, false},
			{"user_id = ?", []any{userID}, &model.ApiCallLog{}, false},
			{"user_id = ?", []any{userID}, &model.CreditLedgerEntry{}, false},
			{"user_id = ?", []any{userID}, &model.BillingOrder{}, false},
			{"user_id = ?", []any{userID}, &model.CreditAccount{}, false},
			{"user_id = ?", []any{userID}, &model.UserDailyActivity{}, false},
			{"user_id = ?", []any{userID}, &model.UserDailyUploadUsage{}, false},
			{"user_id = ?", []any{userID}, &model.UserOSSSetting{}, false},
			{"user_id = ?", []any{userID}, &model.UserSkillState{}, false},
			{"user_id = ?", []any{userID}, &model.UserPromptCustomization{}, false},
			{"user_id = ?", []any{userID}, &model.UserAnnouncementRead{}, false},
			{"user_id = ?", []any{userID}, &model.CanvasShare{}, false},
			{"project_id IN ?", []any{projectIDs}, &model.CanvasUnitLink{}, false},
			{"workflow_step_id IN ?", []any{stepIDs}, &model.WorkflowStepTask{}, false},
			{"workflow_instance_id IN ?", []any{instanceIDs}, &model.WorkflowStepInstance{}, false},
			{"project_id IN ?", []any{projectIDs}, &model.WorkflowInstance{}, false},
			{"shot_id IN ?", []any{shotIDs}, &model.ShotAssetReference{}, false},
			{"project_id IN ?", []any{projectIDs}, &model.Shot{}, false},
			{"project_id IN ?", []any{projectIDs}, &model.ProjectAssetLink{}, false},
			{"project_id IN ?", []any{projectIDs}, &model.ProjectAssetCandidate{}, false},
			{"project_id IN ?", []any{projectIDs}, &model.ProjectUnit{}, false},
			{"id IN ?", []any{projectIDs}, &model.Project{}, false},
			{"user_id = ?", []any{userID}, &model.CanvasProject{}, false},
			{"user_id = ?", []any{userID}, &model.StyleProfile{}, false},
			{"asset_version_id IN ?", []any{versionIDs}, &model.CharacterVoiceBinding{}, false},
			{"voice_profile_id IN ?", []any{voiceIDs}, &model.CharacterVoiceBinding{}, false},
			{"asset_version_id IN ?", []any{versionIDs}, &model.AssetRepresentation{}, false},
			{"id IN ?", []any{versionIDs}, &model.AssetVersion{}, false},
			{"id IN ?", []any{assetIDs}, &model.Asset{}, false},
			{"user_id = ?", []any{userID}, &model.VoiceProfile{}, false},
			{"user_id = ?", []any{userID}, &model.Resource{}, false},
			{"channel_id IN ?", []any{channelIDs}, &model.ChannelModel{}, true},
			{"channel_id IN ?", []any{channelIDs}, &model.ModelPricing{}, false},
			{"id IN ?", []any{channelIDs}, &model.ModelChannel{}, true},
		}
		if account.Email != "" {
			deletes = append(deletes, accountDelete{"email = ?", []any{account.Email}, &model.EmailVerificationCode{}, false})
		}
		for _, item := range deletes {
			if err := deleteWhere(tx, item.model, item.query, item.unscoped, item.args...); err != nil {
				return err
			}
		}
		return tx.Delete(&model.User{}, "id = ?", userID).Error
	})
}

func pluckIDs(tx *gorm.DB, value any, query string, args ...any) ([]string, error) {
	ids := []string{}
	for _, arg := range args {
		if idsArg, ok := arg.([]string); ok && len(idsArg) == 0 {
			return ids, nil
		}
	}
	err := tx.Model(value).Where(query, args...).Pluck("id", &ids).Error
	return ids, err
}

func deleteWhere(tx *gorm.DB, value any, query string, unscoped bool, args ...any) error {
	for _, arg := range args {
		if ids, ok := arg.([]string); ok && len(ids) == 0 {
			return nil
		}
	}
	db := tx
	if unscoped {
		db = tx.Unscoped()
	}
	return db.Where(query, args...).Delete(value).Error
}
