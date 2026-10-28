# Plan Implementasi Dokumentasi API Backend Kiroku

## 0. Metadata

- **Tujuan plan:** menjadi instruksi implementasi lengkap untuk membangun dokumentasi API backend Kiroku yang akurat, dapat dibuka melalui Swagger UI, dan mencerminkan route yang benar-benar aktif.
- **Target utama:** mengganti `backend/docs/docs.go` yang saat ini memiliki `paths: {}` dengan spesifikasi Swagger 2.0 lengkap.
- **Scope utama:** route yang benar-benar dimount oleh `backend/cmd/server/main.go`, handler yang dipanggil route tersebut, domain/model request-response, autentikasi cookie, status HTTP, error body, parameter query/path, dan verifikasi build/runtime.
- **Scope yang tidak termasuk:** refactor handler, perubahan perilaku API, penambahan endpoint baru, pemasangan middleware auth yang belum ada, pembuatan endpoint dari package internal yang belum dimount, perubahan frontend, perubahan database/migration, perubahan dependency kecuali benar-benar diperlukan.
- **File utama yang kemungkinan berubah:** `backend/docs/docs.go`.
- **File yang dibaca sebagai sumber kebenaran:** `backend/cmd/server/main.go`, `backend/internal/delivery/http/*.go`, `backend/internal/auth/*.go`, `backend/internal/domain/*.go`, `backend/internal/models/*.go`, package handler lama yang relevan hanya untuk menemukan gap.
- **Output akhir:** dokumentasi Swagger yang dapat dirender oleh `GET /api/docs/` dan JSON spec yang tidak lagi kosong.

---

## 1. Prinsip dan batasan wajib

1. Dokumentasikan route live berdasarkan wiring di `backend/cmd/server/main.go`, bukan berdasarkan keberadaan file semata.
2. Bedakan route live dari route orphan/unmounted.
3. Jangan mengklaim endpoint internal sebagai tersedia apabila belum dimount.
4. Jangan mengubah kontrak API hanya demi membuat dokumentasi lebih rapi.
5. Semua schema harus mengikuti JSON tag aktual. Jika input struct tidak memiliki JSON tag, verifikasi cara decoder memetakan field sebelum menulis nama properti.
6. Jangan menulis secret, token runtime, PIN, credential, tunnel URL, isi `.env`, database, atau data user ke dokumentasi.
7. Jangan mendokumentasikan nilai token contoh sebagai token nyata.
8. Jangan menambahkan komentar kode baru kecuali memang diminta. Untuk `docs.go`, gunakan struktur data/spec yang sudah sesuai generator.
9. Jangan menambahkan endpoint analytics, receipts, import/export, backup, atau Telegram webhook jika route tersebut tidak dimount oleh `main.go`.
10. Jangan menyatakan seluruh `/api/v1` terlindungi auth jika source saat ini tidak memasang `AuthMiddleware` pada semua group.
11. Security metadata Swagger harus menggambarkan intent dan kondisi aktual. Jika route secara desain membutuhkan login tetapi wiring belum memasang middleware, catat gap tersebut di plan/hasil akhir dan jangan menyamarkan fakta source.
12. Jalankan verifikasi dari direktori yang benar: backend.
13. Perubahan harus kecil dan terarah; default hanya `backend/docs/docs.go`, kecuali verifikasi menunjukkan generator/registrasi memang perlu perubahan.

---

## 2. Kondisi awal yang sudah dikonfirmasi

### 2.1 Swagger saat ini

`backend/docs/docs.go`:

- package `docs`.
- memakai `github.com/swaggo/swag`.
- `SwaggerInfo` berisi:
  - version `1.0`;
  - host `localhost:8080`;
  - base path `/api/v1`;
  - title `Kiroku API`;
  - description `Personal Finance Tracker API`;
  - security definition `CookieAuth` dengan api key cookie `auth_token`.
- template saat ini memiliki `paths: {}`.
- file menyatakan generated, tetapi tidak boleh mengandalkan generator jika source annotations belum lengkap. Implementasi harus memastikan hasil spec terdaftar dan build tetap berhasil.

### 2.2 Wiring live di `main.go`

Route live yang benar-benar dimount:

- `GET /healthz`.
- `GET /api/docs/*` melalui `httpSwagger.Handler`.
- `/shared` melalui `SharedHandler.Routes()`.
- `/api/v1/auth` melalui auth handler.
- `/api/v1/accounts` melalui delivery accounts handler.
- `/api/v1/categories` melalui delivery categories handler.
- `/api/v1/transactions` melalui delivery transactions handler.
- `/api/v1/templates` melalui delivery templates handler.
- `/api/v1/scheduled` melalui delivery scheduled handler.
- `/api/v1/reports` melalui delivery reports handler.
- `/api/v1/settings` melalui delivery settings handler.
- `/api/v1/fx` melalui delivery FX handler.

Middleware global yang aktif:

- request ID;
- real IP wrapper yang saat ini hanya meneruskan request;
- CORS middleware;
- recoverer;
- logging middleware.

