package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Celebrity struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"not null"`
	Age        int    `json:"age" gorm:"not null"`
	Country    string `json:"country" gorm:"not null"`
	Profession string `json:"profession" gorm:"not null"`
}

var db *gorm.DB

func getDatabaseURL() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	return "postgres://postgres:postgres@localhost:5432/celebrities?sslmode=disable"
}

func initDB() {
	var err error
	db, err = gorm.Open(postgres.Open(getDatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	if err = db.AutoMigrate(&Celebrity{}); err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	log.Println("PostgreSQL (GORM) инициализирован успешно")
}

func getAllCelebrities(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /Celebrities/All")

	var celebrities []Celebrity
	if err := db.Find(&celebrities).Error; err != nil {
		log.Printf("Ошибка запроса: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(celebrities)
	log.Printf("Возвращено %d записей", len(celebrities))
}

func getCelebrityByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}
	log.Printf("GET /Celebrities/%d", id)

	var c Celebrity
	if err := db.First(&c, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Элемент с Id=%d не найден", id)
			http.Error(w, "Элемент не найден", http.StatusNotFound)
			return
		}

		log.Printf("Ошибка запроса: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c)
}

func addCelebrity(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /Celebrities")

	var c Celebrity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	var exists int64
	if err := db.Model(&Celebrity{}).Where("id = ?", c.Id).Count(&exists).Error; err != nil {
		log.Printf("Ошибка проверки Id: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		log.Printf("Элемент с Id=%d уже существует", c.Id)
		http.Error(w, "Элемент с таким Id уже существует", http.StatusConflict)
		return
	}

	if err := db.Create(&c).Error; err != nil {
		log.Printf("Ошибка вставки: %v", err)
		http.Error(w, "Ошибка записи в БД", http.StatusInternalServerError)
		return
	}

	log.Printf("Добавлен элемент Id=%d, Name=%s", c.Id, c.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

func updateCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}
	log.Printf("PUT /Celebrities/%d", id)

	var existing Celebrity
	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Элемент с Id=%d не найден для обновления", id)
			http.Error(w, "Элемент не найден", http.StatusNotFound)
			return
		}

		log.Printf("Ошибка проверки Id: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}

	var c Celebrity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	c.Id = id

	if err := db.Model(&existing).Updates(map[string]interface{}{
		"name":       c.Name,
		"age":        c.Age,
		"country":    c.Country,
		"profession": c.Profession,
	}).Error; err != nil {
		log.Printf("Ошибка обновления: %v", err)
		http.Error(w, "Ошибка обновления в БД", http.StatusInternalServerError)
		return
	}

	log.Printf("Обновлён элемент Id=%d", id)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c)
}

func deleteCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}
	log.Printf("DELETE /Celebrities/%d", id)

	result := db.Delete(&Celebrity{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("Ошибка удаления: %v", result.Error)
		http.Error(w, "Ошибка удаления из БД", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("Элемент с Id=%d не найден для удаления", id)
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}

	log.Printf("Удалён элемент Id=%d", id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	log.Println("Запуск сервера GO06_01...")
	initDB()

	r := mux.NewRouter()

	r.HandleFunc("/Celebrities/All", getAllCelebrities).Methods("GET")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", getCelebrityByID).Methods("GET")
	r.HandleFunc("/Celebrities", addCelebrity).Methods("POST")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", updateCelebrity).Methods("PUT")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", deleteCelebrity).Methods("DELETE")

	log.Println("Сервер запущен на порту 3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
