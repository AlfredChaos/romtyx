package database

type Files struct {
	ModelBase
	Path        string `gorm:"type:VARCHAR(255)" json:"path"`
	LastVersion int64  `gorm:"type:BIGINT;" json:"last_version"`
	FileHash    []byte `gorm:"type:VARBINARY(255)" json:"file_hash"`
	OwnerID     int64  `gorm:"type:BIGINT;" json:"owner_id"`
}

type FileFilters struct {
	OwnerID *int64
}

func (f *Files) TableName() string {
	return "files"
}

func (f *Files) Create() error {
	return DbClient().Create(f).Error
}

func (f *Files) Delete(id int64) error {
	if f.ID == 0 {
		f.ID = id
	}
	return DbClient().Delete(f).Error
}

func (f *Files) Update(id int64, values map[string]interface{}) error {
	if f.ID == 0 {
		f.ID = id
	}
	return DbClient().Model(f).Updates(values).Error
}

func (f *Files) Get(id int64) error {
	return DbClient().First(f, "id = ?", id).Error
}

func (f *Files) List(filters *FileFilters, files []Files) error {
	if filters != nil {
		return nil
	}
	return nil
}