Auth middleware saat ini terlihat di auth handler untuk route auth tertentu. `main.go` tidak menunjukkan `AuthMiddleware` dipasang pada mounted accounts/categories/transactions/templates/scheduled/reports/settings/fx routes. Ini adalah known gap keamanan dari `AGENTS.md`; dokumentasi tidak boleh menghapus atau menyembunyikan gap itu.

### 2.3 Route internal yang tidak boleh otomatis didokumentasikan

Pencarian seluruh backend menemukan handler/package lain, termasuk:

- analytics;
- receipts;
- import/export CSV;
- backup;
- labels;
- package lama `internal/accounts`, `internal/transactions`, `internal/settings`, `internal/reports`, `internal/templates`, `internal/scheduled`, `internal/shared`;
- Telegram webhook.

Route-route tersebut tidak otomatis live hanya karena memiliki `Routes()`. Sebelum masuk Swagger, harus dibuktikan dimount di `main.go`. Jika tidak, masukkan ke bagian gap/out-of-scope, bukan ke `paths`.

---

## 3. Hasil yang harus dicapai

### 3.1 Hasil fungsional

- Swagger UI di `/api/docs/` berhasil memuat JSON spec.
- JSON spec memiliki `swagger: "2.0"`.
- `paths` berisi semua route live yang disepakati.
- Tidak ada path duplicate atau parameter path yang tidak didefinisikan.
- Semua operation memiliki `summary`, `description`, `tags`, parameter yang relevan, dan response.
- Semua operation memakai schema reusable jika response/request berulang.
- Cookie auth tersedia sebagai `securityDefinitions.CookieAuth`.
- Route public tidak diberi security requirement.
- Route auth terbuka diberi security sesuai perilaku aktual.
- Body error terdokumentasi sebagai `text/plain` apabila handler memakai `http.Error`, bukan dipaksa menjadi JSON tanpa perubahan handler.
- Contoh response tidak berisi secret nyata atau data runtime.

### 3.2 Hasil kualitas

- Nama properti JSON sama dengan output aktual.
- Tipe numerik, boolean, string, array, object, nullable/optional ditulis konsisten.
- Format waktu didokumentasikan sebagai `string` format `date-time` jika memakai `time.Time`/RFC3339.
- Format tanggal non-waktu didokumentasikan sebagai `string` format `date` jika handler memakai `YYYY-MM-DD`.
- Enum diturunkan dari konstanta yang benar-benar ada; jangan mengarang enum yang tidak divalidasi handler.
- Parameter default yang nyata, misalnya default base FX atau fallback setting, ditulis dalam description/default bila dapat dipastikan dari source.
- Perilaku endpoint yang hanya mengembalikan `{status:"ok"}` dibedakan dari endpoint yang mengembalikan entity.

---

## 4. Inventaris endpoint yang wajib diverifikasi dan didokumentasikan

Gunakan path absolut sesuai mounting. Di bawah ini daftar target minimum; setiap item wajib dicocokkan ulang dengan source sebelum implementasi.

### 4.1 Health

1. `GET /healthz`
   - Public.
   - Response sukses `200` JSON `{"status":"ok"}`.
   - Dokumentasikan content type JSON.
   - Tidak perlu memasukkan secret/security.

### 4.2 Auth: `/api/v1/auth`

Dari `Handler.Routes()`:

1. `POST /api/v1/auth/login`
   - Public.
   - Body model login dari `internal/models`.
   - PIN wajib enam digit numerik berdasarkan `ValidatePIN`.
   - Sukses membuat auth cookies; body `{status:"ok"}`.
   - Error invalid body/PIN `400`, PIN salah `401`, internal `500`.
   - Jangan menulis PIN contoh yang terlihat seperti credential runtime.

2. `POST /api/v1/auth/logout`
   - Public dari routing.
   - Menghapus auth cookies.
   - Sukses `{status:"ok"}`.

3. `POST /api/v1/auth/pin/reset-request`
   - Public.
   - Body request email.
   - Token dikirim melalui Telegram jika account dan binding ditemukan.
   - Sukses `{status:"token_sent"}`.
   - Catat error `400` untuk body/email/account/binding, `500` untuk internal.

4. `POST /api/v1/auth/pin/reset-verify`
   - Public.
   - Body token + PIN baru; baca sisa handler sampai selesai sebelum menulis schema.
   - Token memiliki expiry; jangan tampilkan contoh token yang dapat disalahartikan sebagai token nyata.
   - Dokumentasikan `400`, `500`, dan sukses aktual.

5. `GET /api/v1/auth/me`
   - Auth required oleh group.
   - Sukses `{user_id:string}`.
   - `401` bila user ID tidak tersedia.

6. `POST /api/v1/auth/pin`
   - Auth required.
   - Body current PIN + new PIN.
   - Sukses `{status:"ok"}`.
   - `400` validasi, `401` current PIN salah, `500` internal.

