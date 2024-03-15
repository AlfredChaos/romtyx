package database

type Devices struct {
	ModelBase
	UserID int64  `gorm:"type:BIGINT;" json:"user_id"`
	Name   string `gorm:"type:VARCHAR(255)" json:"name"`
}

type DeviceFileters struct {
}

func (d *Devices) TableName() string {
	return "devices"
}

func (d *Devices) Create() error {
	return DbClient().Create(d).Error
}

func (d *Devices) Delete(id int64) error {
	if d.ID == 0 {
		d.ID = id
	}
	return DbClient().Delete(d).Error
}

func (d *Devices) Update(id int64, values map[string]interface{}) error {
	if d.ID == 0 {
		d.ID = id
	}
	return DbClient().First(d, "id = ?", id).Error
}

func (d *Devices) Get(id int64) error {
	return DbClient().First(d, "id = ?", id).Error
}

func (d *Devices) List(filters *DeviceFileters, devices []Devices) error {
	if filters != nil {
		return nil
	}
	return nil
}
