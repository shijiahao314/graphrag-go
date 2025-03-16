package global

const (
	KBDir = "kb"
)

var (
	Host               string // 主机地址
	Port               int    // 端口号
	PythonServerPort   int    // Python 服务端口号
	WorkDir            string // 本项目的绝对路径
	ExampleSettingFile string // 示例 Settings 文件路径
	PythonPath         string // Conda 环境下 Python 路径
)
