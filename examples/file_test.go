package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetProjectPath(t *testing.T) {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	filePath := path.Dir(filename)
	packagePath := path.Dir(filePath)
	root := path.Dir(packagePath)
	fmt.Println("project path=", root)
}

func TestReadAndWirteFile(t *testing.T) {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	filePath := path.Dir(filename)
	packagePath := path.Dir(filePath)
	dir := packagePath + "/resources/channels"
	// 读取文件内容
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断文件是否为 JSON 文件
		if filepath.Ext(path) == ".json" {
			// 读取 JSON 文件内容
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// 解析 JSON 数据
			var jsonData map[string]interface{}
			err = json.Unmarshal(data, &jsonData)
			if err != nil {
				return err
			}
			config := jsonData["server"].(map[string]interface{})
			deviceConfigs := config["devices"].([]interface{})
			for _, deviceConifg := range deviceConfigs {
				device := deviceConifg.(map[string]interface{})
				tf := device["tags"].([]interface{})

				if tt, ok := device["product_name"]; ok {
					if tt != nil {
						device["deviceTypeName"] = tt
						delete(device, "product_name")
					}
				}
				if _, ok := device["deviceTypeName"]; !ok {
					device["deviceTypeName"] = "go测试自设备"
				}
				for _, tagConf := range tf {
					tag := tagConf.(map[string]interface{})
					extension := tag["extension"].(map[string]interface{})
					if dt, ok := extension["type"]; ok {
						tag["dataType"] = dt
						delete(extension, "type")
					}
					m := tag["dataType"].(string)
					if m == "int" {
						tag["dataType"] = "Integer"
					} else {
						tag["dataType"] = strings.Title(m)
					}
					if tt, ok := tag["tag_type"]; ok {
						tag["tagType"] = tt
						delete(tag, "tag_type")
					}
					if f, ok := tag["scanPeriodInMillis"]; ok {
						tag["frequency"] = f
						delete(tag, "scanPeriodInMillis")
					}
					if r, ok := tag["readOnly"]; ok {
						read := r.(bool)
						if read {
							tag["readWrite"] = 1
						} else {
							tag["readWrite"] = 3
						}
					}

				}
			}
			marshal, _ := json.Marshal(jsonData)
			err = os.WriteFile(path, marshal, 0644)
			if err != nil {
				fmt.Println("写入文件失败:", err)
				return err
			}
		}

		return nil
	})

}
