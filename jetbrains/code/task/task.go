package task

import (
	v2 "license/jetbrains/code/service/v2"
	"license/logger"
)

type Task struct {
	ProductService *v2.ProductService
	PluginService  *v2.PluginService
}

func NewTask() *Task {
	return &Task{
		ProductService: v2.NewProductService(),
		PluginService:  v2.NewPluginService(),
	}

}

func (t *Task) FetchProductLatest() {
	err := t.ProductService.FetchLatest()
	if err != nil {
		logger.Error("Failed to fetch latest product:", err)
	}
	err = t.PluginService.FetchLatest()
	if err != nil {
		logger.Error("Failed to fetch latest plugin:", err)
	}
}
