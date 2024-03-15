package database

type FileVersions struct {
	ModelBase
	FileID        int64 `gorm:"type:BIGINT;" json:"file_id"`
	VersionNumber int64 `gorm:"type:BIGINT;" json:"version_number"`
}

type FileVersionFilters struct {
}

func (fv *FileVersions) Create() error {
	return DbClient().Create(fv).Error
}

func (fv *FileVersions) Delete(id int64) error {
	if fv.ID == 0 {
		fv.ID = id
	}
	return DbClient().Delete(fv).Error
}

func (fv *FileVersions) Update(id int64, values map[string]interface{}) error {
	if fv.ID == 0 {
		fv.ID = id
	}
	return DbClient().Model(fv).Updates(values).Error
}

func (fv *FileVersions) Get(id int64) error {
	return DbClient().First(fv, "id = ?", id).Error
}

func (fv *FileVersions) List(filters *FileVersionFilters, fvs []FileVersions) error {
	if filters != nil {
		return nil
	}
	return nil
}