7. `POST /api/v1/auth/pin/change-request`
   - Auth required.
   - Body new PIN.
   - Token change dikirim Telegram.
   - Sukses `{status:"token_sent"}`.
   - Dokumentasikan `400`, `401`, `500` sesuai handler.

8. `POST /api/v1/auth/pin/verify`
   - Auth required.
   - Body token + new PIN.
   - Sukses `{status:"ok"}`.
   - `400` token invalid/expired atau PIN invalid, `403` token user mismatch, `401` middleware, `500`.

9. `POST /api/v1/auth/telegram/register-token`
   - Auth required.
   - Tanpa body.
   - Sukses menghasilkan token, `telegram_url`, `bot_username`; jika sudah linked dapat status `already_linked` dan chat ID.
   - Jangan menyamarkan token response; schema boleh menyatakan string token, tetapi contoh harus placeholder.

10. `GET /api/v1/auth/telegram/status`
    - Auth required.
    - Sukses `{linked:boolean, chat_id:string}`.
    - Pertimbangkan privacy: contoh chat ID harus placeholder atau kosong.

### 4.3 Accounts: `/api/v1/accounts`

Dari delivery handler:

1. `GET /api/v1/accounts/` — array `AccountWithBalance`.
2. `POST /api/v1/accounts/` — body create account; sukses `{id:string}`.
3. `GET /api/v1/accounts/summary` — object map currency ke balance.
4. `GET /api/v1/accounts/{id}` — `AccountWithBalance`; `404` jika tidak ditemukan.
5. `PUT /api/v1/accounts/{id}` — body partial update; sukses `{status:"ok"}`; currency tidak boleh diubah jika account punya transaksi (`400`).
6. `DELETE /api/v1/accounts/{id}` — sukses `{status:"ok"}`.

Schema wajib mencakup:

- `Account`;
- `AccountWithBalance`;
- `CreateAccountInput`;
- `UpdateAccountInput`;
- `CurrencyBalanceSummary` bila memang response aktual memakai array; jika handler mengubahnya menjadi object map, dokumentasikan response actual map dan jelaskan domain type tidak dipakai langsung.

Account type enum boleh memakai nilai dari domain: `bank`, `cash`, `card`, `investment`, `other`, dengan catatan handler tidak terlihat melakukan validasi enum jika memang demikian.

### 4.4 Categories: `/api/v1/categories`

1. `GET /api/v1/categories/` — array category.
2. `GET /api/v1/categories/tree` — tree category.
3. `POST /api/v1/categories/` — body create; name wajib; sukses `{id}`.
4. `PUT /api/v1/categories/reorder` — body `{items:[{id,sort_order}]}`; sukses `{status:"ok"}`.
5. `GET /api/v1/categories/{id}` — category atau `404`.
6. `PUT /api/v1/categories/{id}` — partial update; sukses `{status:"ok"}`.
7. `DELETE /api/v1/categories/{id}` — sukses `{status:"ok"}`.

Catatan penting:

- `Kind` dan `ParentID` punya rule leaf/children tertentu; tulis pada description, bukan mengarang validation schema jika handler hanya melakukan sebagian validation.
- `Children` optional recursive.
- `ReorderItem` harus menggunakan JSON property aktual; source struct tanpa tag perlu diverifikasi decoder naming behavior dan repository.

### 4.5 Transactions: `/api/v1/transactions`

1. `GET /api/v1/transactions/`
   - Query: `q`, `sort`, `order`, `page`, `limit`, `account_id`, `category_id`, `type`, `status`, `date_from`, `date_to`.
   - `date_from`/`date_to` diparse RFC3339; invalid value saat ini diabaikan oleh handler, jangan klaim selalu 400 kecuali source berubah.
   - Response array transaction.

2. `POST /api/v1/transactions/`
   - Body create transaction.
   - Sukses `{id}`.

3. `POST /api/v1/transactions/batch`
   - Body `BatchInput` dengan IDs/action/category/labels.
   - Sukses `{affected:number,errors:[]}`.
   - Enum action berdasarkan konstanta: `delete`, `recategorize`, `relabel`.

4. `GET /api/v1/transactions/search`
   - Query `q`.
   - Limit internal 50.
   - Array transaction.

5. `GET /api/v1/transactions/autocomplete`
   - Query `q`.
   - Limit internal 10.
   - Array string.

6. `GET /api/v1/transactions/summary`
   - Query `category_id`, `account_id`, `type`.
   - Array `SummaryItem`.

7. `GET /api/v1/transactions/{id}` — transaction atau `404`.
8. `PUT /api/v1/transactions/{id}` — partial update; sukses `{status:"ok"}`.
9. `DELETE /api/v1/transactions/{id}` — sukses `{status:"ok"}`.

Schema wajib mencakup:

- `Transaction`;
- `Attachment`;
- `CreateTransactionInput`;
- `UpdateTransactionInput`;
- `TransactionFilter` sebagai dokumentasi parameter/konsep, tetapi jangan mengklaim filter body dipakai GET jika source memakai query;
- `SummaryItem`.

