package kb

import (
	"fmt"
	"graphraggo/internal/global"
	"os"
)

const (
	BaseDir = "/kb"
)

type KB struct {
	Name string `json:"name,omitempty"`
}

// ReadKB wd/kb 文件夹下所有的文件夹
//
//	wd: 工作目录
func ReadKB() ([]*KB, error) {
	path := fmt.Sprintf("%s/%s", global.WorkDir, global.KBDir)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	kbs := []*KB{}
	for _, file := range files {
		if file.Type().IsDir() {
			kbs = append(kbs, &KB{
				Name: file.Name(),
			})
		}
	}

	return kbs, nil
}