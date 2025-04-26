# ✨ Modular Monolith Template 🔥

Halo halo~ (≧◡≦) ♡  
Ini adalah panduan resmi dan kawaii untuk menggunakan Modular Monolith Template dengan bersih, cepat, dan rapi!

---

## 📚 Workflow Step-by-Step

1. **Gunakan Hanya Module yang Dibutuhkan**  
   Tempatkan module kamu di folder:

/module

> Hanya masukkan module yang benar-benar diperlukan agar project tetap ringan dan clean~ (๑˃ᴗ˂)ﻭ

1. **Konfigurasi Credential Secara Mandiri**  
Atur semua credential seperti database, API keys, dan lainnya di:

/server/config

> Setiap environment (dev, staging, prod) bisa punya config sendiri agar lebih aman dan fleksibel~ ✨

1. **Tulis Definisi Struct dengan Cepat**  
Buat definisi struct-mu menggunakan syntax Proto dengan mudah di:

/mock

> Cuma perlu file `.txt` sederhana dengan format proto! (灬º‿º灬)

1. **Generate Struct Otomatis dengan Validator**  
Gunakan tools super ajaib kita:

/scripts/generate_struct.exe


> File `.exe` ini akan mengubah syntax Proto kamu menjadi Go struct lengkap dengan validasi ready to use~ (ﾉ◕ヮ◕)ﾉ*:･ﾟ✧

1. **Pindahkan Struct yang Dihasilkan ke Modul Kamu**  
Setelah generate, pindahkan hasilnya ke module tujuan. Misalnya:

/module/user/model


> Jangan lupa diorganize supaya project makin rapi dan kawaii~ ( ˘ ³˘)♥

1. **Register Hanya Module yang Dipakai ke Gateway**  
Daftarkan module-mu yang aktif ke dalam router utama:

/server/api/gateway.go


> Modul yang tidak didaftarkan = tidak aktif, sehingga performa server tetap optimal~ 💨

1. **Jalankan Atau Build Project**  
Untuk running atau build gampang banget, cukup:

```go
make run

atau

make build
```

> Biar codinganmu langsung terbang tanpa drama~ (づ｡◕‿‿◕｡)づ

---

## ⚡ Tentang `generate_struct.exe`

- Lokasi tools:
/scripts/generate_struct.exe


- Fungsi:  
Mengubah file definisi proto di `/mock/*.txt` menjadi Go struct yang sudah siap pakai dan sudah ada tag validator-nya.
- Output:  
File Go baru yang akan kamu pindahkan ke:
/module/[nama-modul]/model


- Contoh flow kerja:
1. Tulis proto di `/mock/user.txt`
2. Jalankan `generate_struct.exe`
3. Ambil hasil generate di `/model`
4. Pindahkan ke `/module/user/model`

✨ Mudah kan? Let's keep our codebase super clean and super kawaii together~!! (´｡• ᵕ •｡`) ♡