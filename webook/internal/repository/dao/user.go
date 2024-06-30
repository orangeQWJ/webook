package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("唯一约束字段冲突")
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
			return ErrUserDuplicate
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

func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	// err 数据没找到/数据库出错
	return u, err
}
func (dao *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	// err 数据没找到/数据库出错
	return u, err
}

func (dao *UserDao) UpdateProfile(ctx context.Context, u User) error {
	result := dao.db.Model(&User{}).Where("id = ?", u.Id).Update("Nickname", u.Nickname)
	result = dao.db.Model(&User{}).Where("id = ?", u.Id).Update("AboutMe", u.AboutMe)
	result = dao.db.Model(&User{}).Where("id = ?", u.Id).Update("Birthday", u.Birthday)
	return result.Error
}

// User 直接对应数据库表结构
// 有些人叫做 entity / model / po (persistent object)
type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string
	Birthday string
	AboutMe  string
	Phone    sql.NullString `gorm:"unique"` //唯一索引允许有多个空时,这样设置
	Ctime    int64 // 创建时间
	Utime    int64 // 更新时间
	//Phone    string `gorm:"unique"` // 空字符串相互冲突

}
