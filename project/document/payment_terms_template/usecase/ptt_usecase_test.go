package payment_terms_t_ucase_test

// import (
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/gen/mocks"
// 	"erp/internal/domain"
// 	"erp/pkg/logger"
// 	payment_terms_t_ucase "erp/project/document/payment_terms_template/usecase"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type PttTestSuite struct {
// 	suite.Suite
// 	VariableThatShouldStartAtFive int

// 	defaultData []dto.PaymentTermsTemplateDto
// 	reqContext  *common.RequestContext
// }

// func (suite *PttTestSuite) SetupTest() {
// 	suite.VariableThatShouldStartAtFive = 5
// 	suite.defaultData = []dto.PaymentTermsTemplateDto{
// 		{Name: "50 -50"},
// 	}
// 	suite.reqContext = &common.DEFAULT_REQ_CONTEXT
// }

// func (suite *PttTestSuite) buildPttUseCase(args ...interface{}) payment_terms_t_ucase.PaymentTermsTemplateUcase {
// 	var (
// 		logger     logger.Logger = logger.New()
// 		mockRepo   *mocks.MockPaymentTermsTemplateRepo = new(mocks.MockPaymentTermsTemplateRepo)
// 		core       *mocks.MockCoreService       = new(mocks.MockCoreService)
// 		permission *mocks.MockPermissionService = new(mocks.MockPermissionService)
// 		fsm        *mocks.MockFsmState          = new(mocks.MockFsmState)
// 		bus        *mocks.MockBus               = new(mocks.MockBus)
// 		c          *mocks.MockContainer         = new(mocks.MockContainer)
// 	)
// 	for _, arg := range args {
// 		switch v := arg.(type) {
// 		case *mocks.MockPaymentTermsTemplateRepo:
// 			mockRepo = v
// 		case *mocks.MockCoreService:
// 			core = v
// 		case *mocks.MockPermissionService:
// 			permission = v
// 		case *mocks.MockFsmState:
// 			fsm = v
// 		case *mocks.MockBus:
// 			bus = v
// 		case *mocks.MockContainer:
// 			c = v
// 		}
// 	}
// 	return payment_terms_t_ucase.NewUseCase(
// 		logger,
// 		core,
// 		permission,
// 		mockRepo,
// 		fsm,
// 		bus,
// 		c,
// 	)
// }

// func (suite *PttTestSuite) TestGreet() {
// 	suite.T().Run("Pass arg", func(t *testing.T) {
// 		mockRepo := new(mocks.MockPaymentTermsTemplateRepo)
// 		mockRepo.On("Greet", mock.AnythingOfType("string")).Return("Hello", nil)
// 		u := suite.buildPttUseCase(mockRepo)
// 		got, err := u.Greet("Mike")
// 		assert.Nil(t, err)
// 		assert.Equal(t, got, "Hello, Mike")
// 		// )
// 		// assert.Equal(t, )
// 	})
// }

// func (suite *PttTestSuite) TestGetPaymentTermsTemplates() {
// 	suite.T().Run("success", func(t *testing.T) {
// 		mockRepo := new(mocks.MockPaymentTermsTemplateRepo)
// 		mockRepo.On("GetPaymentTermsTemplates", mock.Anything, mock.Anything).Return(suite.defaultData, nil)
// 		mockRepo.On("GetFilterOptions", mock.AnythingOfType("string")).Return([]dto.FilterOptionDto{}, nil)
// 		mockPermission := new(mocks.MockPermissionService)
// 		mockPermission.On("CheckPermissionEntity", mock.Anything, mock.Anything, 
// 		mock.Anything, mock.Anything).Return(nil)
// 		u := suite.buildPttUseCase(mockRepo, mockPermission)
// 		got, err := u.GetPaymentTermsTemplates(suite.reqContext, dto.PaymentTermsTemplateRequest{})
// 		assert.NoError(t, err)
// 		assert.Len(t, got.Body.Result, len(suite.defaultData))
// 	})
// 	suite.T().Run("fail permission", func(t *testing.T) {
// 		mockPermission := new(mocks.MockPermissionService)
// 		mockPermission.On("CheckPermissionEntity", mock.Anything, mock.Anything, 
// 		mock.Anything, mock.Anything).Return(domain.ACTION_NOT_ALLOWED)
// 		u := suite.buildPttUseCase(mockPermission)
// 		_, err := u.GetPaymentTermsTemplates(suite.reqContext, dto.PaymentTermsTemplateRequest{})
// 		assert.Error(t, err,domain.ACTION_NOT_ALLOWED)
// 		// assert.Empty(t, got.Body.Result)
// 	})
// }

// func TestPaymentTermsTemplateSuite(t *testing.T) {
// 	suite.Run(t, new(PttTestSuite))

// }
