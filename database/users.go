package database

import "time"

type Users struct {
	ModelBase
	LastSignIn   time.Time `gorm:"type:datetime" json:"last_sign_in"`
	UserName     string    `gorm:"type:VARCHAR(255)" json:"user_name"`
	Email        string    `gorm:"type:VARCHAR(255)" json:"email"`
	PasswordHash []byte    `gorm:"type:VARBINARY(255)" json:"password_hash"`
}

type UserFilters struct {
	UserName     *string
	UserNameList []string
	Email        *string
	EmailList    []string
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) Create() error {
	return DbClient().Create(u).Error
}

func (u *Users) Delete(id int64) error {
	if u.ID == 0 {
		u.ID = id
	}
	return DbClient().Delete(u).Error
}

func (u *Users) Update(id int64, values map[string]interface{}) error {
	if u.ID == 0 {
		u.ID = id
	}
	return DbClient().Model(u).Updates(values).Error
}

func (u *Users) Get(id int64) error {
	return DbClient().First(u, "id = ?", id).Error
}

func (u *Users) GetByUserName(name string) error {
	return DbClient().First(u, "user_name = ?", name).Error
}

func (u *Users) GetByEmail(email string) error {
	return DbClient().First(u, "email = ?", email).Error
}

func (u *Users) List(filters *UserFilters, users []Users) error {
	if filters != nil {
		condition := ""
		if filters.UserName != nil && filters.Email == nil {
			condition = "user_name LIKE ?"
			return DbClient().Where(condition, *filters.UserName).Find(&users).Error
		}
		if filters.UserName == nil && filters.Email != nil {
			condition = "email LIKE ?"
			return DbClient().Where(condition, *filters.Email).Find(&users).Error
		}
		if len(filters.EmailList) != 0 {
			condition = "email IN ?"
			return DbClient().Where(condition, filters.EmailList).Find(&users).Error
		}
		if len(filters.UserNameList) != 0 {
			condition = "user_name IN ?"
			return DbClient().Where(condition, filters.UserNameList).Find(&users).Error
		}
	}
	return DbClient().Find(&users).Error
}
