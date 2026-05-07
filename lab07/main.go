package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	precision = 2 // Точность по умолчанию
	precMtx   sync.RWMutex
)

// Структуры JSON-RPC 2.0
type RPCReq struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id,omitempty"`
}

type RPCRes struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Result struct {
	Value float64 `json:"Value"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rpc", handleRPC).Methods("POST")

	log.Println("🚀 GO07_01 RPC server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

// Обработчик маршрута /rpc
func handleRPC(w http.ResponseWriter, r *http.Request) {
	var raw json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		log.Printf("❌ Parse error: %v", err)
		writeJSON(w, http.StatusBadRequest, RPCRes{JSONRPC: "2.0", Error: &RPCError{Code: -32700, Message: "Parse error"}})
		return
	}

	if len(raw) == 0 {
		writeJSON(w, http.StatusBadRequest, RPCRes{JSONRPC: "2.0", Error: &RPCError{Code: -32600, Message: "Invalid Request"}})
		return
	}

	// Определение: пакетный запрос начинается с '['
	if raw[0] == '[' {
		var batch []RPCReq
		if err := json.Unmarshal(raw, &batch); err != nil {
			writeJSON(w, http.StatusBadRequest, RPCRes{JSONRPC: "2.0", Error: &RPCError{Code: -32700, Message: "Parse error"}})
			return
		}
		handleBatch(w, batch)
	} else {
		var req RPCReq
		if err := json.Unmarshal(raw, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, RPCRes{JSONRPC: "2.0", Error: &RPCError{Code: -32700, Message: "Parse error"}})
			return
		}
		log.Printf("📥 Received single request: method=%s, id=%v", req.Method, req.ID)
		res := processRequest(&req)
		if res != nil {
			writeJSON(w, http.StatusOK, *res)
		} else {
			// Определение: уведомление не требует ответа
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// Обработка пакетного запроса
func handleBatch(w http.ResponseWriter, batch []RPCReq) {
	var responses []RPCRes
	for i := range batch {
		res := processRequest(&batch[i])
		if res != nil {
			responses = append(responses, *res)
		}
	}
	if len(responses) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, responses)
}

// Логика выполнения метода
func processRequest(req *RPCReq) *RPCRes {
	res := &RPCRes{JSONRPC: "2.0", ID: req.ID}

	if req.JSONRPC != "2.0" {
		res.Error = &RPCError{Code: -32600, Message: "Invalid Request"}
		return res
	}

	switch req.Method {
	case "sum", "sub", "mul", "div":
		x, y, err := extractXY(req.Params)
		if err != nil {
			res.Error = &RPCError{Code: -32602, Message: "Invalid params"}
			return res
		}

		var val float64
		switch req.Method {
		case "sum":
			val = x + y
		case "sub":
			val = x - y
		case "mul":
			val = x * y // В задании ошибка: указано x-y, по логике умножения используется *
		case "div":
			if y == 0 {
				res.Error = &RPCError{Code: -32603, Message: "Division by zero"}
				return res
			}
			val = x / y
		}

		precMtx.RLock()
		p := precision
		precMtx.RUnlock()

		res.Result = Result{Value: roundToPrecision(val, p)}

	case "pre":
		n, err := extractN(req.Params)
		if err != nil {
			res.Error = &RPCError{Code: -32602, Message: "Invalid params"}
			return res
		}
		precMtx.Lock()
		precision = n
		precMtx.Unlock()
		log.Printf("⚙️ Precision changed to %d", n)
		res.Result = "ok"

	default:
		res.Error = &RPCError{Code: -32601, Message: "Method not found"}
		return res
	}

	// Если id отсутствует или равен null -> это уведомление, ответ не отправляем
	if req.ID == nil {
		return nil
	}
	return res
}

// Парсинг параметров x, y (форматы 1 и 2)
func extractXY(params interface{}) (float64, float64, error) {
	switch p := params.(type) {
	case []interface{}:
		if len(p) != 2 {
			return 0, 0, fmt.Errorf("array must have exactly 2 elements")
		}
		x, ok1 := p[0].(float64)
		y, ok2 := p[1].(float64)
		if !ok1 || !ok2 {
			return 0, 0, fmt.Errorf("elements must be valid numbers")
		}
		return x, y, nil
	case map[string]interface{}:
		x, ok1 := p["x"].(float64)
		y, ok2 := p["y"].(float64)
		if !ok1 || !ok2 {
			return 0, 0, fmt.Errorf("x and y must be valid numbers")
		}
		return x, y, nil
	default:
		return 0, 0, fmt.Errorf("invalid params format")
	}
}

// Парсинг параметра N (формат 3)
func extractN(params interface{}) (int, error) {
	if m, ok := params.(map[string]interface{}); ok {
		if n, ok := m["N"].(float64); ok {
			return int(n), nil
		}
	}
	return 0, fmt.Errorf("N must be an integer")
}

// Округление до заданного количества знаков после запятой
func roundToPrecision(val float64, prec int) float64 {
	if prec < 0 {
		prec = 0
	}
	shift := math.Pow(10, float64(prec))
	return math.Round(val*shift) / shift
}

// Вспомогательная функция для отправки JSON
func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, _ := json.Marshal(v)
	w.Write(data)
}