Enum transaction type: `income`, `expense`, `transfer`. Enum status: `cleared`, `pending`, `reconciled`. Semua harus diverifikasi terhadap penggunaan repository.

### 4.6 Templates: `/api/v1/templates`

1. `GET /api/v1/templates/` — array template.
2. `POST /api/v1/templates/` — create body; `{id}`.
3. `GET /api/v1/templates/{id}` — template atau `404`.
4. `PUT /api/v1/templates/{id}` — partial update; `{status:"ok"}`.
5. `DELETE /api/v1/templates/{id}` — `{status:"ok"}`.
6. `POST /api/v1/templates/{id}/apply` — source delivery saat ini membaca template lalu mengembalikan template, bukan membuat transaction. Dokumentasikan perilaku aktual, jangan menamai response sebagai created transaction.

Schema template memakai field transaction-like: name, type, amount, currency, account IDs, category, labels, note, payee/payer, payment type, warranty, status, place, coordinates, attachments.

### 4.7 Scheduled: `/api/v1/scheduled`

1. `GET /api/v1/scheduled/` — array scheduled transaction.
2. `POST /api/v1/scheduled/` — create body; `{id}`.
3. `GET /api/v1/scheduled/{id}` — entity atau `404`.
4. `PUT /api/v1/scheduled/{id}` — partial update; `{status:"ok"}`.
5. `DELETE /api/v1/scheduled/{id}` — `{status:"ok"}`.
6. `POST /api/v1/scheduled/{id}/toggle` — source delivery memanggil toggle dengan `true`; dokumentasikan status response aktual.

Perhatian:

- Package `internal/scheduled/handler.go` juga memiliki `/preset`, tetapi delivery handler yang dimount di `main.go` tidak memiliki route tersebut. Jangan masukkan `/preset` tanpa bukti wiring berubah.
- `transaction_template` adalah `json.RawMessage`; schema boleh memakai object/map fleksibel dengan description field expected.
- RRule, max occurrences, enabled, next run, occurrences done wajib didokumentasikan.

### 4.8 Reports: `/api/v1/reports`

1. `GET /api/v1/reports/spending`
   - Query `period`, `group_by`, `account_id`.
   - Response `{period:string, groups:[...]}`.
   - `groups` menggunakan hasil repository; baca `models/reports.go` dan implementation untuk property aktual.

2. `GET /api/v1/reports/trend`
   - Query `months`, `account_id`.
   - Response `{series:[...]}`.
   - `months` diparse integer; perilaku default/invalid harus mengikuti source.

Schema report minimal:

- `Group`/SpendingGroup equivalent: id, name, color, kind, amount, currency, count, children.
- `MonthData`/TrendPoint equivalent: month, income, expense, net, currency.
- Jangan membuat schema bernama domain yang tidak ada tanpa alasan; gunakan definition yang jelas dan konsisten.

### 4.9 Settings: `/api/v1/settings`

1. `GET /api/v1/settings/`
   - Response object map string ke string.
   - Handler mengisi default: theme, language, base_currency, date_format, number_format, week_start, backup_passphrase_hash, totp_enabled.
   - Jangan menampilkan contoh hash yang realistis; gunakan empty/placeholder.

2. `PUT /api/v1/settings/`
   - Body object map string ke string.
   - Sukses `{status:"ok"}`.

3. `GET /api/v1/settings/currencies`
   - Response array string; fallback `USD` jika tidak tersedia.

4. `POST /api/v1/settings/currencies`
   - Body `{code:string}`; response array string.
   - `400` jika body/code invalid.

5. `GET /api/v1/settings/{key}`
   - Response `{key:string,value:string}`.
   - Default known keys dapat dikembalikan; unknown bisa `404`.

6. `PUT /api/v1/settings/{key}`
   - Body `{value:string}`.
   - Sukses `{status:"ok"}`.

Perhatikan route ordering: `/currencies` harus dipahami sebagai static route, sedangkan `/{key}` generic route. Swagger path keduanya tetap ditulis terpisah.

### 4.10 FX: `/api/v1/fx`

1. `GET /api/v1/fx/rates`
   - Query `base`, `date`.
   - Default base `JPY` pada delivery handler.
   - `date=latest` berarti current date behavior; date valid memakai `YYYY-MM-DD`.
   - Response `{base:string,date:string,rates:object[string]number}`.

2. `GET /api/v1/fx/convert`
   - Query wajib `from`, `to`, `amount`.
   - Invalid parameter `400`.
   - Same currency response converted equals amount.
   - General response `{from,to,amount,converted}`.
   - Rate fallback behavior harus ditulis secara hati-hati: handler fallback rate 1.0 jika rate lookup tidak menghasilkan valid pair; jangan menyembunyikan perilaku ini.

### 4.11 Shared: `/shared`

