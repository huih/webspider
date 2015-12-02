package task

// 基础包
import (
	"spider"
)

func init() {
	csdn.Register()
}

var csdn = &spider.Spider{
	Name:        "csdn数据",
	Description: "csdn数据，获取最新的数据",
	PauseTime: 10,
	Address:"http://blog.csdn.net/",
}
