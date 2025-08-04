package module

import (
	"context"
	// "fmt"
	// "erp/api/dto"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/system"
	module_event "erp/project/core/module/handler/event"
	module_rest "erp/project/core/module/handler/rest"
	module_fsm "erp/project/core/module/pkg/fsm"
	module_repo "erp/project/core/module/repository"
	module_ucase "erp/project/core/module/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	fsm := module_fsm.NewModuleFsm()
	moduleEventRepo := module_repo.NewModuleEventRepo()

	moduleRepo := module_repo.NewModuleRepository(svc.DBConn(), svc.Helpers())
	moduleUseCase := module_ucase.NewModuleUcase(
		svc.Logger(), moduleRepo, svc.PermissionService(), svc.CoreService(), fsm,
	)
	module_rest.NewHandler(svc.HumaApi(),
		svc.Helpers(), moduleUseCase, huma.Middlewares{svc.Middlewares().Authenticate},
		svc.PermissionService())
	module_event.NewModuleEventHandler(svc.Logger(),moduleEventRepo,svc.EventBus())

	svc.Container().AddScoped(domain.ModuleUseCase, func(c di.Container) (any, error) {
		return moduleUseCase,nil
	})
	
	// err = moduleEventRepo.CreateCompany(svc.DBConn().GetQ().Begin(), context.Background(), 275, []dto.CompanyModule{
	// 	{Label: "Contabilidad", Name: "Accounting"},
	// 	{Label: "Cuentas por Pagar", Name: "Payables"},
	// 	{Label: "Cuentas por Cobrar", Name: "Receivables"},
	// 	{Label: "Compra", Name: "Buying"},
	// 	{Label: "Venta", Name: "Selling"},
	// 	{Label: "Inventario", Name: "Stock"},
	// })
	// if err != nil {
	// 	fmt.Println("MODULE ERROR", err)
	// }
	return nil
}