1. `POST /shared/`
   - Source menggunakan `userID := "system"`; ini harus dicatat sebagai implementasi saat ini, bukan dipresentasikan sebagai user identity yang berasal dari auth.
   - Body account ID dan optional expires_at RFC3339.
   - Response id, token, share_url.
   - Contoh token harus placeholder.

2. `GET /shared/`
   - Response array share.

3. `DELETE /shared/{id}`
   - Response `{status:"revoked"}` atau `404`.

4. `GET /shared/public/{token}`
   - Public share lookup.
   - Response account, transactions, attribution.
   - Token parameter wajib path.
   - Jangan menaruh token sample nyata.

Catatan keamanan wajib di hasil audit: top-level shared create/list/revoke saat ini tidak terlihat dibatasi auth dan memakai system user. Ini bukan diperbaiki dalam task dokumentasi kecuali user meminta security fix terpisah.

---

## 5. Audit source sebelum implementasi

### Langkah 5.1 — Baca semua route mounted

- Buka `backend/cmd/server/main.go`.
- Catat prefix dan handler constructor.
- Cocokkan setiap constructor dengan file delivery/auth.
- Buat tabel route live: method, full path, handler method, auth middleware, response status.
- Pastikan trailing slash tidak salah di Swagger. Chi menerima route `/`; dokumentasi dapat memakai `/api/v1/accounts/` atau path canonical tanpa trailing slash, tetapi konsisten dan tidak membuat dua operation untuk endpoint yang sama.

### Langkah 5.2 — Baca handler hingga perilaku response

Untuk tiap handler:

- catat semua query parameter yang dibaca;
- catat path parameter;
- catat body struct dan validasi;
- catat status success aktual;
- catat error string/status aktual;
- catat `Content-Type` yang diset;
- catat bentuk response map/array/entity;
- catat default value;
- catat side effect yang relevan untuk description;
- catat route yang secara source tidak diberi auth.

Jangan hanya mengandalkan nama method atau domain interface.

### Langkah 5.3 — Baca model/domain

- Cari semua struct yang dipakai langsung handler.
- Baca JSON tags.
- Tandai field pointer sebagai optional/nullable.
- Tandai `omitempty` sebagai optional pada output.
- Tandai field input tanpa JSON tag untuk diverifikasi terhadap `encoding/json`.
- Baca const enum yang benar-benar digunakan.
- Hindari schema duplikat untuk struct identik jika reusable definition cukup.

### Langkah 5.4 — Audit package orphan

- Search `Routes()` di seluruh `backend/internal`.
- Hanya masukkan route jika constructor/route benar-benar dipanggil dari `main.go` atau router yang dipanggil dari sana.
- Buat daftar route yang ditemukan tetapi tidak live.
- Masukkan daftar itu dalam laporan plan/hasil implementasi sebagai excluded endpoints.
- Jangan mengubah route wiring selama pekerjaan docs.

---

## 6. Desain format Swagger

### 6.1 Root document

Pertahankan:

- `swagger: "2.0"`;
- `info.title: Kiroku API`;
- `info.version: "1.0"`;
- `info.description` yang menjelaskan personal finance API dan cookie auth tanpa secret;
- `host: localhost:8080` sebagai default development host;
- `basePath: /api/v1` untuk API v1.

Karena `/healthz` dan `/shared` berada di luar base path `/api/v1`, verifikasi perilaku renderer Swagger. Swagger 2.0 `basePath` berlaku global; jika root basePath membuat top-level route tidak dapat dipanggil dengan benar di UI, dokumentasikan tradeoff dan pertimbangkan spec structure/description yang jelas. Jangan mengarang host/path.

### 6.2 Tags

Gunakan tag konsisten, minimal:

- Health
- Authentication
- Accounts
- Categories
- Transactions
- Templates
- Scheduled Transactions
- Reports
- Settings
- Foreign Exchange
- Sharing

Setiap operation tepat memiliki tag domain utama.

### 6.3 Security

Definisi:

```json
"CookieAuth": {
  "type": "apiKey",
  "name": "auth_token",
  "in": "cookie"
}
```

Aturan dokumentasi:

- login, logout, reset-request, reset-verify: tanpa `security` requirement karena route public.
- `/auth/me`, change PIN, Telegram authenticated operations: `security: [{"CookieAuth": []}]`.
- accounts/categories/transactions/templates/scheduled/reports/settings/fx: verifikasi source. Bila dokumentasi menampilkan intended auth, wajib menulis gap dalam hasil. Jangan mengubah source secara diam-diam.
- `/shared` dan `/healthz`: tanpa security karena mounted public.

### 6.4 Parameters

- Path parameter: `in: path`, `required: true`.
- Query: `in: query`, required hanya jika handler benar-benar menolak jika kosong.
- Body: satu `in: body`, schema reusable.
- Cookie auth jangan ditulis sebagai body/header parameter manual jika security definition cukup.
- Date-time: `type:string`, `format:date-time`.
- Date only: `type:string`, `format:date`.
- Currency/code: string, contoh aman seperti `USD`/`JPY`.
- Amount: number/format double.
- Page/limit/months/count: integer.

