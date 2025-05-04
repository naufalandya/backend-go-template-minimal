# ✨ Clean Code Workflow for Fiber Controller - `CreateUser` ✨

Hallo minna-san! (｡>﹏<｡)  
Di sini kita mau bahas gimana cara membuat controller yang clean, kawaii, dan maintainable menggunakan contoh fungsi `CreateUser` di Fiber framework.  
Yuk kita lihat alur kerjanya! (っ˘ω˘ς )

---

## 🌸 Workflow Controller: `CreateUser`

### 1. Parsing Request Body

Controller selalu dimulai dengan mencoba **membaca isi request body** (`c.BodyParser`) dan mengubahnya menjadi struct input (`models.UserInput`).

```go
if err := c.BodyParser(&input); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
        Code:    fiber.StatusBadRequest,
        Status:  false,
        Message: fmt.Sprintf("Cannot parse your body~ (´；ω；｀) : %s", err.Error()),
    })
}
```

Supaya data dari client bisa langsung diproses sebagai objek Go yang rapi dan aman dipakai dalam program kita~ (⁄ ⁄>⁄ ▽ ⁄<⁄ ⁄)

### 2.  Basic Input Sanitization

Setelah body berhasil diparsing, kita lakukan sanitasi untuk mendeteksi input berbahaya menggunakan function FuckOffHackerByJSON.

```go
if err := functions.FuckOffHackerByJSON(&input); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
        Code:    fiber.StatusBadRequest,
        Status:  false,
        Message: fmt.Sprintf("Suspicious input detected~ (｀_´) : %s", err.Error()),
    })
}
```

Agar bisa mencegah serangan-serangan berbahaya kayak XSS atau SQL Injection yang bisa bersembunyi dalam JSON request! (｀・ω・´)

### 3. Validation

Validasi semua field di struct agar memenuhi aturan yang sudah ditentukan.

```go
if errs := functions.ValidateStruct(input); errs != nil {
    return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
        Code:    fiber.StatusBadRequest,
        Status:  false,
        Message: fmt.Sprintf("Validation failed~ (｀_´) : %s", errs[0]),
    })
}
```

Karena kita perlu memastikan semua data yang masuk itu valid, lengkap, dan sesuai format sebelum diproses lebih lanjut~ (´｡• ᵕ •｡`)

### 4. Main Logic (Business Logic)

Kalau semuanya sudah aman, baru jalankan main logic di sini.

Important Rule:

Jangan tuliskan detail teknis seperti query database, kalkulasi berat, atau pengiriman email langsung di dalam controller.

Buat function terpisah (service/helper) agar controller tetap clean, fokus pada alur utama aja.

```go
// Do Something With Data
// Always Clean Code!
// Always Put Main Logic Inside Controller, Don't Detail It!
```

### 5. Return Response

Terakhir, kita kembalikan response ke client dengan format yang konsisten.

```go
return c.Status(fiber.StatusCreated).JSON(global.Apiresponse{
    Code:    fiber.StatusCreated,
    Status:  true,
    Message: "Success",
    Data:    input,
})
```

Agar semua response dari API kita mudah dipahami, seragam, dan predictable buat front-end atau pengguna lain. (´｡• ω •｡`)

### 🐾 Mini Summary Workflow

```go
Parse Request Body → Sanitize Input → Validate Input → Execute Main Logic → Return API Response
```

### Example Directory Structure Suggestion

Supaya project tetap rapi dan maintainable, berikut contoh struktur folder yang bisa digunakan:

```go
/controllers
    - user_controller.go
/services
    - user_service.go
/models
    - user_input.go
/helpers
    - validation.go
    - sanitizer.go
/global
    - api_response.go
```