package task_test

import (
	"bytes"
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

func TestUpdate_UpdatesOnlyProvidedFields(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	desc := "original description"
	task := models.Task{Title: "Original title", Description: &desc}
	require.NoError(t, db.Create(&task).Error)

	body, err := json.Marshal(map[string]any{"title": "Updated title"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/tasks/%d", task.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Updated title", resp.Data.Title)
	assert.Equal(t, "original description", resp.Data.Description)
}

func TestUpdate_CanExplicitlyClearANullableField(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	desc := "will be cleared"
	task := models.Task{Title: "Task", Description: &desc}
	require.NoError(t, db.Create(&task).Error)

	body, err := json.Marshal(map[string]any{"description": nil})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/tasks/%d", task.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data struct {
			Description *string `json:"description"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Nil(t, resp.Data.Description)
}

func TestUpdate_RejectsAnEmptyTitle(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	task := models.Task{Title: "Task"}
	require.NoError(t, db.Create(&task).Error)

	body, err := json.Marshal(map[string]any{"title": ""})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/tasks/%d", task.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdate_Returns404WhenMissing(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	body, err := json.Marshal(map[string]any{"title": "Doesn't matter"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/api/tasks/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