### 6.5 Responses

Buat reusable definitions untuk response umum jika Swagger 2.0 file tetap mudah dirawat:

- `StatusOK`;
- `IDResponse`;
- `TokenSentResponse`;
- `ErrorResponse` atau `ErrorText`.

Namun jangan memalsukan JSON error: handler memakai `http.Error`, yang secara default menulis text/plain dengan newline. Pilihan dokumentasi:

- response schema `string` + `text/plain`; atau
- schema object hanya jika handler nyata mengembalikan JSON.

Untuk tiap operation tulis minimal response yang mungkin dari handler. Jangan menambahkan 404 pada endpoint yang source tidak pernah mengembalikan 404 hanya karena konvensi REST, kecuali description menyatakan kemungkinan repository error dipetakan ke 500.

### 6.6 Examples

- Gunakan examples yang kecil dan representatif.
- Tidak perlu memasukkan contoh lengkap 20 field pada setiap endpoint jika reusable schema sudah jelas, tetapi response penting harus punya contoh.
- Semua ID/token/email/chat ID bersifat placeholder non-secret.
- Jangan gunakan PIN yang terlihat sebagai PIN default.
- Jangan gunakan URL tunnel.
- Share URL cukup `/shared/{token}` atau `https://example.invalid/shared/{token}` hanya jika Swagger membutuhkan full URL; prefer relative placeholder dan jangan mengarang URL produksi.

---

## 7. Strategi implementasi

### Opsi A — Update generated `docs.go` secara langsung

Gunakan jika:

- dependency swag tersedia;
- tidak ada source annotations yang cukup untuk generator;
- target user hanya membutuhkan docs yang aktif dan buildable.

Langkah:

1. Pertahankan package/import/register existing.
2. Replace `docTemplate` dengan JSON template valid yang memuat `paths`, `definitions`, `tags`, `securityDefinitions`.
3. Gunakan escaping compatible dengan `swag.Spec` (`marshal`, `escape`) bila diperlukan.
4. Pastikan template dapat dirender oleh `swag` tanpa invalid JSON.
5. Pastikan `SwaggerInfo` fields tetap benar.

### Opsi B — Tambahkan annotations lalu generate

Gunakan hanya jika repository sudah menyediakan workflow generator yang terverifikasi:

1. Cari `swag init`, `go:generate`, Makefile/script, atau dependency CLI.
2. Pastikan annotations dapat menjangkau semua handler dan models.
3. Tambahkan annotations hanya pada source yang diperlukan.
4. Generate file memakai command yang memang tersedia.
5. Review diff generated.
6. Jangan mengklaim generated docs jika command tidak tersedia/reproducible.

Default plan: pilih opsi dengan diff paling kecil dan hasil paling dapat diverifikasi. Jangan menambah dependency hanya untuk menghindari penulisan spec jika existing `swag.Spec` dapat memuat template lengkap.

### Struktur internal yang disarankan

- Root info/security.
- `paths` dikelompokkan secara alfabetis per prefix.
- `definitions` dikelompokkan per domain.
- Reusable response definitions setelah domain schemas.
- Tidak ada komentar tambahan di file kecuali diperlukan oleh generator.

---

## 8. Detail schema yang wajib dipetakan

### Account

Properties minimum dari domain:

- id: string;
- name: string;
- currency: string;
- type: string;
- icon: string;
- color: string;
- archived: boolean;
- created_at: string date-time;
- updated_at: string date-time.

`AccountWithBalance` menambah `balance:number double`.

Input create/update harus dibedakan: update memakai nullable/pointer optional properties.

### Category

- id, parent_id, name, icon, color, kind, sort_order, depth, path, children, archived, created_at, updated_at.
- `parent_id` dan `kind` nullable.
- `children` recursive optional.

### Transaction

- id, type, date_time, amount, currency;
- optional from/to account IDs;
- category ID;
- labels;
- note, payee_payer, payment_type;
- warranty_until;
- status;
- place_name, place_lat, place_lng;
- attachments;
- template_id;
- created_at, updated_at, deleted_at.

Dokumentasikan coordinate sebagai double nullable.

### Template

Fields transaction-like ditambah `name`.

### ScheduledTransaction

- id;
- rrule;
- max_occurrences nullable;
- transaction_template object/raw JSON;
- enabled;
- next_run nullable date-time;
- occurrences_done integer;
- created_at/updated_at.

### Reports

Gunakan source `models/reports.go`:

- Group: id, name, color, kind optional, amount, currency, count, children optional recursive.
- MonthData: month, income, expense, net, currency.

### Settings

- `Settings`: object additionalProperties string.
- `SettingResponse`: key/value.
- `SettingReq`: value.

### FX

- Rates response: base, date, rates map string→double.
- Convert response: from, to, amount, converted.

### Share

- id, account_id, account_name optional, token, expires_at nullable, created_at, last_accessed_at nullable.
- Shared view: account, transactions, attribution.

### Auth

