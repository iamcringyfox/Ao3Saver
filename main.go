package main

import (
    "fmt"
    "log"
    "net/http"
	  "os"
    "os/exec"
)

// Middleware для обработки CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Устанавливаем заголовки CORS
        w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем любой origin
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // Если это предварительный запрос (OPTIONS), отправляем пустой ответ
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Вызываем следующий обработчик
        next.ServeHTTP(w, r)
    }
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
    // Читаем тело запроса
    body := make([]byte, r.ContentLength)
    r.Body.Read(body)
    defer r.Body.Close()

    // Парсим JSON
    var url string
    fmt.Sscanf(string(body), `{"url":"%s"}`, &url)

    if url == "" {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Выводим полученную ссылку
    log.Printf("Received URL: %s", url)

    // Отправляем успешный ответ клиенту НЕМЕДЛЕННО
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "URL received and processing started")

    // Запускаем gallery-dl в фоновой Goroutine
    go func() {
        cmd := exec.Command("gallery-dl", url)
        cmd.Stdout = os.Stdout // Выводим логи в консоль
        cmd.Stderr = os.Stderr // Выводим ошибки в консоль
        err := cmd.Run()

        if err != nil {
            log.Printf("Error running gallery-dl for URL %s: %v", url, err)
        } else {
            log.Printf("Successfully processed URL: %s", url)
        }
    }()
}

func main() {
    // Создаем HTTP-сервер с CORS-поддержкой
    http.Handle("/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handlePostRequest(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))

    log.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
