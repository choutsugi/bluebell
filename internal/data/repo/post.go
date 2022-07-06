package repo

import (
	"bluebell/internal/entity"
	"bluebell/internal/pkg/errx"
	"gorm.io/gorm"
)

var _ PostRepo = (*postRepo)(nil)

type PostRepo interface {
	Insert(post *entity.Post) (err error)
	Delete(post *entity.Post) (err error)
	Update(post *entity.Post) (err error)
	FetchByID(postId int64) (post *entity.Post, err error)
	FetchAll() (posts []*entity.Post, err error)
	FetchList(offset, limit int) (posts []*entity.Post, err error)
}

type postRepo struct {
	db *gorm.DB
}

func (repo *postRepo) FetchList(offset, limit int) (posts []*entity.Post, err error) {
	posts = make([]*entity.Post, 0)

	//后端计算分页时：
	if err = repo.db.Offset(offset).Limit(limit).Find(posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errx.ErrPostNotFound
		}
		return nil, err
	}
	return
}

func (repo *postRepo) Insert(post *entity.Post) (err error) {
	return repo.db.Save(&post).Error
}

func (repo *postRepo) Delete(post *entity.Post) (err error) {
	return repo.db.Delete(&post).Error
}

func (repo *postRepo) Update(post *entity.Post) (err error) {
	return repo.db.Save(&post).Error
}

func (repo *postRepo) FetchByID(postId int64) (post *entity.Post, err error) {
	post = new(entity.Post)
	if err = repo.db.Find(post, postId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errx.ErrPostNotFound
		}
		return nil, err
	}
	return
}

func (repo *postRepo) FetchAll() (posts []*entity.Post, err error) {
	posts = make([]*entity.Post, 0)
	err = repo.db.Find(&posts).Error
	return
}

func NewPostRepo(db *gorm.DB) PostRepo {
	return &postRepo{db: db}
}
