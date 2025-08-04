package module_repo

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"fmt"

	// "erp/internal/domain/event"
	"strings"
)

type ModuleEventRepo interface {
	// CreateCompany(ctx context.Context,d event.CreatedCompanyEventData)(error)
	CreateCompany(tx *query.QueryTx, ctx context.Context,companyID int64, companyModules []dto.CompanyModule) (err error)
}

type moduleEventRepo struct {
	modules map[string][]SectionEntity
}

func NewModuleEventRepo() ModuleEventRepo {
	return &moduleEventRepo{}
}

// func (r *moduleEventRepo) CreateCompany(ctx context.Context,payload event.CreatedCompanyEventData)(err error) {
func (r *moduleEventRepo) CreateCompany(tx *query.QueryTx, ctx context.Context,
	companyID int64, companyModules []dto.CompanyModule) (err error) {
	// tx := payload.Tx
	fmt.Println("COMPANY MODULES",companyModules)
	defaultModules := DefaultSectionEntities()
	var modules []*model.Module
	var moduleSections []*model.ModuleSection
	for _, module := range companyModules {
		sections, ok := defaultModules[module.Name]
		if !ok {
			continue
		}

		moduleQ := tx.Module
		moduleID, err := moduleQ.WithContext(ctx).InsertParty(proto.PartyType_module.String())
		if err != nil {
			return err
		}
		moduleM := &model.Module{
			ID:        moduleID,
			Label:     module.Label,
			Href:      strings.ToLower(module.Name),
			CompanyID: companyID,
			IconCode: &module.IconCode,
			IconName: &module.IconName,
			Priority: module.Priority,
		}
		modules = append(modules, moduleM)
		for _, section := range sections {
			moduleSection := &model.ModuleSection{
				Name:     section.SectionName,
				ModuleID: moduleM.ID,
				EntityID: int32(section.EntityID),
			}
			moduleSections = append(moduleSections, moduleSection)
		}
	}
	err = tx.Module.CreateInBatches(modules, len(modules))
	if err != nil {
		return
	}
	err = tx.ModuleSection.CreateInBatches(moduleSections, len(moduleSections))
	if err != nil {
		return
	}
	return
}

type SectionEntity struct {
	SectionName string
	EntityID    int64
}
