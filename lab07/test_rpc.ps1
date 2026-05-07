# ============================================
# Тестовый скрипт для лабораторной работы ПИС-07
# Сервер: GO07_01 (JSON-RPC 2.0, порт 3000)
# Запуск: .\test_rpc.ps1
# ============================================

$BASE_URL = "http://localhost:3000/rpc"
$HEADERS = @{ "Content-Type" = "application/json" }

Write-Host "🚀 Запуск тестов для GO07_01 RPC сервера" -ForegroundColor Cyan
Write-Host ""

# --------------------------------------------
# БЛОК 1: Тесты с массивом параметров [x, y]
# --------------------------------------------
Write-Host "📦 БЛОК 1: Параметры в виде массива [x, y]" -ForegroundColor Green

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sum","params":[5.05, 3.33],"id":1}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sub","params":[5.05, 3.33],"id":2}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"mul","params":[5.05, 3.33],"id":3}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"div","params":[5.05, 3.33],"id":4}'
Write-Host ""
Write-Host ""

# --------------------------------------------
# БЛОК 2: Тесты с объектом параметров {"x": , "y": }
# --------------------------------------------
Write-Host "📦 БЛОК 2: Параметры в виде объекта {`"x`": , `"y`": }" -ForegroundColor Green

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sum","params":{"x":5.05, "y":3.33},"id":5}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sub","params":{"x":5.05, "y":3.33},"id":6}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"mul","params":{"x":5.05, "y":3.33},"id":7}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"div","params":{"x":5.05, "y":3.33},"id":8}'
Write-Host ""
Write-Host ""

# --------------------------------------------
# БЛОК 3: Уведомление pre (без id) - установка точности
# --------------------------------------------
Write-Host "📦 БЛОК 3: Уведомление pre (установка точности)" -ForegroundColor Green

Write-Host "Установка точности: 3 знака" -ForegroundColor Yellow
curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"pre","params":{"N":3}}'
Write-Host "`n(Уведомление: ответа нет, статус 204)" -ForegroundColor DarkGray
Write-Host ""

# --------------------------------------------
# БЛОК 4: Проверка точности = 3 (массив параметров)
# --------------------------------------------
Write-Host "📦 БЛОК 4: Проверка результатов с точностью 3 знака" -ForegroundColor Green

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sum","params":[5.05, 3.33],"id":9}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sub","params":[5.05, 3.33],"id":10}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"mul","params":[5.05, 3.33],"id":11}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"div","params":[5.05, 3.33],"id":12}'
Write-Host ""
Write-Host ""

# --------------------------------------------
# БЛОК 5: Смена точности на 1 знак
# --------------------------------------------
Write-Host "📦 БЛОК 5: Установка точности = 1 знак" -ForegroundColor Green

Write-Host "Установка точности: 1 знак" -ForegroundColor Yellow
curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"pre","params":{"N":1}}'
Write-Host "`n(Уведомление: ответа нет, статус 204)" -ForegroundColor DarkGray
Write-Host ""

# --------------------------------------------
# БЛОК 6: Проверка точности = 1 (объект параметров)
# --------------------------------------------
Write-Host "📦 БЛОК 6: Проверка результатов с точностью 1 знак" -ForegroundColor Green

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sum","params":{"x":5.05, "y":3.33},"id":13}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"sub","params":{"x":5.05, "y":3.33},"id":14}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"mul","params":{"x":5.05, "y":3.33},"id":15}'
Write-Host ""

curl.exe -X POST $BASE_URL -Headers $HEADERS -d '{"jsonrpc":"2.0","method":"div","params":{"x":5.05, "y":3.33},"id":16}'
Write-Host ""
Write-Host ""

# --------------------------------------------
# БЛОК 7: Пакетный запрос (Batch)
# --------------------------------------------
Write-Host "📦 БЛОК 7: Пакетная обработка (Batch Request)" -ForegroundColor Green

# Используем here-string для чистого JSON без экранирования
$batchJson = @'
[
  {"jsonrpc":"2.0","method":"pre","params":{"N":4}},
  {"jsonrpc":"2.0","method":"sum","params":{"x":5.05, "y":3.33},"id":17},
  {"jsonrpc":"2.0","method":"sub","params":{"x":5.05, "y":3.33},"id":18},
  {"jsonrpc":"2.0","method":"mul","params":{"x":5.05, "y":3.33},"id":19},
  {"jsonrpc":"2.0","method":"div","params":{"x":5.05, "y":3.33},"id":20}
]
'@

curl.exe -X POST $BASE_URL -Headers $HEADERS -d $batchJson
Write-Host ""
Write-Host ""

# --------------------------------------------
# Завершение
# --------------------------------------------
Write-Host "✅ Все тесты завершены!" -ForegroundColor Cyan