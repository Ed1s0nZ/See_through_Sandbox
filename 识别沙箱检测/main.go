package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

func See_through() {
	// 1. 延时运行
	timeSleep1, _ := timeSleep()
	// 2. 检测开机时间
	bootTime1, _ := bootTime()
	// 3. 检测物理内存
	physicalMemory1, _ := physicalMemory()
	// 4. 检测CPU核心数
	numberOfCPU1, _ := numberOfCPU()
	// 5. 检测临时文件数
	numberOfTempFiles1, _ := numberOfTempFiles()
	level := timeSleep1 + bootTime1 + physicalMemory1 + numberOfCPU1 + numberOfTempFiles1 // 有五个等级，等级越趋向于5，越像真机
	fmt.Println("level:", level)
	if level < 4 {
		fmt.Println("可能是沙箱！")
	} else {
		fmt.Println("可能是真机！")
	}
}

// 1. 延时运行
func timeSleep() (int, error) {
	startTime := time.Now()
	time.Sleep(10 * time.Second)
	endTime := time.Now()
	sleepTime := endTime.Sub(startTime)
	if sleepTime >= time.Duration(10*time.Second) {
		fmt.Println("睡眠时间为:", sleepTime)
		return 1, nil
	} else {
		return 0, nil
	}
}

// 2. 检测开机时间
// 许多沙箱检测完毕后会重置系统，我们可以检测开机时间来判断是否为真实的运行状况。
func bootTime() (int, error) {
	var kernel = syscall.NewLazyDLL("Kernel32.dll")
	GetTickCount := kernel.NewProc("GetTickCount")
	r, _, _ := GetTickCount.Call()
	if r == 0 {
		return 0, nil
	}
	ms := time.Duration(r * 1000 * 1000)
	fmt.Println("开机时常为:", ms)
	tm := time.Duration(30 * time.Minute)
	if ms < tm {
		return 0, nil
	} else {
		return 1, nil
	}

}

// 3. 检测物理内存
// 当今大多数pc具有4GB以上的RAM,我们可以检测RAM是否大于4GB来判断是否是真实的运行机器。
func physicalMemory() (int, error) {
	var mod = syscall.NewLazyDLL("kernel32.dll")
	var proc = mod.NewProc("GetPhysicallyInstalledSystemMemory")
	var mem uint64
	proc.Call(uintptr(unsafe.Pointer(&mem))) // ret, _, err := proc.Call(uintptr(unsafe.Pointer(&mem)))
	mem = mem / 1048576                      // 该值以千字节为单位给出，因此除以 1048576 (1024*1024) 以获得以千兆字节为单位的值。
	if mem < 4 {
		fmt.Printf("物理内存为%dG\n", mem)
		return 0, nil // 小于4GB返回0
	}
	fmt.Printf("物理内存为%dG\n", mem)
	return 1, nil // 大于4GB返回1
}

// 4. 检测CPU核心数
// 大多数pc拥有4核心cpu，许多在线检测的虚拟机沙盘是2核心，我们可以通过核心数来判断是否为真实机器或检测用的虚拟沙箱。
func numberOfCPU() (int, error) {
	a := runtime.NumCPU()
	fmt.Println("CPU核心数为:", a)
	if a < 4 {
		return 0, nil // 小于4核心数,返回0
	} else {
		return 1, nil // 大于4核心数，返回1
	}
}

// 5. 检测临时文件数
// 正常使用的系统，其中用户的临时文件夹中有一定数量的临时文件，可以通过判断临时文件夹内的文件数量来检测是否在沙箱中运行。
func numberOfTempFiles() (int, error) {
	conn := os.Getenv("temp") // 通过环境变量读取temp文件夹路径
	var k int
	if conn == "" {
		fmt.Println("未找到temp文件夹，或temp文件夹不存在")
		return 0, nil
	} else {
		local_dir := conn
		err := filepath.Walk(local_dir, func(filename string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				return nil
			}
			k++
			// fmt.Println("filename:", filename)  // 输出文件名字
			return nil
		})
		fmt.Println("Temp总共文件数量:", k)
		if err != nil {
			// fmt.Println("路径获取错误")
			return 0, nil
		}
	}
	if k < 30 {
		return 0, nil
	}
	return 1, nil

}

func main() {
	See_through()
}
