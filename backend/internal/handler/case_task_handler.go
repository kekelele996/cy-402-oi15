package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"cylawcase/internal/constants"
	"cylawcase/internal/dto"
	"cylawcase/internal/middleware"
	"cylawcase/internal/service"
	"cylawcase/internal/util"

	"github.com/gin-gonic/gin"
)

// CaseTaskHandler 案件待办 HTTP 处理器。
type CaseTaskHandler struct {
	svc    *service.CaseTaskService
	logger *slog.Logger
}

// NewCaseTaskHandler 构造案件待办处理器。
func NewCaseTaskHandler(svc *service.CaseTaskService, logger *slog.Logger) *CaseTaskHandler {
	return &CaseTaskHandler{svc: svc, logger: logger}
}

// ListByCase 按案件查看待办（支持 group=overdue/today/pending/done 过滤，返回分组统计）。
func (h *CaseTaskHandler) ListByCase(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTask list: invalid case id")
		return
	}
	tasks, summary, err := h.svc.ListByCase(caseID, c.Query("group"))
	if err != nil {
		h.wrapError(c, err, "CaseTask list by case failed")
		return
	}
	OK(c, gin.H{"list": tasks, "summary": summary})
}

// Create 创建待办事项。
func (h *CaseTaskHandler) Create(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTask create: invalid case id")
		return
	}
	var req dto.CaseTaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTask create: "+err.Error())
		return
	}
	dueDate, err := dto.ParseDueDate(valueOrEmpty(req.DueDate))
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTask create: invalid due_date")
		return
	}
	t, err := h.svc.Create(caseID, middleware.GetUserID(c), req.AssigneeID, req.Title, dueDate)
	if err != nil {
		h.wrapError(c, err, "CaseTask[case_id="+strconv.FormatUint(caseID, 10)+"] create failed")
		return
	}
	OKWithMessage(c, constants.MsgCaseTaskCreated, t)
}

// Complete 完成待办事项。
func (h *CaseTaskHandler) Complete(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTask complete: invalid task id")
		return
	}
	t, err := h.svc.Complete(taskID)
	if err != nil {
		h.wrapError(c, err, "CaseTask complete failed")
		return
	}
	OKWithMessage(c, constants.MsgCaseTaskCompleted, t)
}

// Assignees 待办负责人候选（该案件主办 + 协办人员）。
func (h *CaseTaskHandler) Assignees(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTask assignees: invalid case id")
		return
	}
	users, err := h.svc.AssigneeCandidates(caseID)
	if err != nil {
		h.wrapError(c, err, "CaseTask assignees failed")
		return
	}
	OK(c, users)
}

func (h *CaseTaskHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("case task handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("case task handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
