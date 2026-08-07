package task_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/st-man-hori/go-feature-template/internal/models"
	"github.com/st-man-hori/go-feature-template/internal/testutil"
)

func TestShow_ReturnsTheTask(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	task := models.Task{Title: "Write the quarterly report"}
	require.NoError(t, db.Create(&task).Error)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data struct {
			ID    uint   `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, task.ID, resp.Data.ID)
	assert.Equal(t, "Write the quarterly report", resp.Data.Title)
}

func TestShow_Returns404WhenMissing(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "not_found", resp.Error.Code)
}
