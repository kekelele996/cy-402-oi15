package router

import (
	"testing"

	"cylawcase/internal/config"
	"cylawcase/internal/handler"
	"cylawcase/internal/repository"
	"cylawcase/internal/service"
	"cylawcase/internal/util"

	"log/slog"
)

// TestSetupRegistersCaseTaskRoutes 验证路由装配不 panic，且案件待办路由已注册。
// 构造期不发生任何数据库查询，传 nil db 即可。
func TestSetupRegistersCaseTaskRoutes(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHours: 1, RateLimitPerMinute: 1000}
	logger := util.NewLogger(slog.LevelError)

	userRepo := repository.NewUserRepository(nil)
	clientRepo := repository.NewClientRepository(nil)
	caseRepo := repository.NewCaseRepository(nil)
	caseTaskRepo := repository.NewCaseTaskRepository(nil)
	documentRepo := repository.NewDocumentRepository(nil)
	billingRepo := repository.NewBillingRepository(nil)

	userSvc := service.NewUserService(userRepo, logger)
	clientSvc := service.NewClientService(clientRepo, caseRepo, logger)
	caseSvc := service.NewCaseService(caseRepo, clientRepo, userRepo, caseTaskRepo, logger)
	caseTaskSvc := service.NewCaseTaskService(caseTaskRepo, caseRepo, userRepo, logger)
	documentSvc := service.NewDocumentService(documentRepo, caseRepo, logger)
	billingSvc := service.NewBillingService(billingRepo, caseRepo, clientRepo, logger)

	r := New(cfg, nil, logger,
		handler.NewUserHandler(userSvc, logger),
		handler.NewClientHandler(clientSvc, logger),
		handler.NewCaseHandler(caseSvc, logger),
		handler.NewCaseTaskHandler(caseTaskSvc, logger),
		handler.NewDocumentHandler(documentSvc, logger),
		handler.NewBillingHandler(billingSvc, logger),
		handler.NewUploadHandler(cfg, logger),
		handler.NewAuditLogHandler(nil, logger),
	)
	engine := r.Setup()

	registered := map[string]bool{}
	for _, ri := range engine.Routes() {
		registered[ri.Method+" "+ri.Path] = true
	}
	want := []string{
		"GET /api/v1/cases/:id/tasks",
		"POST /api/v1/cases/:id/tasks",
		"POST /api/v1/cases/:id/tasks/:task_id/complete",
		"GET /api/v1/cases/:id/assignees",
		"POST /api/v1/cases/:id/status",
	}
	for _, w := range want {
		if !registered[w] {
			t.Errorf("route %s not registered", w)
		}
	}
}
