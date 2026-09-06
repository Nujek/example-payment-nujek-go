# Contoh project Go memakai SDK merchant

Project CLI ini menggunakan `github.com/Nujek/sdk-payment-nujek-go`.

## Menjalankan

Salin template konfigurasi lalu isi kredensial merchant:

```bash
cp .env.example .env
```

Ganti `your-api-key` dan `your-api-secret` dengan credential merchant yang sebenarnya. Placeholder tersebut akan ditolak sebelum request dikirim.

Project otomatis membaca `.env`. Environment variable yang sudah diset dari shell tetap memiliki prioritas.

```bash
go run . balance
go run . create-bill -external-id order-123 -channel NOBU_QRIS -total 150000.00
go run . bill <bill-uuid>
go run . list-bills -page 1 -per-page 20
go run . create-payout -external-id payout-123 -bank BNI -account 1234567890 -name "Nama Merchant" -amount 50000.00
go run . qris-static
go run . qris-static <qris-uuid>
```

SDK otomatis membuat signature HMAC dan tiga header autentikasi pada setiap request.
