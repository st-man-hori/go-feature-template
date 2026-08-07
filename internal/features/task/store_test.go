package task_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/st-man-hori/go-feature-template/internal/models"
	"github.com/st-man-hori/go-feature-template/internal/testutil"
)

func TestStore_CreatesATask(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	payload := map[string]any{
		"title":       "Write the quarterly report",
		"description": "Summarize Q2 sales figures",
		"dueDate":     "2026-08-01",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var resp struct {
		Data struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			DueDate     string `json:"dueDate"`
			IsDone      bool   `json:"isDone"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Write the quarterly report", resp.Data.Title)
	assert.Equal(t, "Summarize Q2 sales figures", resp.Data.Description)
	assert.Equal(t, "2026-08-01", resp.Data.DueDate)
	assert.False(t, resp.Data.IsDone)

	var count int64
	require.NoError(t, db.Model(&models.Task{}).Count(&count).Error)
	assert.EqualValues(t, 1, count)
}

func TestStore_RequiresATitle(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	body, err := json.Marshal(map[string]any{"description": "No title given"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp struct {
		Errors map[string][]string `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Contains(t, resp.Errors, "title")

	var count int64
	require.NoError(t, db.Model(&models.Task{}).Count(&count).Error)
	assert.EqualValues(t, 0, count)
}
