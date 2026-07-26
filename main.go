package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Структура для распределения входящего JSON с логином и паролем пользователя
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	// http.HandleFunc - это все регистрация путей, другими словами это функции, которые у нас вызываются при заходе на любой путь (/health; /register и т.д.)
	http.HandleFunc("/health", handler)
	http.HandleFunc("/register", registerHandler)

	log.Printf("[INFO] Auth Project")
	log.Printf("[INFO] ver. 1.0.0")
	log.Printf("[INFO] Starting...")

	log.Fatal(http.ListenAndServe(":8080", nil)) // Запуск локального сервера с портом 8080, nil - использование настроек по умолчанию
}

// Обработчик /handler
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// Обработчик регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		// Если наш метод не POST, то возвращаем 405 Method Not Allowed (Метод не поддерживается), т.е. мы используем здесь только POST (Отправить данные)
		return
	}

	var data RegisterRequest

	// Читаем json и записываем (декодируем) его в переменную data, иначе получаем ошибку в err
	err := json.NewDecoder(r.Body).Decode(&data) // &data - для того, чтобы записывать в саму переменную, а не делать копию
	if err != nil {
		http.Error(w, "Ошибка чтения json", http.StatusBadRequest)
		return
	}

	// Проверяем что логин не пустой и длина пароля не менее 8 символов
	if data.Login == "" || len(data.Password) < 8 {
		http.Error(w, "Логин - обязателен, пароль должен быть не менее 8 символов", http.StatusBadRequest)
		return
	}

	// Хэширование пароля с помощью bcrypt, хранение и дальнейшее использование пароля в исходном после его передачи - ЗАПРЕЩЕНО
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost) // bcrypt.GenerateFromPassword - функция преобразования пароля в хэш (работает только с byte, поэтому надо сначала пароль преобразовать в byte, а bcrypt.DefaultCost - грубо говоря значит то, насколько сложным будет хэш, и насколько долго он будет вычисляться)
	if err != nil {
		http.Error(w, "Ошибка обработки пароля", http.StatusInternalServerError)
	}
	fmt.Println(string(hash)) // ВРЕМЕНО, ЧТОБЫ НЕ БЫЛО ОШИБКИ У hash
}
