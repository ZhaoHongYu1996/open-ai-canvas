//go:build windows

package hostupdate

func availableDiskBytes(string) (int64, error) {
	// Host updater 只在 Linux 部署上执行；Windows 本地编译跳过该预检。
	return 3 << 30, nil
}
