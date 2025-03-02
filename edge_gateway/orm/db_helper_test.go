package orm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewDBKit(t *testing.T) {
	// 设置环境变量
	_ = os.Setenv("DB_PATH", "./data.DB")

	// 初始化数据库工具包
	db, err := NewDBKit()
	if err != nil {
		log.Fatalf("Failed to initialize DBKit: %v", err)
	}
	defer func(db *DBKit) {
		_ = db.Close()
	}(db)

	// 新增通道
	channel := Channel{
		Name:     "Channel 1",
		Type:     "Type A",
		ParentID: 0,
		Config:   "{}",
		Desc:     "Test channel",
	}
	id, err := db.ChannelAdd(channel)
	if err != nil {
		log.Fatalf("Failed to add channel: %v", err)
	}
	fmt.Printf("Added channel with ID: %d\n", id)

	// 查询通道
	retrievedChannel, err := db.ChannelById(int(id))
	if err != nil {
		log.Fatalf("Failed to get channel: %v", err)
	}
	fmt.Printf("Retrieved channel: %+v\n", retrievedChannel)

	// 更新通道
	retrievedChannel.Desc = "Updated description"
	if err := db.ChannelUpdate(*retrievedChannel); err != nil {
		log.Fatalf("Failed to update channel: %v", err)
	}
	fmt.Println("Channel updated")

	// 查询所有通道
	channels, err := db.ChannelAll()
	if err != nil {
		log.Fatalf("Failed to get all channels: %v", err)
	}
	fmt.Printf("All channels: %+v\n", channels)

	// 删除通道
	if err := db.ChannelDelete(id); err != nil {
		log.Fatalf("Failed to delete channel: %v", err)
	}
	fmt.Println("Channel deleted")
}

func TestName(t *testing.T) {
	// 示例数据
	channels := []Channel{
		{ID: 1, Name: "Channel 1", Type: "Type A", ParentID: 0, Config: "{}", Desc: "Root Channel 1"},
		{ID: 2, Name: "Channel 2", Type: "Type B", ParentID: 0, Config: "{}", Desc: "Root Channel 2"},
		{ID: 3, Name: "Channel 3", Type: "Type A", ParentID: 1, Config: "{}", Desc: "Child of Channel 1"},
		{ID: 4, Name: "Channel 4", Type: "Type B", ParentID: 1, Config: "{}", Desc: "Child of Channel 1"},
		{ID: 5, Name: "Channel 5", Type: "Type A", ParentID: 2, Config: "{}", Desc: "Child of Channel 2"},
		{ID: 6, Name: "Channel 6", Type: "Type B", ParentID: 3, Config: "{}", Desc: "Child of Channel 3"},
	}

	// 构建树
	tree := BuildTree(channels)

	// 打印树结构
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}