Baca `internal/models` auth request structs lengkap sebelum implementasi:

- LoginReq;
- ChangePINReq;
- PinChangeRequestReq;
- VerifyPinChangeReq;
- RequestPinResetReq;
- VerifyPinResetReq.

Jangan menebak field token/email/new PIN dari nama method jika struct belum dibaca.

---

## 9. Matrix status dan error

Buat matrix sebelum coding:

| Domain | Sukses utama | Bad request | Unauthorized | Forbidden | Not found | Internal |
|---|---:|---:|---:|---:|---:|---:|
| health | 200 | - | - | - | - | - |
| auth login | 200 | 400 | 401 invalid PIN | - | - | 500 |
| auth session | 200 | 400 | 401 | 403 pada token mismatch jika berlaku | - | 500 |
| CRUD resources | 200 | 400 | source/wiring check | source check | handler-specific | 500 |
| shared | 200 | 400 | source public | - | 404 | 500 |
| fx | 200 | 400 convert invalid | source/wiring check | - | - | 500 |

Matrix final harus memakai source aktual, bukan tabel ini secara buta. Endpoint yang tidak mengeluarkan status tertentu tidak perlu memaksa response tersebut ke spec.

---

## 10. Verifikasi implementasi

### 10.1 Verifikasi statis

Dari repo root:

- periksa `git diff -- backend/docs/docs.go`;
- pastikan hanya file yang diminta berubah;
- search `paths: {}` dan pastikan tidak ada lagi pada active template;
- search semua mounted route dan cocokkan jumlah operation;
- search `/preset`, `/scan`, `/export`, `/import`, `/analytics` dalam docs; pastikan hanya muncul jika live route terbukti dimount.

### 10.2 Verifikasi build/test backend

Dari `backend/` jalankan:

```bash
go build ./...
go test ./...
```

Sesuai `AGENTS.md`, keduanya wajib dilakukan untuk backend change.

### 10.3 Verifikasi rendering spec

Jika server bisa dijalankan tanpa secret exposure:

```bash
KIROKU_TUNNEL_ENABLED=false go run ./cmd/server
```

Kemudian cek:

- `/healthz` mengembalikan status ok;
- `/api/docs/` dapat dibuka;
- `/api/docs/swagger.json` atau endpoint JSON yang dirujuk Swagger handler mengandung `paths` non-empty;
- parse JSON dengan tool lokal/std library;
- operation count sesuai inventory;
- tidak ada `$ref` yang menunjuk definition tidak ada;
- tidak ada JSON invalid.

Jangan menyalin atau menampilkan cookie, token, URL tunnel, data database, atau secret output. Jika server membutuhkan config yang belum tersedia, catat sebagai unverifiable.

### 10.4 Verifikasi schema

Untuk setiap operation:

- semua path params terdefinisi di path;
- body schema definition ada;
- response schema definition ada;
- array memiliki `items`;
- map memakai `additionalProperties`;
- nullable pointer behavior tidak dibuat required secara salah;
- contoh response cocok dengan schema;
- route static dan parameterized tidak tertukar.

### 10.5 Verifikasi security

- `CookieAuth` terdaftar tepat sekali.
- Auth group endpoints punya security metadata sesuai keputusan yang terdokumentasi.
- Public routes tidak meminta cookie.
- Missing middleware pada data routes tetap dicatat sebagai known gap, tidak diperbaiki diam-diam.

---

## 11. Acceptance criteria

Implementasi dianggap selesai hanya jika seluruh kondisi berikut terpenuhi:

1. `backend/docs/docs.go` tidak lagi menghasilkan `paths: {}`.
2. Semua route live dari `main.go` tercantum atau secara eksplisit dijelaskan pengecualiannya.
3. Tidak ada route unmounted yang dimasukkan sebagai API live.
4. Auth route lengkap mencakup seluruh method dari `auth.Handler.Routes()`.
5. Accounts, categories, transactions, templates, scheduled, reports, settings, FX, shared, health lengkap.
6. Query/path/body parameter terdokumentasi dengan tipe dan required behavior yang sesuai source.
7. Response success/error tidak mengarang bentuk JSON yang tidak pernah dikirim handler.
8. Schema domain utama tersedia dan reusable.
9. Token/PIN/secret tidak terekspos dalam examples atau output.
10. `CookieAuth` tersedia.
11. `go build ./...` sukses dari `backend/`.
12. `go test ./...` sukses dari `backend/`.
13. Spec JSON valid dan dapat dipakai Swagger UI.
14. Diff tidak menyentuh file unrelated.
15. Laporan akhir menyebutkan:
    - file berubah;
    - route yang didokumentasikan;
    - verification yang sukses;
    - hal yang tidak dapat diverifikasi;
    - known security gap yang ditemukan namun tidak diubah.

---

## 12. Risiko dan mitigasi

### Risiko 1 — `docs.go` generated tidak cocok dengan generator

