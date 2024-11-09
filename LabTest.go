package main 
 
import ( 
 "context" 
 "fmt" 
 "io" 
 "math/rand" 
 "os" 
 "time" 
) 
 
func main() { 
 // Начало отсчета времени 
 start := time.Now() 
 
 // Имя файла и его размер в MB 
 fileName := "r.bin" 
 sizeMB := 30 
 sizeBytes := sizeMB * 1024 * 1024 
 
 
//fileName = "/ergerhgsrthsrth/r.bin" // Устанавливаем недопустимый путь 
file, err := os.Create(fileName) 
if err != nil { 
    fmt.Println("Ошибка создания файла:", err) 
    return 
} 
defer file.Close() 
 
 
 // Запись случайных данных в файл 
 data := make([]byte, sizeBytes) 
 rand.Read(data) 
//file.Close() // Закрытие файла для проверки ошибки записи 
 _, err = file.Write(data) 
 if err != nil { 
  fmt.Println("Ошибка записи в файл:", err) 
  return 
 } 
 
// Удаление файла перед его открытием для вызова ошибки 
//os.Remove(fileName) 
 
// Открытие файла для чтения 
file, err = os.Open(fileName) 
if err != nil { 
    fmt.Println("Ошибка открытия файла:", err) 
    return 
} 
defer file.Close() 
 
 
 // Создание контекста с таймаутом 2 секунды 
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond) 
 defer cancel() 
 
 readBytes := 0 
 buffer := make([]byte, 1024) 
 
 // Канал для чтения в отдельной горутине 
 done := make(chan error) 
 
 go func() { 
  for { 
   n, err := file.Read(buffer) 
   if err == io.EOF { 
    done <- nil // Завершение чтения 
    return 
   } 
   if err != nil { 
    done <- err // Ошибка при чтении 
    return 
   } 
   readBytes += n 
  } 
 }() 
 
select { 
case <-ctx.Done(): 
    fmt.Println("Время вышло:", ctx.Err()) 
    fmt.Printf("Файл не был прочитан полностью, прочитано байт: %d из %d\n", readBytes, sizeBytes) 
case err := <-done: 
    if err != nil { 
        fmt.Println("Ошибка чтения файла:", err) 
    } else { 
       fmt.Printf("Файл прочитан полностью, прочитано байт: %d из %d\n", readBytes, sizeBytes) 
    } 
} 
 
 // Вывод результатов 
 fmt.Println("Прочитано байт:", readBytes) 
 
 // Подсчет времени выполнения 
 elapsed := time.Since(start) 
 fmt.Printf("Время выполнения: %.2f сек.\n", elapsed.Seconds()) 
}