package ftp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var (
	sessions = make(map[string]*session) // Хранение сессий для докачки
)

type session struct {
	fileName string
	fileSize int64
	progress int64
}

func handleUpload(conn net.Conn, command string) {
	fileName := strings.TrimSpace(strings.TrimPrefix(command, "UPLOAD"))
	if fileName == "" {
		conn.Write([]byte("Ошибка: имя файла не указано\n"))
		return
	}

	sizeBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, sizeBuf); err != nil {
		fmt.Println("Ошибка при чтении размера файла:", err)
		return
	}
	fileSize := int64(binary.BigEndian.Uint32(sizeBuf))

	sessionKey := conn.RemoteAddr().String()
	sess, exists := sessions[sessionKey]
	if !exists {
		sess = &session{fileName: fileName, fileSize: fileSize}
		sessions[sessionKey] = sess
	}

	filePath := fmt.Sprintf("%s", fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Ошибка при создании файла: %v\n", err)))
		return
	}
	defer file.Close()

	conn.Write([]byte(fmt.Sprintf("PROGRESS %d\n", sess.progress)))

	startTime := time.Now()
	buffer := make([]byte, 1024)
	for sess.progress < fileSize {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Передача файла завершена")
				break
			}
			fmt.Println("Ошибка при чтении данных:", err)
			return
		}
		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
		sess.progress += int64(n)
	}

	duration := time.Since(startTime)
	bitrate := float64(sess.progress) / duration.Seconds()
	conn.Write([]byte(fmt.Sprintf("Файл успешно загружен. Битрейт: %.2f байт/сек\n", bitrate)))

	delete(sessions, sessionKey)
}

func handleDownload(conn net.Conn, command string) {
	fileName := strings.TrimSpace(strings.TrimPrefix(command, "DOWNLOAD"))
	if fileName == "" {
		conn.Write([]byte("Ошибка: имя файла не указано\n"))
		return
	}

	filePath := fmt.Sprintf("%s", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		conn.Write([]byte("Ошибка: файл не найден\n"))
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	sizeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuf, uint32(fileSize))
	if _, err := conn.Write(sizeBuf); err != nil {
		fmt.Println("Ошибка при отправке размера файла:", err)
		return
	}

	startTime := time.Now()
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Ошибка при чтении файла:", err)
			return
		}
		if _, err := conn.Write(buffer[:n]); err != nil {
			fmt.Println("Ошибка при отправке данных:", err)
			return
		}
	}

	duration := time.Since(startTime)
	bitrate := float64(fileSize) / duration.Seconds()
	conn.Write([]byte(fmt.Sprintf("Файл успешно скачан. Битрейт: %.2f байт/сек\n", bitrate)))
}