- **Risiko:** perubahan manual hilang saat `swag init` dijalankan.
- **Mitigasi:** cek apakah repository memiliki generator workflow. Jika tidak ada, dokumentasikan bahwa file adalah runtime spec maintained manually. Jika ada, tambahkan annotations/generate hanya setelah command diverifikasi.

### Risiko 2 — Base path global untuk top-level route

- **Risiko:** `/shared` dan `/healthz` berada di luar `/api/v1`, sementara Swagger 2.0 hanya memiliki satu basePath.
- **Mitigasi:** verifikasi Swagger UI; jika limitation nyata, jelaskan pada description dan pertimbangkan canonical server URL/path strategy tanpa mengubah runtime route secara sembarangan.

### Risiko 3 — Handler lama dan handler delivery berbeda

- **Risiko:** source memiliki duplicate package dengan route lebih lengkap, misalnya scheduled `/preset`.
- **Mitigasi:** `main.go` adalah source of truth untuk live API; hanya delivery handler yang dimount yang didokumentasikan.

### Risiko 4 — Field input tanpa JSON tags

- **Risiko:** property name Swagger salah.
- **Mitigasi:** baca `encoding/json` behavior dan repository query mapping; bila ambigu, gunakan nama yang sudah dipakai frontend/request payload, lalu catat gap untuk follow-up, bukan menebak diam-diam.

### Risiko 5 — Security metadata menutupi missing middleware

- **Risiko:** docs terlihat aman padahal route belum memakai AuthMiddleware.
- **Mitigasi:** pisahkan intended API security dari actual wiring; tulis known gap eksplisit dalam hasil.

### Risiko 6 — Error body sebenarnya text/plain

- **Risiko:** Swagger mendokumentasikan JSON padahal handler memakai `http.Error`.
- **Mitigasi:** dokumentasikan string/text response, atau ubah handler hanya jika user meminta kontrak error JSON—bukan bagian task ini.

### Risiko 7 — Spec terlalu besar dan sulit dirawat

- **Risiko:** ratusan endpoint/schema menghasilkan file tidak terbaca.
- **Mitigasi:** reusable definitions, tag, operationId optional konsisten, no duplicate inline schema, no undocumented orphan routes.

---

## 13. Urutan kerja implementasi yang direkomendasikan

1. Baca `AGENTS.md` dan jangan buka `.env`.
2. Baca `main.go` dan catat route live.
3. Baca seluruh delivery handler yang dimount.
4. Baca seluruh auth handler yang dimount.
5. Baca model/domain request-response terkait.
6. Baca reports implementation untuk bentuk group/trend aktual.
7. Audit duplicate/orphan route dan daftar exclusion.
8. Buat endpoint/status/schema matrix.
9. Tentukan apakah docs manual atau generator-based berdasarkan repository tooling.
10. Update `backend/docs/docs.go` dengan root metadata, security, tags, paths, definitions.
11. Jalankan static checks terhadap path count dan unresolved refs.
12. Jalankan `go build ./...` dari backend.
13. Jalankan `go test ./...` dari backend.
14. Jalankan server tanpa tunnel jika environment memungkinkan.
15. Fetch/parse Swagger JSON tanpa mencetak secret.
16. Review diff final.
17. Catat known gaps dan verification limitation.

---

## 14. Perintah verifikasi final

Dari repo root untuk review diff:

```bash
git status --short
git diff -- backend/docs/docs.go
```

Dari backend:

```bash
go build ./...
go test ./...
```

Opsional jika config runtime tersedia dan aman:

```bash
KIROKU_TUNNEL_ENABLED=false go run ./cmd/server
```

Jangan menjalankan perintah yang mencetak `.env`, token, cookie, database, `backend/data/`, atau tunnel URL.

---

## 15. Follow-up di luar scope

Tidak dikerjakan dalam implementasi docs kecuali diminta terpisah:

- memasang `AuthMiddleware` pada accounts/categories/transactions/templates/scheduled/reports/settings/fx;
- mengamankan `/shared` create/list/revoke;
- mengganti `userID := "system"` pada shared handler;
- memperbaiki CORS wildcard reflection dengan credentials;
- memperbaiki default PIN;
- memasang package analytics/receipts/impexp/backup yang belum dimount;
- menambah endpoint `/scheduled/preset`, receipts, analytics, import/export;
- mengubah error response menjadi JSON;
- membuat OpenAPI 3 jika runtime saat ini Swagger 2.0;
- menambahkan contract tests atau integration tests baru.

Semua follow-up tersebut boleh menjadi issue terpisah setelah dokumentasi selesai dan route contract disepakati.

---

## 16. Definition of done plan ini

Plan ini siap dipakai implementer apabila:

- endpoint live sudah dipisahkan dari endpoint orphan;
- setiap domain memiliki daftar method/path/parameter/response;
- schema utama dan auth model sudah diidentifikasi;
- security gap tercatat;
- acceptance criteria dan command verifikasi tersedia;
- implementer tidak perlu menebak route berdasarkan nama file;
- implementer tahu mana yang boleh diubah dan mana yang out-of-scope.
