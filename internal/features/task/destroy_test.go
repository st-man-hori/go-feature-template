package task_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/st-man-hori/go-feature-template/internal/models"
	"github.com/st-man-hori/go-feature-template/internal/testutil"
)

func TestDestroy_DeletesTheTask(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	task := models.Task{Title: "Task to delete"}
	require.NoError(t, db.Create(&task).Error)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	require.NoError(t, db.Model(&models.Task{}).Count(&count).Error)
	assert.EqualValues(t, 0, count)
}

func TestDestroy_Returns404WhenMissing(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
