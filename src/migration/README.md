# Dynamic Migration System

Sistem migration dinamis yang dapat menampung beberapa model sekaligus.

## Cara Kerja

Sistem ini menggunakan **Registry Pattern** untuk mendaftarkan semua model yang perlu di-migrate. Semua model yang terdaftar akan otomatis di-migrate saat aplikasi start.

## Menambahkan Model Baru

### 1. Buat Model di `src/model/`

```go
// src/model/order.go
package model

import (
	"time"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Total     float64        `json:"total" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:pending"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Order) TableName() string {
	return "orders"
}
```

### 2. Daftarkan Model di `src/migration/create_users_table.go`

```go
// src/migration/create_users_table.go
func init() {
	RegisterModel(&model.User{})
	RegisterModel(&model.Product{})
	RegisterModel(&model.Order{}) // Tambahkan model baru di sini
}
```

Selesai! Model akan otomatis di-migrate saat aplikasi start.

## Fungsi yang Tersedia

### `RunMigrations(db *gorm.DB) error`
Menjalankan migration untuk semua model yang terdaftar.

```go
if err := migration.RunMigrations(db); err != nil {
    log.Fatal("Failed to migrate:", err)
}
```

### `RegisterModel(model interface{})`
Mendaftarkan model baru untuk migration (biasanya dipanggil di `init()`).

```go
migration.RegisterModel(&model.User{})
```

### `GetRegisteredModels() []string`
Mengembalikan daftar nama model yang terdaftar.

```go
models := migration.GetRegisteredModels()
// Output: ["User", "Product", "Order"]
```

## Contoh Penggunaan

### Di main.go

```go
import "boilerblade/src/migration"

func main() {
	app, _ := server.NewApp()
	
	// Run migrations
	if app.Config.Database != nil {
		if err := migration.RunMigrations(app.Config.Database); err != nil {
			log.Fatal("Failed to migrate:", err)
		}
		log.Printf("Migrated models: %v", migration.GetRegisteredModels())
	}
}
```

## Logging

Sistem migration akan log:
- Jumlah model yang akan di-migrate
- Setiap model yang sedang di-migrate
- Status sukses/gagal untuk setiap model
- Summary setelah semua migration selesai

## Keuntungan

1. **Otomatis**: Semua model terdaftar otomatis di-migrate
2. **Fleksibel**: Mudah menambahkan model baru
3. **Terpusat**: Semua model terdaftar di satu tempat
4. **Logging**: Detail logging untuk debugging
5. **Error Handling**: Error handling yang baik untuk setiap model

## Model yang Sudah Terdaftar

- `User` - User entity
- `Product` - Product entity
- (Tambahkan model lain sesuai kebutuhan)
