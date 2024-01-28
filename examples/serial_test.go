package examples

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"testing"
)

var serialPortNameA = "/dev/ttys004"
var serialPortNameB = "/dev/ttys005"

/*
创建虚拟串口
brew install socat
socat -d -d pty,raw,echo=0 pty,raw,echo=0
*/
func TestOpenSerial(t *testing.T) {
	// 配置串口参数
	options := serial.OpenOptions{
		PortName:        serialPortNameA,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// 打开串口
	port, err := serial.Open(options)
	if err != nil {
		log.Fatal("open error: ", err)
	}

	// 关闭串口
	defer func(port io.ReadWriteCloser) {
		err := port.Close()
		if err != nil {
			log.Println("serial close failed,", err)
		}
	}(port)
}

func TestReadSerial(t *testing.T) {
	// 配置串口参数
	options := serial.OpenOptions{
		PortName:        serialPortNameA,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// 打开串口
	port, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}

	// 关闭串口
	defer func(port io.ReadWriteCloser) {
		err := port.Close()
		if err != nil {
			log.Println("serial close failed,", err)
		}
	}(port)

	// 读取数据
	buf := make([]byte, 128)
	n, err := port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// 输出读取到的数据
	log.Printf("Read %d bytes: %v", n, string(buf[:n]))
}

func TestWriteSerial(t *testing.T) {
	// 配置串口参数
	options := serial.OpenOptions{
		PortName:        serialPortNameB, // 虚拟串口设备的路径，根据实际情况进行修改
		BaudRate:        9600,            // 波特率
		DataBits:        8,               // 数据位
		StopBits:        1,               // 停止位
		MinimumReadSize: 1,               // 最小读取字节数
	}

	// 打开串口
	port, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}

	// 关闭串口
	defer func(port io.ReadWriteCloser) {
		err := port.Close()
		if err != nil {
			log.Println("serial close failed,", err)
		}
	}(port)

	// 写入数据
	buf := []byte("Hello, Serial!")
	n, err := port.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

	// 输出写入的字节数
	log.Printf("Write %d bytes: %v", n, buf)
}

// 生成的串口对能相互读写
func TestAToB(t *testing.T) {
	// 配置虚拟串口参数
	writeOptions := serial.OpenOptions{
		PortName:        serialPortNameA, // 虚拟串口设备的路径，根据实际情况进行修改
		BaudRate:        9600,            // 波特率
		DataBits:        8,               // 数据位
		StopBits:        1,               // 停止位
		MinimumReadSize: 1,               // 最小读取字节数
	}
	// 配置虚拟串口参数
	readOptions := serial.OpenOptions{
		PortName:        serialPortNameB, // 虚拟串口设备的路径，根据实际情况进行修改
		BaudRate:        9600,            // 波特率
		DataBits:        8,               // 数据位
		StopBits:        1,               // 停止位
		MinimumReadSize: 1,               // 最小读取字节数
	}

	// 打开虚拟串口
	writePort, err := serial.Open(writeOptions)
	if err != nil {
		log.Fatalf("无法打开虚拟串口：%v", err)
	}
	// 打开虚拟串口
	readPort, err := serial.Open(readOptions)
	if err != nil {
		log.Fatalf("无法打开虚拟串口：%v", err)
	}
	defer func(writePort io.ReadWriteCloser, readPort io.ReadWriteCloser) {
		err := writePort.Close()
		if err != nil {

		}
		err = readPort.Close()
		if err != nil {
			return
		}
	}(writePort, readPort)

	// 写入数据到虚拟串口
	writeData := []byte("Hello, Serial Port!")
	n, err := writePort.Write(writeData)
	if err != nil {
		log.Fatalf("无法写入数据到虚拟串口：%v", err)
	}
	fmt.Printf("成功写入 %d 字节数据\n", n)

	// 读取虚拟串口数据
	readData := make([]byte, 128)
	n, err = readPort.Read(readData)
	if err != nil && err != io.EOF {
		log.Fatalf("无法读取虚拟串口数据：%v", err)
	}
	fmt.Printf("读取到的数据：%s\n", readData[:n])
}
