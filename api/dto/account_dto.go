package dto

import (
	"erp/gen/db/model"
)

type SignInRequest struct {
	Body struct {
		Email    string `json:"email" doc:"Email of the user"`
		Password string `json:"password" minLength:"8" doc:"Password of the user"`
	}
}

type SignInResponse struct {
	Body struct {
		AccessToken  string          `json:"access_token" doc:"Access token of the user"`
		User         UserDto         `json:"user"`
		U            model.User      `json:"-"`
		UserRelation UserRelationDto `json:"user_relation"`
	}
}

type ChangePasswordRequest struct {
	AcceptLanguageHeader
	Body struct {
		Password string `json:"password"`
		Token    string `json:"token"`
	}
}

type UpdatePasswordRequest struct {
	AuthParams
	Body struct {
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}
}

type ResetPasswordRequest struct {
	AcceptLanguageHeader
	Body struct {
		Email string `json:"email"`
	}
}

type AccountResponse struct {
	Body struct {
		User            UserDto            `json:"user"`
		Profile         ProfileDto         `json:"profile"`
		Role            RoleDto            `json:"role"`
		Company         CompanyDto         `json:"company"`
		CompanyDefaults CompanyDefaultsDto `json:"company_defaults"`
		RoleActions     []RoleActionDto    `json:"role_actions"`
	}
}

type UserRelationDto struct {
	UUID    string     `json:"uuid"`
	Company CompanyDto `json:"company"`
	Profile ProfileDto `json:"profile"`
	Role    RoleDto    `json:"role"`
}

func UserRelationDtoFromModel(m *model.UserRelation) UserRelationDto {
	r := UserRelationDto{}
	r.UUID = m.UUID
	if m.Company.ID != 0 {
		r.Company = CompanyDTOFromModel(&m.Company)
	}
	if m.Role.ID != 0 {
		r.Role = RoleDTOFromModel(&m.Role)
	}
	if m.Profile.ID != 0 {
		r.Profile = ProfileDTOFromModel(&m.Profile)
	}
	return r
}
