package task

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/st-man-hori/go-feature-template/internal/apperror"
	"github.com/st-man-hori/go-feature-template/internal/models"
)

// UseCase holds the Task feature's business logic — the equivalent of the
// classes under Laravel's UseCases/, one execute() each. Following the same
// rule the Laravel template states explicitly ("no Repository pattern —
// UseCases talk to Eloquent directly"), UseCase here talks to *gorm.DB
// directly instead of going through a repository interface.
type UseCase struct {
	db *gorm.DB
}

func NewUseCase(db *gorm.DB) *UseCase {
	return &UseCase{db: db}
}

func (uc *UseCase) Index(req IndexTaskRequest) ([]models.Task, error) {
	query := uc.db.Model(&models.Task{})
	if req.IsDone != nil {
		query = query.Where("is_done = ?", *req.IsDone)
	}

	var tasks []models.Task
	if err := query.Order("created_at DESC").Limit(req.PerPage).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (uc *UseCase) Show(id uint) (models.Task, error) {
	return uc.find(id)
}

func (uc *UseCase) Store(req StoreTaskRequest) (models.Task, error) {
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
	}

	if req.DueDate != nil {
		d, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return models.Task{}, err // unreachable: StoreTaskRequest already validated the format
		}
		task.DueDate = &d
	}

	if err := uc.db.Create(&task).Error; err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (uc *UseCase) Update(id uint, req UpdateTaskRequest) (models.Task, error) {
	task, err := uc.find(id)
	if err != nil {
		return models.Task{}, err
	}

	if req.Title.Set {
		task.Title = req.Title.Value
	}

	if req.Description.Set {
		task.Description = req.Description.Value
	}

	if req.DueDate.Set {
		if req.DueDate.Value == nil {
			task.DueDate = nil
		} else {
			d, err := time.Parse("2006-01-02", *req.DueDate.Value)
			if err != nil {
				return models.Task{}, err // unreachable: UpdateTaskRequest.Validate already checked the format
			}
			task.DueDate = &d
		}
	}

	if req.IsDone.Set {
		task.IsDone = req.IsDone.Value
	}

	if err := uc.db.Save(&task).Error; err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (uc *UseCase) Destroy(id uint) error {
	task, err := uc.find(id)
	if err != nil {
		return err
	}

	return uc.db.Delete(&task).Error
}

func (uc *UseCase) find(id uint) (models.Task, error) {
	var task models.Task

	err := uc.db.First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, apperror.NotFound("Task")
	}
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}
