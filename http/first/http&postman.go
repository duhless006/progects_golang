package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var message []int
var messages = make(map[int]string)
var count int
var mtx sync.Mutex

func massageheandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ошибка чтения", err)
		return
	}
	mass := string(httpRequestBody)
	fmt.Println("new massage:\n", mass)

}

func saveMassageHeandler(w http.ResponseWriter, r *http.Request) {

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ошибка чтения", err)
	}

	massages := string(httpRequestBody)
	lines := strings.Split(strings.TrimSpace(massages), "\n")

	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		fmt.Println("строка не может быть пустой", err)
		return
	}

	mtx.Lock()
	for _, line := range lines {
		if line != "" {
			count++
			messages[count] = line
			message = append(message, count)

		}
	}
	mtx.Unlock()

	fmt.Println("Все сообщения прочитаны и сохранены:")
	for _, id := range message {
		fmt.Printf("ID: %d, Сообщение: %s\n", id, messages[id])
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Сохранено %d сообщение(й)", len(lines))))

}

func deleteHeandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ошибка чтения", err)
		return
	}

	idStr := strings.TrimSpace(string(httpRequestBody))
	if idStr == "" {
		fmt.Println("пустая строка")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("ошибка преобразования в integer", err)
		return
	}
	mtx.Lock()
	if _, ok := messages[id]; !ok {
		mtx.Unlock()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Message not found"))
		return
	}

	delete(messages, id)

	for i, msgId := range message {
		if msgId == id {
			message = append(message[:i], message[i+1:]...)
			break
		}
	}
	mtx.Unlock()

	fmt.Printf("Сообщение с ID %d удалено\n", id)

	w.Write([]byte(fmt.Sprintf("Сообщение с ID %d успешно удалено", id)))

	for _, v := range messages {
		fmt.Println(v)
	}

}

func main() {

	messages = make(map[int]string)
	message = make([]int, 0)

	http.HandleFunc("/massage", massageheandler)
	http.HandleFunc("/save", saveMassageHeandler)
	http.HandleFunc("/delete", deleteHeandler)

	port := ":9091"
	fmt.Printf("Сервер запускается на порту %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Ошибка сервера:", err)
		return
	}
}
