package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	// todo Create返回值的含义
	err := dao.db.WithContext(ctx).Create(&u).Error
	// 与mysql数据库强耦合
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	// err 数据没找到/数据库出错
	return u, err
}

// User 直接对应数据库表结构
// 有些人叫做 entity / model / po (persistent object)
type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
