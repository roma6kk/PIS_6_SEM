package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Celebrity struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Country    string `json:"country"`
	Profession string `json:"profession"`
}

var db *sql.DB

func getDatabaseURL() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	return "postgres://postgres:postgres@localhost:5432/celebrities?sslmode=disable"
}

// initDB - инициализация базы данных и создание таблицы
func initDB() {
	var err error
	db, err = sql.Open("pgx", getDatabaseURL())
	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS Celebrities (
		Id         INTEGER PRIMARY KEY,
		Name       TEXT    NOT NULL,
		Age        INTEGER NOT NULL,
		Country    TEXT    NOT NULL,
		Profession TEXT    NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	log.Println("PostgreSQL инициализирован успешно")
}

func getAllCelebrities(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /Celebrities/All")

	rows, err := db.Query("SELECT Id, Name, Age, Country, Profession FROM Celebrities")
	if err != nil {
		log.Printf("Ошибка запроса: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	celebrities := []Celebrity{}
	for rows.Next() {
		var c Celebrity
		if err := rows.Scan(&c.Id, &c.Name, &c.Age, &c.Country, &c.Profession); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		celebrities = append(celebrities, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(celebrities)
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
	err = db.QueryRow("SELECT Id, Name, Age, Country, Profession FROM Celebrities WHERE Id = $1", id).
		Scan(&c.Id, &c.Name, &c.Age, &c.Country, &c.Profession)

	if err == sql.ErrNoRows {
		log.Printf("Элемент с Id=%d не найден", id)
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Ошибка запроса: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func addCelebrity(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /Celebrities")

	var c Celebrity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM Celebrities WHERE Id = $1", c.Id).Scan(&exists)
	if err != nil {
		log.Printf("Ошибка проверки Id: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		log.Printf("Элемент с Id=%d уже существует", c.Id)
		http.Error(w, "Элемент с таким Id уже существует", http.StatusConflict) // 409
		return
	}

	_, err = db.Exec("INSERT INTO Celebrities (Id, Name, Age, Country, Profession) VALUES ($1, $2, $3, $4, $5)",
		c.Id, c.Name, c.Age, c.Country, c.Profession)
	if err != nil {
		log.Printf("Ошибка вставки: %v", err)
		http.Error(w, "Ошибка записи в БД", http.StatusInternalServerError)
		return
	}

	log.Printf("Добавлен элемент Id=%d, Name=%s", c.Id, c.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func updateCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}
	log.Printf("PUT /Celebrities/%d", id)

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM Celebrities WHERE Id = $1", id).Scan(&exists)
	if err != nil {
		log.Printf("Ошибка проверки Id: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		log.Printf("Элемент с Id=%d не найден для обновления", id)
		http.Error(w, "Элемент не найден", http.StatusNotFound) // 404
		return
	}

	var c Celebrity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	c.Id = id

	_, err = db.Exec("UPDATE Celebrities SET Name=$1, Age=$2, Country=$3, Profession=$4 WHERE Id=$5",
		c.Name, c.Age, c.Country, c.Profession, id)
	if err != nil {
		log.Printf("Ошибка обновления: %v", err)
		http.Error(w, "Ошибка обновления в БД", http.StatusInternalServerError)
		return
	}

	log.Printf("Обновлён элемент Id=%d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func deleteCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}
	log.Printf("DELETE /Celebrities/%d", id)

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM Celebrities WHERE Id = $1", id).Scan(&exists)
	if err != nil {
		log.Printf("Ошибка проверки Id: %v", err)
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		log.Printf("Элемент с Id=%d не найден для удаления", id)
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM Celebrities WHERE Id = $1", id)
	if err != nil {
		log.Printf("Ошибка удаления: %v", err)
		http.Error(w, "Ошибка удаления из БД", http.StatusInternalServerError)
		return
	}

	log.Printf("Удалён элемент Id=%d", id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	log.Println("Запуск сервера GO05_01...")

	initDB()
	defer db.Close()

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
