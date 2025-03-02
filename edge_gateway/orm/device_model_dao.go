package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/enums"
)

// DeviceThingsModel 定义设备物模型表结构
type DeviceThingsModel struct {
	ID       int64                 `json:"id"`
	Code     string                `json:"code" binding:"required"`
	Name     string                `json:"name" binding:"required"`
	Type     enums.ThingsModelType `json:"type" binding:"required"`
	DeviceId int                   `json:"DeviceId" binding:"required"`
	Config   string                `json:"config"`
	Desc     string                `json:"description"`
}

// ThingsModelAdd 新增设备物模型
func (d *DBKit) ThingsModelAdd(deviceThingsModel DeviceThingsModel) (int64, error) {
	query := `INSERT INTO device_things_model (code,name,type,device_id, config, desc) VALUES (?, ?, ?, ?, ?);`
	stmt, _ := d.db.Prepare(query)

	result, err := stmt.Exec(deviceThingsModel.Code, deviceThingsModel.Name, deviceThingsModel.Type, deviceThingsModel.DeviceId, deviceThingsModel.Config, deviceThingsModel.Desc)
	if err != nil {
		return 0, fmt.Errorf("failed to add deviceThingsModel: %v", err)
	}
	return result.LastInsertId()
}

// ThingsModelDelete 删除设备物模型
func (d *DBKit) ThingsModelDelete(id int64) error {
	query := `DELETE FROM device_things_model WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete DeviceThingsModel: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no DeviceThingsModel found with the given id")
	}
	return nil
}

// ThingsModelById 查询设备物模型
func (d *DBKit) ThingsModelById(id int64) (*DeviceThingsModel, error) {
	query := `SELECT id, code,name,type, device_id, config, desc FROM device_things_model WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	row := stmt.QueryRow(id)

	var deviceThingsModel DeviceThingsModel
	err := row.Scan(&deviceThingsModel.ID, &deviceThingsModel.Code, &deviceThingsModel.Name, &deviceThingsModel.Type, deviceThingsModel.DeviceId, &deviceThingsModel.Config, &deviceThingsModel.Desc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("deviceThingsModel not found")
		}
		return nil, fmt.Errorf("failed to get deviceThingsModel: %v", err)
	}
	return &deviceThingsModel, nil
}

// ThingsModelUpdate 更新设备物模型
func (d *DBKit) ThingsModelUpdate(deviceThingsModel DeviceThingsModel) error {
	query := `
	UPDATE device_things_model
	SET name = ?, type = ?, device_id = ?, config = ?, desc = ?
	WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	result, err := stmt.Exec(deviceThingsModel.Name, deviceThingsModel.Type, deviceThingsModel.DeviceId, deviceThingsModel.Config, deviceThingsModel.Desc, deviceThingsModel.ID)
	if err != nil {
		return fmt.Errorf("failed to update deviceThingsModel: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no deviceThingsModel found with the given id")
	}
	return nil
}

// ThingsModelByDevice 查询设备物模型
func (d *DBKit) ThingsModelByDevice(DeviceId int64) ([]DeviceThingsModel, error) {
	query := `SELECT id,code, name, type, device_id, config, desc FROM device_things_model 
                                          where device_id=?
                                          order by id asc;`
	prepare, _ := d.db.Prepare(query)
	rows, err := prepare.Query(DeviceId)
	if err != nil {
		return nil, fmt.Errorf("failed to query DeviceThingsModels: %v", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var DeviceThingsModels []DeviceThingsModel
	for rows.Next() {
		var deviceThingsModel DeviceThingsModel
		if err := rows.Scan(&deviceThingsModel.ID, &deviceThingsModel.Code, &deviceThingsModel.Name, &deviceThingsModel.Type, &deviceThingsModel.DeviceId, &deviceThingsModel.Config, &deviceThingsModel.Desc); err != nil {
			return nil, fmt.Errorf("failed to scan deviceThingsModel: %v", err)
		}
		DeviceThingsModels = append(DeviceThingsModels, deviceThingsModel)
	}
	return DeviceThingsModels, nil
}
