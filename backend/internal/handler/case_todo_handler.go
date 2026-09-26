package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"cylawcase/internal/constants"
	"cylawcase/internal/dto"
	"cylawcase/internal/service"
	"cylawcase/internal/util"

	"github.com/gin-gonic/gin"
)

// CaseTodoHandler 案件待办 HTTP 处理器。
type CaseTodoHandler struct {
	svc    *service.CaseTodoService
	logger *slog.Logger
}

// NewCaseTodoHandler 构造案件待办处理器。
func NewCaseTodoHandler(svc *service.CaseTodoService, logger *slog.Logger) *CaseTodoHandler {
	return &CaseTodoHandler{svc: svc, logger: logger}
}

// ListByCase 某案件待办看板（逾期/今日/未完成/已完成分组 + 统计）。
func (h *CaseTodoHandler) ListByCase(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTodo list by case: invalid case id")
		return
	}
	board, err := h.svc.ListByCase(caseID)
	if err != nil {
		h.wrapError(c, err, "CaseTodo list by case failed")
		return
	}
	OK(c, board)
}

// Assignees 某案件可担当负责人的员工（主办 + 协办）。
func (h *CaseTodoHandler) Assignees(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTodo assignees: invalid case id")
		return
	}
	users, err := h.svc.Assignees(caseID)
	if err != nil {
		h.wrapError(c, err, "CaseTodo assignees failed")
		return
	}
	OK(c, users)
}

// Create 创建待办事项。
func (h *CaseTodoHandler) Create(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTodo create: invalid case id")
		return
	}
	var req dto.CaseTodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTodo[case_id="+strconv.FormatUint(caseID, 10)+"] create: "+err.Error())
		return
	}
	dueDate, err := dto.ParseDueDate(*req.DueDate)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTodo create: invalid due_date")
		return
	}
	t, err := h.svc.Create(caseID, req.AssigneeID, req.Title, dueDate)
	if err != nil {
		h.wrapError(c, err, "CaseTodo[case_id="+strconv.FormatUint(caseID, 10)+"] create failed")
		return
	}
	OKWithMessage(c, constants.MsgCaseTodoCreated, t)
}

// Complete 完成待办事项。
func (h *CaseTodoHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "CaseTodo[id] complete: invalid id")
		return
	}
	t, err := h.svc.Complete(id)
	if err != nil {
		h.wrapError(c, err, "CaseTodo complete failed")
		return
	}
	OKWithMessage(c, constants.MsgCaseTodoCompleted, t)
}

func (h *CaseTodoHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("case todo handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("case todo handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
