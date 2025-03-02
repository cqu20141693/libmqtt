package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/enums"
)

// Channel 定义通道表结构
type Channel struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name" binding:"required"`
	Type     string        `json:"type" binding:"required"`
	ParentID int64         `json:"parentId"`
	Config   string        `json:"config"`
	Desc     string        `json:"description"`
	Children []Channel     `json:"children"`
	MenuType enums.BizType `json:"menuType"`
}

// BuildTree 将通道列表转换为树结构
func BuildTree(channels []Channel) []Channel {
	// 创建一个映射，用于快速查找子节点
	channelMap := make(map[int64][]Channel)
	for _, channel := range channels {
		channelMap[channel.ParentID] = append(channelMap[channel.ParentID], channel)
	}

	// 递归构建树
	var build func(parentID int64) []Channel
	build = func(parentID int64) []Channel {
		nodes := channelMap[parentID]
		for i := range nodes {
			nodes[i].Children = build(nodes[i].ID) // 递归查找子节点
		}
		return nodes
	}

	// 从根节点（ParentID 为 0）开始构建树
	return build(-1)
}

// ChannelAdd 新增通道
func (d *DBKit) ChannelAdd(channel Channel) (int64, error) {
	query := `INSERT INTO channels (name, type, parent_id, config, desc) VALUES (?, ?, ?, ?, ?);`
	stmt, _ := d.db.Prepare(query)

	result, err := stmt.Exec(channel.Name, channel.Type, channel.ParentID, channel.Config, channel.Desc)
	if err != nil {
		return 0, fmt.Errorf("failed to add channel: %v", err)
	}
	return result.LastInsertId()
}

// ChannelDelete 删除通道
func (d *DBKit) ChannelDelete(id int64) error {
	query := `DELETE FROM channels WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no channel found with the given id")
	}
	return nil
}

// ChannelById 查询通道
func (d *DBKit) ChannelById(id int) (*Channel, error) {
	query := `SELECT id, name, type, parent_id, config, desc FROM channels WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	row := stmt.QueryRow(id)

	var channel Channel
	err := row.Scan(&channel.ID, &channel.Name, &channel.Type, &channel.ParentID, &channel.Config, &channel.Desc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("channel not found")
		}
		return nil, fmt.Errorf("failed to get channel: %v", err)
	}
	return &channel, nil
}

// ChannelUpdate 更新通道
func (d *DBKit) ChannelUpdate(channel Channel) error {
	query := `
	UPDATE channels
	SET name = ?, type = ?, parent_id = ?, config = ?, desc = ?
	WHERE id = ?;`
	stmt, _ := d.db.Prepare(query)
	result, err := stmt.Exec(channel.Name, channel.Type, channel.ParentID, channel.Config, channel.Desc, channel.ID)
	if err != nil {
		return fmt.Errorf("failed to update channel: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no channel found with the given id")
	}
	return nil
}

// ChannelAll 查询所有通道
func (d *DBKit) ChannelAll() ([]Channel, error) {
	query := `SELECT id, name, type, parent_id, config, desc FROM channels order by id asc;`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %v", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var channels []Channel
	for rows.Next() {
		var channel Channel
		if err := rows.Scan(&channel.ID, &channel.Name, &channel.Type, &channel.ParentID, &channel.Config, &channel.Desc); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %v", err)
		}
		if channel.Type == constants.ChannelGroupKey {
			channel.MenuType = enums.ChannelGroup
		} else {
			channel.MenuType = enums.Channel
		}
		channels = append(channels, channel)
	}

	return BuildTree(channels), nil
}
