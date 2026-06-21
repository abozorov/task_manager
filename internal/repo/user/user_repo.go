package repo

import (
	"context"
	"time"

	"github.com/abozorov/task_manager/internal/models"
	"github.com/abozorov/task_manager/pkg/errs"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) Create(ctx context.Context, user models.User) error {
	tx := u.db.WithContext(ctx).Create(&user)
	return errs.PostgresToErrs(tx.Error)
}

func (u *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)
	tx := u.db.WithContext(ctx).Select("id", "name", "deleted_at").Find(&users)

	return users, errs.PostgresToErrs(tx.Error)
}

func (u *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	tx := u.db.WithContext(ctx).Select("id", "name", "deleted_at").Where("id = ?", id).First(user)
	return user, errs.PostgresToErrs(tx.Error)
}

func (u *UserRepo) Update(ctx context.Context, user models.User) error {
	tx := u.db.WithContext(ctx).Model(&user).Updates(models.User{
		ID:   user.ID,
		Name: user.Name,
	})
	return errs.PostgresToErrs(tx.Error)
}

func (u *UserRepo) DeleteByID(ctx context.Context, id int) error {
	user, err := u.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !user.DeletedAt.IsZero() {
		return errs.ErrNotFound
	}
	tx := u.db.WithContext(ctx).Model(&user).Update("deleted_at", time.Now())
	return errs.PostgresToErrs(tx.Error)
}
