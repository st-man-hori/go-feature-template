package task_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/st-man-hori/go-feature-template/internal/models"
	"github.com/st-man-hori/go-feature-template/internal/testutil"
)

func TestIndex_ReturnsAListOfTasks(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	require.NoError(t, db.Create(&[]models.Task{
		{Title: "Task A"},
		{Title: "Task B"},
		{Title: "Task C"},
	}).Error)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 3)
	assert.Contains(t, resp.Data[0], "id")
	assert.Contains(t, resp.Data[0], "dueDate")
	assert.Contains(t, resp.Data[0], "isDone")
	assert.Contains(t, resp.Data[0], "createdAt")
}

func TestIndex_FiltersTasksByIsDone(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	require.NoError(t, db.Create(&[]models.Task{
		{Title: "Not done 1"},
		{Title: "Not done 2"},
		{Title: "Done", IsDone: true},
	}).Error)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks?isDone=true", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 1)
}

func TestIndex_RejectsPerPageOutOfRange(t *testing.T) {
	db := testutil.NewDB(t)
	router := testutil.NewRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks?perPage=0", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
