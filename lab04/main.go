package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Celebrity struct {
	Id         int    `json:"Id"`
	Name       string `json:"Name"`
	Profession string `json:"Profession"`
	Country    string `json:"Country"`
	BirthYear  int    `json:"BirthYear"`
}

const dataFile = "Celebrities.json"

func loadCelebrities() ([]Celebrity, error) {
	log.Printf("[loadCelebrities] Чтение файла: %s", dataFile)
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var celebrities []Celebrity
	if err := json.Unmarshal(data, &celebrities); err != nil {
		return nil, err
	}
	log.Printf("[loadCelebrities] Загружено %d записей", len(celebrities))
	return celebrities, nil
}

func saveCelebrities(celebrities []Celebrity) error {
	log.Printf("[saveCelebrities] Сохранение %d записей в файл: %s", len(celebrities), dataFile)
	data, err := json.MarshalIndent(celebrities, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func getAllCelebrities(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET] /Celebrities/All — запрос от %s", r.RemoteAddr)
	celebrities, err := loadCelebrities()
	if err != nil {
		log.Printf("[ERROR] getAllCelebrities: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(celebrities)
	log.Printf("[GET] /Celebrities/All — отправлено %d записей", len(celebrities))
}

func getCelebrityById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Printf("[GET] /Celebrities/%s — запрос от %s", idStr, r.RemoteAddr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ERROR] getCelebrityById: некорректный id=%s", idStr)
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}

	celebrities, err := loadCelebrities()
	if err != nil {
		log.Printf("[ERROR] getCelebrityById: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	for _, c := range celebrities {
		if c.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			log.Printf("[GET] /Celebrities/%d — найден: %s", id, c.Name)
			return
		}
	}

	log.Printf("[GET] /Celebrities/%d — не найден", id)
	http.Error(w, "Элемент не найден", http.StatusNotFound)
}

func addCelebrity(w http.ResponseWriter, r *http.Request) {
	log.Printf("[POST] /Celebrities — запрос от %s", r.RemoteAddr)

	var newCelebrity Celebrity
	if err := json.NewDecoder(r.Body).Decode(&newCelebrity); err != nil {
		log.Printf("[ERROR] addCelebrity: ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	celebrities, err := loadCelebrities()
	if err != nil {
		log.Printf("[ERROR] addCelebrity: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	for _, c := range celebrities {
		if c.Id == newCelebrity.Id {
			log.Printf("[POST] /Celebrities — дублирование Id=%d (409 Conflict)", newCelebrity.Id)
			http.Error(w, "Элемент с таким Id уже существует", http.StatusConflict)
			return
		}
	}

	celebrities = append(celebrities, newCelebrity)
	if err := saveCelebrities(celebrities); err != nil {
		log.Printf("[ERROR] addCelebrity: ошибка сохранения: %v", err)
		http.Error(w, "Ошибка сохранения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCelebrity)
	log.Printf("[POST] /Celebrities — добавлен: Id=%d, Name=%s", newCelebrity.Id, newCelebrity.Name)
}

func updateCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Printf("[PUT] /Celebrities/%s — запрос от %s", idStr, r.RemoteAddr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ERROR] updateCelebrity: некорректный id=%s", idStr)
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}

	var updatedCelebrity Celebrity
	if err := json.NewDecoder(r.Body).Decode(&updatedCelebrity); err != nil {
		log.Printf("[ERROR] updateCelebrity: ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	celebrities, err := loadCelebrities()
	if err != nil {
		log.Printf("[ERROR] updateCelebrity: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	found := false
	for i, c := range celebrities {
		if c.Id == id {
			updatedCelebrity.Id = id
			celebrities[i] = updatedCelebrity
			found = true
			log.Printf("[PUT] /Celebrities/%d — обновлён: %s", id, updatedCelebrity.Name)
			break
		}
	}

	if !found {
		log.Printf("[PUT] /Celebrities/%d — не найден (404)", id)
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}

	if err := saveCelebrities(celebrities); err != nil {
		log.Printf("[ERROR] updateCelebrity: ошибка сохранения: %v", err)
		http.Error(w, "Ошибка сохранения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCelebrity)
}

func deleteCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Printf("[DELETE] /Celebrities/%s — запрос от %s", idStr, r.RemoteAddr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ERROR] deleteCelebrity: некорректный id=%s", idStr)
		http.Error(w, "Некорректный Id", http.StatusBadRequest)
		return
	}

	celebrities, err := loadCelebrities()
	if err != nil {
		log.Printf("[ERROR] deleteCelebrity: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	found := false
	newList := make([]Celebrity, 0, len(celebrities))
	for _, c := range celebrities {
		if c.Id == id {
			found = true
			log.Printf("[DELETE] /Celebrities/%d — удалён: %s", id, c.Name)
		} else {
			newList = append(newList, c)
		}
	}

	if !found {
		log.Printf("[DELETE] /Celebrities/%d — не найден (404)", id)
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}

	if err := saveCelebrities(newList); err != nil {
		log.Printf("[ERROR] deleteCelebrity: ошибка сохранения: %v", err)
		http.Error(w, "Ошибка сохранения данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	log.Println("=== Запуск GO04_01 Web-сервера ===")

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		log.Fatalf("[FATAL] Файл данных не найден: %s", dataFile)
	}

	r := mux.NewRouter()

	r.HandleFunc("/Celebrities/All", getAllCelebrities).Methods(http.MethodGet)
	r.HandleFunc("/Celebrities/{id:[0-9]+}", getCelebrityById).Methods(http.MethodGet)
	r.HandleFunc("/Celebrities", addCelebrity).Methods(http.MethodPost)
	r.HandleFunc("/Celebrities/{id:[0-9]+}", updateCelebrity).Methods(http.MethodPut)
	r.HandleFunc("/Celebrities/{id:[0-9]+}", deleteCelebrity).Methods(http.MethodDelete)

	port := ":3000"
	log.Printf("Сервер слушает на порту %s", port)
	log.Printf("Маршруты:")
	log.Printf("  GET    /Celebrities/All")
	log.Printf("  GET    /Celebrities/{id}")
	log.Printf("  POST   /Celebrities")
	log.Printf("  PUT    /Celebrities/{id}")
	log.Printf("  DELETE /Celebrities/{id}")

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("[FATAL] Ошибка запуска сервера: %v", err)
	}
}
