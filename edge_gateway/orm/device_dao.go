package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/enums"
)

// Device 定义设备表结构
type Device struct {
	ID        int64         `json:"id"`
	Code      string        `json:"code" binding:"required"`
	Name      string        `json:"name" binding:"required"`
	Type      string        `json:"type" `
	ChannelId int           `json:"ChannelId" binding:"required"`
	Config    string        `json:"config"`
	Desc      string        `json:"description"`
	MenuType  enums.BizType `json:"menuType"`
}

// DeviceAdd 新增设备
func (d *DBKit) DeviceAdd(device Device) (int64, error) {
	query := `INSERT INTO devices (code,name,type,  channel_id, config, desc) VALUES (?,?, ?, ?, ?, ?);`
	stmt, _ := d.db.Prepare(query)

	result, err := stmt.Exec(device.Code, device.Name, device.Type, device.ChannelId, device.Config, device.Desc)
	if err != nil {
		return 0, fmt.Errorf("failed to add Device: %v", err)
	}
	return result.LastInsertId()
}

// DeviceDelete 删除设备
func (d *DBKit) DeviceDelete(id int64) error {
	query := `DELETE FROM devices WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete Device: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no Device found with the given id")
	}
	return nil
}

// DeviceById 查询设备
func (d *DBKit) DeviceById(id int) (*Device, error) {
	query := `SELECT id, code,name,type, channel_id, config, desc FROM devices WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	row := stmt.QueryRow(id)

	var device Device
	err := row.Scan(&device.ID, &device.Code, &device.Name, &device.Type, device.ChannelId, &device.Config, &device.Desc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("device not found")
		}
		return nil, fmt.Errorf("failed to get device: %v", err)
	}
	return &device, nil
}

// DeviceUpdate 更新设备
func (d *DBKit) DeviceUpdate(device Device) error {
	query := `
	UPDATE devices
	SET name = ?, type = ?, channel_id = ?, config = ?, desc = ?
	WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	result, err := stmt.Exec(device.Name, device.Type, device.ChannelId, device.Config, device.Desc, device.ID)
	if err != nil {
		return fmt.Errorf("failed to update Device: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no Device found with the given id")
	}
	return nil
}

// DeviceByChannel 查询所有设备根据通道
func (d *DBKit) DeviceByChannel(channelId int64) ([]Device, error) {
	query := `SELECT id,code, name, channel_id, config, desc FROM devices where channel_id=? order by id asc;`
	prepare, _ := d.db.Prepare(query)
	rows, err := prepare.Query(channelId)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %v", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var devices []Device
	for rows.Next() {
		var device Device
		if err := rows.Scan(&device.ID, &device.Code, &device.Name, &device.ChannelId, &device.Config, &device.Desc); err != nil {
			return nil, fmt.Errorf("failed to scan device: %v", err)
		}
		device.MenuType = enums.Device
		devices = append(devices, device)
	}
	return devices, nil
}
