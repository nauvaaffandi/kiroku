# Labels Feature — Implementation Plan

## 0. Dokumen, status, dan aturan eksekusi

- Status: rencana implementasi; belum merupakan klaim bahwa fitur sudah selesai.
- Tujuan dokumen: menjadi spesifikasi implementasi yang cukup rinci agar pekerjaan dapat dilanjutkan tanpa menebak kontrak, lokasi file, edge case, atau verifikasi.
- Sasaran produk: pengguna dapat membuat, mengatur, mencari, memilih, memakai, dan menghapus label transaksi.
- Sumber kebenaran aturan proyek: `AGENTS.md`, lalu `DESIGN.md`, lalu pola kode yang telah terpasang. Jika rencana ini bertentangan dengan keduanya, perbaiki implementasi/rencana sebelum coding.
- Rencana lama: `.working/plans/labels-feature.md`. Dokumen ini menggantikannya sebagai instruksi eksekusi yang lebih spesifik, tetapi tidak menghapus rencana lama.
- Instruksi penting: baca `AGENTS.md` penuh sebelum perubahan; untuk UI baca `DESIGN.md` penuh; jangan buka `.env`; jangan menampilkan, mengubah, menghapus, atau mengandalkan isi `backend/data/`, database, token, PIN, maupun URL tunnel.
- Jangan melakukan commit kecuali diminta eksplisit.
- Kerjakan berurutan per tahap. Setelah tiap tahap, perbarui status checklist dan jalankan verifikasi tahap terkait sebelum melangkah.
- Perubahan worktree yang sudah ada: `backend/internal/db/migrations/015_labels.sql`, `016_transaction_labels.sql`, `backend/internal/domain/label.go`, `backend/internal/repository/sqlite/labels.go` tampak sebagai file untracked; inspeksi isinya, pertahankan yang benar, revisi seperlunya. Jangan menganggapnya sudah benar atau menghapusnya tanpa alasan.
- Ada risiko kode parsial sebelumnya pada `backend/internal/repository/sqlite/transactions.go`; inspect `git diff` dan file sebelum melanjutkan. Pastikan jangan mempertahankan penggandaan argumen ID di `Update`.
- Nomor migrasi: jangan menerima asumsi bahwa 015/016 otomatis aman hanya dari nama file. Inventaris seluruh `backend/internal/db/migrations/*.sql`, aturan `runMigrations()` di `backend/internal/db/db.go`, dan file yang sudah ada pada repo. Nomor 008 memang duplikat. Jangan mengubah migration yang sudah dirilis/applied. Tetapkan nomor baru yang aman berdasarkan kebijakan migrasi proyek dan kondisi repo; jika 015/016 memang file lokal yang belum pernah dirilis, pastikan urutannya tidak mengubah database user. Bila perlu, gunakan nomor berurutan setelah nomor tertinggi yang dipakai dan jelaskan keputusan di PR/summary.

---

## 1. Definisi fitur

### 1.1 Istilah

- **Label**: metadata datar/non-hierarkis untuk klasifikasi transaksi lintas kategori. Contoh: `Liburan`, `Langganan`, `Kerja`, `Pajak`.
- **Label ID**: identifier stabil yang tersimpan pada junction; transaksi tidak menyimpan label berdasarkan nama sebagai sumber utama.
- **Label nama legacy**: string lama pada `transactions.labels`, umumnya nama label yang dipisah koma. String ini tetap dibaca dan dipertahankan agar transaksi lama tidak kehilangan tampilan.
- **Junction**: relasi many-to-many `transaction_labels(transaction_id, label_id)`.
- **Inline create**: membuat label baru dari picker pada form transaksi lalu otomatis memilih label baru.

### 1.2 Perilaku yang diminta

1. CRUD label melalui API dan halaman Labels.
2. Label memiliki nama unik, warna, ikon opsional, deskripsi opsional, urutan manual.
3. Banyak label dapat ditempel ke satu transaksi; satu label dapat digunakan pada banyak transaksi.
4. Form transaksi dapat memilih dan menghapus label terpilih, mencari label, dan membuat label baru langsung.
5. Edit transaksi dapat memuat dan mengganti label ID.
6. Menghapus label menghapus relasi junction secara konsisten; transaksi tetap ada.
7. Nilai legacy comma-separated `transactions.labels` tetap dapat ditampilkan, tidak ditafsirkan sebagai ID, dan tidak hilang ketika perubahan label ID dilakukan.
8. Sidebar menampilkan Labels di kelompok Organization; Accounts, Categories, Labels, Scheduled berada di kelompok itu; Reports dan Shared Access di kelompok Insights.
9. Semua teks UI baru diterjemahkan dalam `en`, `ja`, `id`.
10. Semua route data baru wajib mengikuti middleware autentikasi yang diwajibkan `AGENTS.md`, bukan menambah route terbuka baru.

### 1.3 Bukan sasaran (kecuali keputusan produk eksplisit)

- Hierarki label, parent ID, atau label bertingkat.
- Pemfilteran transaksi berdasarkan label di UI/API pada iterasi pertama, kecuali backend filter yang sudah tersedia dan masuk scope dengan aman.
- Sistem permission multi-user, karena aplikasi single-user saat ini.
- Mengubah nama label untuk menulis ulang teks legacy historis.
- Menambah dependency drag-and-drop baru.
- Mengubah desain global atau migrasi massal komponen UI.
- Analitik pengeluaran per label sebagai requirement wajib. Audit widget yang ada dahulu; integrasi chart opsional dan jangan menambahkan chart endpoint baru bila backend tidak menyediakannya.

---

## 2. Konteks codebase yang sudah diketahui

### 2.1 Backend

- Layer aktif: `backend/internal/domain/`, `backend/internal/repository/sqlite/`, `backend/internal/delivery/http/`, wiring `backend/cmd/server/main.go`.
- Database menggunakan `database/sql`, SQLite pure Go, foreign keys aktif melalui DSN di `backend/internal/db/db.go`.
- SQL migrations di-embed (`//go:embed migrations/*.sql`) dan diterapkan berdasarkan nomor lebih besar dari `schema_version`; file dengan nomor duplikat dapat diskip setelah nomor pertama.
- `domain.CategoryRepository`, `repository/sqlite/categories.go`, `delivery/http/categories.go` merupakan pola CRUD/reorder terdekat.
- `domain.Transaction` memiliki `Labels []string` JSON `labels`; kolom database lama adalah teks comma-separated.
- `TransactionsRepo` melakukan query langsung untuk list/get/search/summary; method `scanTransactions` scan kolom berurutan. Setiap modifikasi SELECT harus mempertahankan konsistensi jumlah dan urutan kolom dengan `Scan`.
- Transaksi dihapus lunak (`deleted_at`). Relasi harus cascade hanya ketika transaksi benar-benar terhapus, bukan soft-delete.
- API route pada `main.go` saat ini memasang sejumlah data handler tanpa middleware auth (gap known). Route baru tidak boleh menambah risiko ini; gunakan `auth.AuthMiddleware(cfg)` pada subtree labels sesuai signature middleware saat ini, dan jangan mengubah route existing tanpa scope/approval terpisah.
- `TransactionRepository` sekarang dipakai oleh HTTP transaction handler; menambah method interface akan memaksa semua implementasi/mock yang memenuhi interface diperbarui. Search semua implementasi sebelum perubahan.
- Package internal yang belum diwire tidak boleh dijadikan dasar pola fitur.

### 2.2 Frontend

- Vue 3, `<script setup lang="ts">`, Pinia setup stores, axios instance `@/api/client`, alias `@/`.
- `frontend/src/stores/categories.ts` dan `frontend/src/widgets/categories/` pola store/UI referensi; komponen kategori legacy memiliki teks hard-coded dan bukan standar yang boleh ditiru untuk string baru.
- `frontend/src/types/index.ts` barrel exports types; cek bentuk `frontend/src/widgets/index.ts` dan barrel domain sebelum membuat barrel baru.
- `frontend/src/components/ui/index.ts` saat inspeksi hanya mengekspor `BaseButton`, `BaseInput`, `BaseSelect`, `BaseModal`; tidak ada `BaseCombobox` yang terverifikasi. Jangan mengimpor komponen yang belum ada. Gunakan kontrol HTML accessible atau buat komponen minimal hanya jika benar-benar dibutuhkan.
- Route didefinisikan eksplisit di `frontend/src/router/index.ts`; halaman lazily loaded dan harus menandai `meta.requiresAuth: true`.
- Sidebar berada di `frontend/src/components/layout/Sidebar.vue`; nav item arrays menggunakan `t(...)`, icon components, serta route matching.
- Transaction form yang nyata berada di `frontend/src/widgets/transactions/TransactionFormModal.vue`, bukan path pada plan lama yang disebutkan tetapi tidak ada. Props saat ini: `modelValue`, `transaction`, `accounts`, `categories`; submit mengeluarkan `Partial<Transaction>`.
- Tipe frontend saat ini `Transaction.labels: string`; create/update request juga `labels?: string`. Store menerima `Partial<Transaction>` dan axios mengirim data itu.
- i18n files: `frontend/src/locales/en.ts`, `ja.ts`, `id.ts`; flat object groups `common`, `nav`, `transactions`, etc. `useI18n()` berasal dari `@/locales`.
- UI baru wajib pakai `<style scoped>` dan baris pertama style `@reference "tailwindcss";`; Tailwind hanya dalam `@apply`, class template semantic; gunakan token CSS yang terdaftar di `AGENTS.md`/`DESIGN.md`; jangan tambah hard-coded UI strings.
- `@lucide/vue` sudah terpasang; gunakan `Tag`/ikon yang benar-benar tersedia pada dependency. Hindari inline SVG baru.
- `npm run build` menjalankan vue-tsc dan Vite; `npm run test:run` terkonfigurasi tapi folder `frontend/tests/` mungkin belum ada. Verifikasi kondisi aktual sebelum menjanjikan test.

### 2.3 Titik yang harus diperiksa ulang saat mulai

- `git status --short`, `git diff`, dan isi perubahan yang sudah ada; jangan merusak perubahan user.
- Seluruh migrations dan urutan filename; pastikan nomor dan strategi upgrade.
- Semua implementasi `domain.TransactionRepository`, mock, callsites input transaction, dan JSON decoder handler.
- `transactions` schema, tipe kolom, legacy `labels`, pemakaian `labels` di frontend/API/filter/batch/export/import/backup.
- Wiring `AuthMiddleware` signature dan cara middleware diterapkan di chi saat ini.
- `DESIGN.md` bagian komponen/navigation/chips/forms/modal/mobile; jangan mengarang token/radius/spacing.
- Scripts package frontend, Go version/module path, test ignore rules.

---

## 3. Kontrak data dan keputusan teknis wajib

### 3.1 Schema label

Bentuk konseptual:

```sql
CREATE TABLE labels (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  color TEXT NOT NULL DEFAULT '#6366f1',
  icon TEXT,
  description TEXT,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Rincian:
- Name di-trim pada trust boundary; kosong/tidak valid ditolak.
- Keputusan case sensitivity harus eksplisit: SQLite `UNIQUE` default case-sensitive; plan awal tidak menetapkan normalisasi. Tinjau requirement existing API. Pilihan rekomendasi: trim whitespace, pertahankan kapitalisasi, gunakan unique constraint saat ini; jangan tiba-tiba case-fold kecuali diputuskan sebagai aturan produk.
- `color`: default `#6366f1` menurut plan awal. Validasi format `#RRGGBB` agar CSS injection/invalid input tidak masuk; jangan menerima arbitrary CSS string.
- `icon`: opsional, input plain text terbatas; jangan menerima markup/SVG HTML. Jika UI menggunakan emoji, maknanya harus user-chosen label icon, bukan dekorasi UI.
- `description`: opsional plain text; length limit wajar dan konsisten.
- `sort_order`: integer. Reorder harus atomik.
- Timestamp: gunakan konvensi SQLite saat ini. Trigger `updated_at` perlu diuji agar tidak recursive/bermasalah dengan recursive triggers; jika pola trigger tidak aman/terlalu rumit, perbarui di UPDATE query repository. Jangan menambah trigger tanpa test.
- Nama constraint/index unik: `UNIQUE(name)` telah menciptakan index implisit; jangan menambah index duplikat pada name tanpa query-plan alasan. `sort_order` index berguna hanya bila order/query pattern membutuhkannya.

### 3.2 Junction

```sql
CREATE TABLE transaction_labels (
  transaction_id TEXT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  label_id TEXT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
  PRIMARY KEY (transaction_id, label_id)
);
```

Rincian:
- Primary key mencegah relasi ganda.
- FK harus bekerja pada koneksi SQLite yang foreign_keys ON.
- Index tambahan pada kolom yang sudah merupakan prefix primary key mungkin redundan; periksa query usage sebelum membuat index redundant. Index `label_id` berguna untuk menghapus label / mendapatkan transaksi label.
- Junction tidak mengubah soft-delete semantics. Soft-delete transaksi tidak menghapus link agar restore kelak tetap konsisten; ini harus dipastikan dengan kebutuhan trash/restore dan didokumentasikan.
- Semua relasi yang diterima dari client harus tervalidasi: ID label tidak dikenal membuat operasi gagal secara atomik, bukan menyimpan subset diam-diam.

### 3.3 Legacy compatibility policy

Pilih dan implementasikan aturan berikut secara konsisten, bukan perkiraan:

1. `transactions.labels` kolom lama tetap ada; migration jangan rename/drop.
2. Existing legacy string tetap output `labels` sebagai string array/string sesuai backend JSON contract aktual. Frontend type harus cocok dengan contract yang benar-benar dikirim; jangan mengubah `labels` menjadi array tanpa migration API callsites.
3. Tambahkan `label_ids` sebagai `[]string` untuk relasi baru.
4. `label_ids` berasal dari junction. Ambil hanya untuk transaksi yang ditampilkan; hindari query N+1 bila hasil list banyak. Implementasi batching/join agregasi harus sesuai SQLite version dan scan model; maksimal ukuran list API perlu diperhitungkan.
5. Untuk create baru, tulis IDs hanya ke junction. Kolom legacy `labels` tetap kompatibel: isi hanya jika request lama mengirim string/nama legacy. Jangan menulis ID ke kolom nama.
6. Untuk update, bedakan `label_ids` tidak dikirim dari dikirim array kosong. JSON decoding ke `[]string` biasa tidak cukup membedakan absent vs `[]`; gunakan pointer/sentinel DTO (`*[]string`) atau custom presence field pada request boundary. Semantik: absent = jangan sentuh junction; `[]` = hapus semua junction; nonempty = replace set seluruhnya.
7. Jangan merge legacy names dengan junction IDs menjadi satu jenis array ambigu. Pertahankan dua fields API (`labels` legacy, `label_ids` canonical). UI boleh menampilkan legacy labels terpisah sebagai chip teks readonly/compatibility indicator jika data itu ada; jangan map legacy strings ke IDs berdasarkan nama otomatis, karena nama mungkin tidak unik semantis, ambigu, atau berubah.
8. Backfill legacy values ke `labels` catalog/junction hanya jika ada migrasi data yang deterministik dan data ownership/rules disetujui. Plan lama menyebut migrasi nilai lama tetapi tidak menetapkan warna, duplikasi, parsing koma dalam nama, normalisasi, atau collision. Default aman untuk iterasi ini: jangan melakukan lossy automatic backfill; tetap tampilkan legacy values, user dapat mengatur ulang label ID baru secara eksplisit. Jika produk mensyaratkan backfill, buat subplan dan migrasi teruji sebelum implementasi.
9. Jangan menyimpan label names dan label IDs secara bersamaan sebagai sumber kebenaran baru.
10. Edit/delete/reorder tidak boleh merusak data transaksi yang tidak terkait.

### 3.4 API contracts

Gunakan prefix `/api/v1` dan axios bersama pada frontend.

| Method | Route | Perilaku sukses | Input |
|---|---|---|---|
| GET | `/api/v1/labels` | `200`, array Label (termasuk `[]`) | optional query `q`, `sort`, `order` hanya bila repo mendukung dengan benar |
| POST | `/api/v1/labels` | `201` atau pola status repository yang konsisten, body `{id}` atau label; pilih satu contract dan konsisten | name required; color/icon/description/sort_order optional |
| GET | `/api/v1/labels/{id}` | `200` label | ID path |
| PUT | `/api/v1/labels/{id}` | `200` status atau updated label; pilih contract konsisten | patch fields optional |
| DELETE | `/api/v1/labels/{id}` | `200` JSON status atau `204`; konsisten existing handlers | ID path |
| PUT | `/api/v1/labels/reorder` | `200` JSON status | `{items:[{id,sort_order}]}` |
| GET | `/api/v1/transactions/{id}/labels` | diputuskan apakah dibutuhkan; plan lama mendeskripsikan method repo tapi handler tidak menyertakan endpoint | transaction ID |

Rincian kontrak:
- Chi static path `/reorder` harus didaftarkan sebelum route `/{id}` agar route variable tidak menangkap `reorder`.
- Request decode harus menolak JSON rusak; batas body jika sudah ada pola middleware. Validasi semua fields sebelum database write.
- `400` invalid input, `404` item tidak ada, `409` duplicate name/constraint conflict, `500` generic internal error tanpa detail SQL. Jangan menyamarkan error DB sebagai 404.
- Semua success response JSON set `Content-Type` bila body ada.
- Handler bergantung pada `domain` interface, tidak `database/sql`.
- API routes labels wajib di belakang `auth.AuthMiddleware(cfg)` sesuai instruksi `AGENTS.md`. Sebelum memasang, verifikasi middleware menerima handler/chi.Router seperti apa; jika auth middleware currently only defined for a specific route group, tentukan mount yang benar. Jangan memperluas auth bypass.
- Karena ada known security gap pada route data lain, summary implementasi harus jelas bahwa labels route aman sendiri tetapi aplikasi masih punya endpoint lain unauthenticated. Jangan melebar memperbaiki route lain tanpa request.
- API filter search: jika q disediakan, gunakan parameterized query (LIKE dengan escaping wildcard jika diharapkan literal); jangan concatenate user input ke SQL identifier.
- Sort/order adalah identifier SQL: whitelist `name`, `sort_order`, `created_at` dan arah `ASC`/`DESC`; jangan interpolasi mentah.

### 3.5 Domain contract

`backend/internal/domain/label.go` memuat:
- `Label` dengan JSON fields `id`, `name`, `color`, `icon,omitempty`, `description,omitempty`, `sort_order`, timestamps.
- `CreateLabelInput`, `UpdateLabelInput` pointer fields.
- Repository interface untuk `List`, `Get`, `Create`, `Update`, `Delete`, `Reorder`, relation read/write hanya jika menjadi dependency yang benar.
- Reorder item: `domain.ReorderItem` sudah ada di `category.go`; jangan deklarasikan ulang type global tersebut. Gunakan type existing atau pindahkan hanya dengan pertimbangan compile dan perubahan minimal.
- Interface relation jangan duplikat kemampuan. Pilih satu: `LabelRepository.GetByTransaction/SetTransactionLabels` atau transaction repository methods. Plan lama menyebut relasi pada LabelRepository, sementara partial implementation menambahkan `GetLabels/SetLabels` ke TransactionRepository. Putuskan berdasarkan cohesive ownership dan endpoint; rekomendasi: junction relation methods berada pada `LabelsRepo`, dipakai `TransactionsRepo` via SQL transaction? Hindari repository saling konstruksi tanpa kebutuhan. Untuk atomicity create/update transaction dengan junction, implementasi relation writes dalam `TransactionsRepo` memakai db yang sama dalam transaction; expose read IDs to transaction scan through helper there. `LabelsRepo.GetByTransaction` only needed by dedicated route/UI, not mandatory.
- Do not add unused interfaces/methods just because old plan listed them.

### 3.6 Transaction update atomicity

- Create transaction row + label links dalam satu SQL transaction.
- Update scalar fields + junction replacement dalam satu SQL transaction bila keduanya disertakan. Jangan commit transaksi row lalu gagal set links dan meninggalkan partial state.
- Update label-only harus tetap berfungsi walau tidak ada scalar field updates.
- Empty scalar updates + absent label IDs should return current item (`Get`) without issuing malformed `UPDATE`.
- Build `WHERE id = ?` args exactly once. Cegah bug duplicate `id` arg yang muncul sebagai risiko pada perubahan parsial sebelumnya.
- Ensure transaction exists and is not soft-deleted before replacing links; return not-found semantics consistently. Do not recreate orphan links.
- Use rollback on every error; defer rollback pattern can simplify safe rollback, commit once after all writes.
- `SetLabels` should validate all IDs before deletion or rely on FK errors inside transaction; FK failure must rollback deleted old links, preserving prior set.
- Deduplicate incoming IDs or reject duplicate input. Prefer de-duplicate in service/repo deterministically while preserving no semantic duplicates; test repeated IDs.

---

## 4. Phase 0 — Workspace audit and decisions

### Objective
Establish accurate current state so implementation does not overwrite unrelated work or assume missing components exist.

### Actions
1. Run `git status --short`; record user changes.
2. Inspect `git diff --` relevant modified files; inspect untracked partial files directly.
3. Confirm exact migration list and current latest version. Read `db.go`, all migration headers/content relevant to transaction schema and labels.
4. Confirm migration order behavior for duplicate version numbers; decide migration number strategy before schema edits.
5. Inspect transaction schema (migration creating transactions), all usages of `Transaction.Labels`, transaction HTTP DTOs, and frontend create/update callsites.
6. Search all implementations of `domain.TransactionRepository`, `CreateTransactionInput`, `UpdateTransactionInput`, JSON request types, test fakes.
7. Inspect `AuthMiddleware` signature and route mounting. Record concrete labels mounting syntax.
8. Inspect `frontend/src/widgets/transactions/TransactionsView.vue` or parent component to identify where accounts/categories/forms are populated; determine if labels should be passed via prop or picker reads store itself. Prefer one clear source of labels.
9. Inspect `BaseModal`, `BaseInput`, `BaseButton`, `BaseSelect`, project combobox/popover patterns and `@lucide/vue` exports. The plan's `BaseCombobox` is unverified: don't assume it exists.
10. Inspect DESIGN sections for nav, cards, form input, badges/chips, responsive, modal, focus/accessibility, color usage.
11. Inspect existing migration tests and test setup; inspect `.gitignore` rule for `backend/**/*_test.go`. Tests may be ignored; report that, do not modify `.gitignore` autonomously.
12. Decide and document four policy questions before coding:
    - stable migration IDs / versioning;
    - legacy labels migration (no lossy backfill by default);
    - duplicate label name behavior/case sensitivity;
    - transaction update presence semantics for `label_ids`.

### Exit criteria
- No unknown API/composable/component is assumed.
- List of exact files to touch is updated from the tree at execution time.
- No user modifications are overwritten.
- Decisions above are recorded in the implementation summary or issue notes.

---

## 5. Phase 1 — Database migrations

### Objective
Persist label records and transaction-label many-to-many relationships while retaining legacy schema/data.

### Files
- Add the next safe migration file(s) under `backend/internal/db/migrations/`.
- Existing partial `015_labels.sql`, `016_transaction_labels.sql`: inspect and either keep, amend only if unshipped, or supersede with new version if already applied/released. Never edit an applied migration.

### Implementation steps
1. Select migration IDs after Phase 0 audit; account for duplicate 008 and migration runner integer version behavior.
2. Labels migration:
   - create `labels` table with primary key, required trimmed nonempty name enforcement where SQLite constraint can safely enforce it (application validation still required), required color default, optional icon/description, sort order, timestamps;
   - unique name rule;
   - add only indexes that are not redundant and support actual listing/reorder/query.
3. Junction migration:
   - create `transaction_labels` with transaction/label FKs and compound primary key;
   - ensure FK `ON DELETE CASCADE` supports hard deletes only;
   - add `label_id` index only if needed for reverse lookup/delete;
   - do not cascade on soft-delete.
4. Migration legacy behavior:
   - do not drop, rename, or rewrite `transactions.labels`;
   - do not create label catalog items from split strings without a deterministic parsing/backfill plan;
   - migration must succeed on fresh DB and existing DB with currently supported prior schema.
5. Trigger/timestamp decision:
   - either implement `updated_at` in repository update statement, or use trigger proven safe on modernc SQLite with recursive triggers behavior;
   - prefer explicit update assignment unless an established project migration pattern proves trigger behavior safe.
6. Migration must be transaction-safe under migration runner. Migration runner wraps each SQL file in one transaction.
7. Do not seed sample labels in production.

### Validation
- Build a temporary disposable database through existing DB migration code (outside `backend/data/`, do not inspect user DB); assert schema and foreign keys.
- Apply migrations twice / open twice; idempotent `IF NOT EXISTS` behavior must not conflict with schema_version logic.
- Confirm fresh DB migration runner reaches selected version exactly once.
- Confirm existing migration path (test fixture / temp DB at prior schema) upgrades without data loss; do not use user's actual DB.
- Confirm duplicate version filenames cannot result in silently skipped feature migration.

### Exit criteria
- Fresh and upgrade paths pass.
- Legacy `transactions.labels` retained.
- Version choice documented and no applied migration edited.

---

## 6. Phase 2 — Backend domain

### Objective
Define typed label entities and transport-independent repository contracts without duplicating existing domain concepts.

### Files
- `backend/internal/domain/label.go` (existing partial file may be reused).
- `backend/internal/domain/transaction.go`.
- Relevant repository interface file(s), likely `repository_transaction.go`; only change labels interface if required.

### Implementation steps
1. Inspect existing `domain.ReorderItem`; reuse it.
2. Add `Label` with correct snake_case JSON serialization and timestamps.
3. Add create/update input models; pointer fields let update distinguish omitted values.
4. Add label repository interface with only methods used by handler/application path.
5. Add `LabelIDs []string` to `Transaction` JSON `label_ids` and both Create/Update inputs; retain legacy `Labels []string` unchanged.
6. Decide `LabelIDs` empty-vs-omitted behavior at the HTTP DTO boundary. Domain update input may need pointer-to-slice to preserve presence, or an explicit boolean; ordinary nil/empty semantics must be tested.
7. Do not make all `TransactionRepository` implementations require new methods unless repository operations truly belong there. Search before interface change.
8. Keep domain free of SQL/chis/HTTP dependencies.

### Validation
- `gofmt` changed Go files.
- `go test ./...` and `go build ./...` after each interface change.
- Search for duplicate `ReorderItem` type names.

### Exit criteria
- Contracts compile and clearly preserve legacy field.
- Presence semantics for `label_ids` explicit.

---

## 7. Phase 3 — SQLite repositories

### Objective
Implement reliable CRUD, ordering, relation reads, and atomic transaction writes.

### Files
- `backend/internal/repository/sqlite/labels.go` (existing partial file).
- `backend/internal/repository/sqlite/transactions.go`.
- Optional tests under backend sqlite tests, subject to ignored test-file caveat.

### 7.1 Labels repository

Implement:
- `List(ctx, filters?)`: stable deterministic sort `sort_order, name`; optional search only if API contract includes it.
- `Get(ctx,id)`: `nil,nil` only for `sql.ErrNoRows`, other errors propagated.
- `Create`: generate ID via existing `genID()`, validate/default values at handler boundary and repository where needed; do not ignore failed read of MAX sort order; ensure SQL insert error returned.
- `Update`: allow only supplied fields; update timestamp; no-op update should return current record; missing ID must not return a fake empty label.
- `Delete`: delete label and let FK cascade remove junction links. Prefer single `DELETE`; don't do separate non-atomic manual cleanup then label delete. Distinguish missing ID if API semantics require 404 (use rows affected).
- `Reorder`: execute all changes in one transaction; validate IDs and duplicate IDs, all-or-nothing; avoid silently ignoring unknown IDs. Ordering payload should be validated before any mutation.
- relation operations only when required; avoid method implementation copied from a draft without actual callsite.

### 7.2 Transaction repository

1. Update transaction request/row mappings only after verifying concrete transaction DB schema.
2. Preserve legacy label column read/write behavior for legacy API callers.
3. Retrieve `label_ids` for `List`, `Get`, `Search`, and any path that returns transactions. Options:
   - collect IDs in one additional query for all transaction IDs in result (bounded list), or
   - aggregate using SQLite JSON functions only if guaranteed available in modernc SQLite and tests; do not introduce N+1 query for each result by default.
4. Guarantee `label_ids` serializes as `[]` or omitted consistently; choose empty array for current transactions if API/UI expects array.
5. Create:
   - begin SQL transaction;
   - insert base transaction preserving legacy `labels` semantics;
   - validate/insert distinct IDs into junction;
   - commit only if all succeed;
   - fetch final transaction after commit through canonical `Get`.
6. Update:
   - support label-only request;
   - update scalar transaction fields and label links in same SQL transaction;
   - append transaction ID to SQL args exactly once;
   - do not early-return before applying `LabelIDs`;
   - absent IDs means unchanged, explicit empty means clear all, nonempty means exact replacement;
   - verify transaction exists, not soft-deleted;
   - FK errors roll back link replacement, preserving old links;
   - fetch result after commit.
7. Delete:
   - soft delete remains soft delete; junction links stay so potential restore retains labels;
   - hard delete cascade is handled by FK.
8. Batch operations:
   - existing `BatchRelabel` modifies legacy names; do not reinterpret its current `Labels []string` as IDs;
   - decide whether batch relabel should eventually set junction; not in scope unless API already exposes canonical label IDs. Avoid silent semantic change.
9. Scan implementation:
   - no extra SQL queries while `rows` remains open if driver connection behavior might deadlock; close result rows before follow-up query or use collected IDs after scanning.
   - stable row order and empty-result behavior.
10. SQL parameters must stay parameterized. Whitelist dynamic order fields where relevant.

### Tests
- CRUD labels: default color, explicit color/icon/description/sort order, timestamps, duplicate names.
- Update only selected fields and no-op update.
- Delete label with linked transactions: relation disappears and transaction remains.
- Reorder valid list, unknown ID, duplicate item, rollback on failure.
- Create transaction with two labels; read via get/list/search and assert IDs.
- Update with label IDs absent preserves set; explicit `[]` clears; new IDs replaces.
- Invalid label ID fails without inserting transaction or modifying prior link set.
- Legacy labels string is unchanged and returned alongside `label_ids`.
- Soft-delete retains junction links; hard-delete cascade if hard-delete path exists.
- Ensure amount conversion behavior unchanged; do not touch money conversion except necessary.

### Exit criteria
- No partial writes on errors.
- All transaction read paths return canonical IDs.
- Legacy behavior is tested.

---

## 8. Phase 4 — HTTP handlers and route wiring

### Objective
Expose validated labels API and authenticated routes.

### Files
- `backend/internal/delivery/http/labels.go` (new).
- `backend/cmd/server/main.go`.
- Relevant auth wiring only as required for new route.

### Handler implementation
1. Constructor receives `domain.LabelRepository`; transaction repo only if an actual endpoint requires it.
2. `Routes()` registers static `/reorder` before `/{id}`.
3. Implement List/Get/Create/Update/Delete/Reorder.
4. Validate:
   - JSON parse;
   - trimmed nonempty name (max length);
   - color strict hex format;
   - icon/description size bounds;
   - sort order valid integer (already typed; reject negative if product says so);
   - reorder items non-nil and unique IDs;
   - update must not accept empty name or invalid color;
   - unknown JSON fields policy should match existing handlers; do not invent global strictness.
5. Error mapping:
   - not found -> 404;
   - duplicate unique name -> 409 (identify SQLite constraint errors safely without leaking SQL text; use driver error type if available, otherwise carefully map only known constraint conditions);
   - malformed/validation -> 400;
   - internal -> generic 500.
6. Return JSON consistent with selected API contracts; don't return nil array for empty list.
7. Do not include raw error internals in body or logs with user-entered sensitive fields.
8. Do not create the optional `GET /transactions/{id}/labels` route unless frontend actually needs it; main transaction representation should carry IDs.

### Wiring
1. Instantiate `labelsRepo := sqlite.NewLabelsRepo(database)` alongside current repositories.
2. Construct `LabelsHandler` and mount under `/api/v1/labels`.
3. Put route behind `auth.AuthMiddleware(cfg)` using current chi pattern. Verify middleware response and cookie/token state in an isolated handler test if feasible.
4. Ensure route ordering/mount path actually serves both `/api/v1/labels` and `/api/v1/labels/` expected by axios.
5. Don't imply this resolves unauthenticated categories/transactions/etc. known gap.

### Validation
- `go test ./...`, `go build ./...` from backend.
- Handler tests: status/body for valid/invalid/duplicate/not-found/reorder.
- Route test confirms unauthenticated labels request rejected and authenticated path reaches handler, using existing auth test utilities without secrets.
- Manual API checks only against local disposable database / local server; don't print credentials or tunnel URLs.

### Exit criteria
- All endpoints implemented, errors stable, endpoint auth covered.

---

## 9. Phase 5 — Frontend types/API/store

### Objective
Create typed frontend access to labels using existing shared client and reactive state patterns.

### Files
- `frontend/src/types/label.ts`.
- `frontend/src/types/transaction.ts`.
- `frontend/src/types/index.ts`.
- `frontend/src/api/labels.ts`.
- `frontend/src/stores/labels.ts`.

### Types
1. `Label`: ID, name, color, optional icon/description, sort_order, timestamps.
2. `LabelFilters`: only include query fields the backend supports. Plan draft `q/sort/order` should not cause API to imply unsupported filtering.
3. Request types for create/update/reorder.
4. `Transaction.label_ids?: string[]`; preserve existing `labels?: string` type unless API contract is deliberately changed with all callers updated.
5. Create/update transaction request types should include `label_ids`; avoid `Partial<Transaction>` mismatch if existing code relies on Partial; improve type only where necessary.
6. Ensure optional vs nullable semantics match JSON decoder on Go backend.

### API client
- Import `api` from `@/api/client`; never hard-code backend URL.
- Methods: list, get, create, update, delete, reorder; optional label-by-transaction endpoint only if planned.
- Use `/api/v1` paths; choose slash convention consistently with current service.
- Query serialization must omit blank filters or use axios params correctly.
- Return typed response payloads, but consumers should read axios `.data` according to existing store convention.

### Pinia store
- `defineStore('labels', ...)`, refs for list/loading/filters.
- Provide fetch/create/update/delete/reorder actions.
- Set loading using `try/finally` for fetch; errors should be surfaced through existing UI store or propagated consistently.
- On create/update/delete/reorder, refresh list once; avoid duplicate refresh from both store and view.
- Keep inline create errors available to picker; no silent swallowing.
- Store default list as `[]`.
- Don't add pagination if backend has no pagination.
- Decide where search filtering occurs. If backend doesn't support q, filter client side against fetched set; don't send unsupported params.
- Avoid hardcoded user-visible error strings in UI components; toast strategy follows existing store patterns. For newly added strings, use i18n where displayed as UI; store error patterns existing may be English but should not multiply hard-coded text.

### Validation
- `npm run build` after type additions.
- Minimal unit tests if test setup exists; if `frontend/tests/` is missing, add only if task includes creating test setup is justified. Plan acceptance should reflect actual tooling and disclose inability.

### Exit criteria
- TypeScript catches API/form contract mismatches.

---

## 10. Phase 6 — Label management UI widgets

### Objective
Provide reusable, localized, accessible primitives for managing labels.

### Design constraints
- Read `DESIGN.md` detailed component rules before styling.
- Follow exact `<template>`, `<script setup lang="ts">`, `<style scoped>` order.
- Style first line `@reference "tailwindcss";`, utilities through `@apply` only.
- Use only CSS variables actually declared/approved and exact design tokens. Note `DESIGN.md` palette names may not match runtime CSS variables; use `AGENTS.md` listed semantic variables, not invented variables.
- Use semantic template classes. No utility class strings.
- No inline SVG. Use installed Lucide icons.
- Focus-visible, button labels, keyboard behavior, touch targets >=44px, contrast, don't communicate label only by color.
- Japanese and Indonesian long strings need wrapping/truncation behavior.
- At 360px, avoid overflow and horizontal scroll.

### 10.1 `frontend/src/widgets/labels/LabelCard.vue`
Responsibilities:
- Display label color indicator, optional icon, name, optional description if design permits.
- Provide edit/delete actions via real buttons; labels via `t()`.
- Show label name text even with color dot; accessible name should include action.
- Emit typed `edit(label)` and `delete(id)`.
- Do not calculate or display use count unless backend provides it reliably; avoid fake count.
- Avoid user-supplied icon rendered as HTML; render plain text only.
- Color swatch is decorative `aria-hidden` while visible name carries identity.

### 10.2 `LabelFormModal.vue`
Responsibilities:
- Controlled `modelValue`, optional `label` prop, typed submit payload.
- Fields: name (required), color (preset swatches + valid custom hex), optional icon plain text, description, sort order only if manual editing intended. Avoid making manual order conflict with drag reorder.
- Preset colors need to be sourced from DESIGN tokens/accessible swatches; default label DB color is indigo but current design accent differs. Decide whether product-specific label color is arbitrary user data or fixed design token set. For custom hex, validate `#RRGGBB` both client and server.
- Label every control using `for`/`id`; show validation; keyboard accessible; modal close/cancel semantics.
- New strings via `labels.*`/`common.*` translation keys.
- Do not put `placeholder="{{ t(...) }}"`; Vue binding is `:placeholder="t(...)"`.
- Reset form on open/close/edit target changes; no stale edit values on create.

### 10.3 `LabelFilters.vue`
Responsibilities:
- Search input controlled via typed `v-model` (`modelValue` + update event or named `q` model following Vue conventions).
- Search by label name/description if supported client-side; placeholder translated.
- Debounce belongs in view/store only if existing patterns support it; don't duplicate debounce at both component and page.
- Accessible label and clear-search action if implemented.

### 10.4 `LabelList.vue`
Responsibilities:
- Render empty/loading/list states.
- Emit edit/delete and reorder.
- Reorder accessible without new dependency. Options: native HTML drag events plus keyboard-accessible move up/down buttons. Required functionality says drag reorder; implement native drag and drop if robust, plus keyboard move controls for accessibility. Do not use drag-only interaction.
- Do not reorder by mutating Pinia state before server success without rollback. Either optimistic reorder with snapshot rollback on failure, or submit then refresh.
- Stable keys by label ID.
- Empty state string and CTA are localized.
- Keep sorting stable; do not allow manual reorder to conflict with active search. Disable drag reorder while filtered, or define reorder payload as entire ordering and ensure filtered ordering doesn't overwrite hidden labels.

### 10.5 Barrel `widgets/labels/index.ts`
- Export only public components used outside the folder, matching neighbors.
- Ensure no circular imports.

### Validation
- Component props/events compile under vue-tsc.
- Manual visual checks light/dark and 360px.
- Keyboard: tab, enter/space, reorder without pointer, escape/cancel modal.

---

## 11. Phase 7 — Labels page and route

### Objective
Add authenticated Labels management page and sidebar destination.

### Files
- `frontend/src/views/LabelsView.vue`.
- `frontend/src/router/index.ts`.
- Optional `frontend/src/widgets/labels/index.ts`.

### Page behavior
1. Page header with `labels.title`, add button `labels.addLabel`.
2. Search/filter control.
3. List renders loading, empty, populated, and error states.
4. Create modal and edit modal should be one modal instance driven by `editingLabel` rather than duplicated instances unless there is a concrete reason. Prevent concurrent edit/create modals.
5. Delete flow uses accessible confirmation modal or existing confirmation pattern. Native `confirm()` may be retained only if existing convention/use and localized message is correct; do not add untranslated prompt.
6. Delete confirmation includes label name; if number of linked transactions is claimed, get a trustworthy count. Otherwise don't show unsupported `usedInTransactions` message. FK cleanup is automatic; user confirmation must clearly explain deleted label associations if usage count is known.
7. On route mount, fetch labels. On search, apply backend filter if implemented or local filtering. On mutation, refresh list once.
8. Error behavior consistent with UI store/API conventions.
9. Route `/labels`, name `Labels`, lazy import, `meta.requiresAuth: true`.
10. Verify active sidebar state uses exact `/labels` and does not mark unrelated route active due to prefix matching.

### Validation
- Route resolves and unauthenticated user is redirected by existing guard.
- Add/edit/delete/reorder/search smoke path.
- Empty collection produces `[]`, not blank/error UI.

---

## 12. Phase 8 — Label picker and transaction integration

### Objective
Make label selection a reliable multi-select workflow in transaction create/edit.

### Files
- `frontend/src/components/common/LabelPicker.vue`.
- `frontend/src/widgets/transactions/TransactionFormModal.vue` (actual form path).
- `frontend/src/views/transactions/TransactionsView.vue` or parent owning modal props, only if labels are passed from parent.
- `frontend/src/stores/transactions.ts` / type file as needed.

### 12.1 Picker interaction contract

- Props: available `Label[]` or internal store access; selected IDs via `v-model:labelIds`; disabled/loading optional only if needed; translated placeholder should be a real prop/string binding.
- Emits selected ID array; never emits labels names in place of IDs.
- Selected chips show name and optionally color dot; remove button has accessible label including label name and min 44px target.
- Dropdown search filters available labels by name (and optionally description).
- Clicking/keyboard-selecting option adds ID once.
- Selection order can be stable by selection order or label sort order; choose and test.
- Dropdown closes after select only if expected UX; keyboard ArrowUp/Down/Enter/Escape behavior specified.
- Inline create: when query is nonempty and no exact existing match, show localized create action, use trim, prevent blank names, avoid duplicate names (client precheck is convenience only; server remains source of truth), call store/API create, then append returned ID and refresh options.
- On server duplicate conflict, show localized/usable message and select existing match only if lookup is unambiguous; otherwise show error, no silent selection.
- No `BaseCombobox` assumption; implement a small accessible combobox using native input/button/list semantics and ARIA roles only if necessary. Prefer native datalist only if it supports create and multi-select requirements; likely does not.
- Do not paste inline SVG. Use lucide icon components if any.
- Scope styles correctly, no template utility classes.

### 12.2 Transaction form data flow

1. Replace legacy `BaseInput v-model=form.labels` in actual `TransactionFormModal.vue` with LabelPicker.
2. Add `label_ids: string[]` to form data, default `[]`.
3. On edit, initialize from `transaction.label_ids || []`; preserve legacy string field separately so editing unrelated fields does not destroy `labels`.
4. If a transaction has legacy `labels` but no IDs, present legacy names as readonly compatibility display if needed; do not pretend they map to selectable new labels. Do not replace legacy string with `''` on submit if server contract expects preserve.
5. On submit, emit `label_ids` explicitly. `[]` must mean clear labels on an edited transaction. Ensure create no labels sends either `[]` or omitted according to contract; both acceptable only if create semantics are clear.
6. Parent modal must load label list at a suitable boundary. Avoid fetching labels on every keystroke or per transaction row. Options: store fetch on modal opening, or parent view loads labels once. Prefer minimal requests and refresh after inline create.
7. Keep existing fields/type compatibility unchanged.
8. Transaction create/update store must pass `label_ids` without serialization transformation. Keep `labels` legacy field if needed.
9. On edit, verify array is preserved through server response and form state.
10. Form labels should be included in duplicate transaction flow if duplication exists; inspect callsites and add IDs only if expected duplicate semantics carry metadata.

### Validation
- New transaction with 0/1/multiple labels.
- Existing transaction labels loaded, removed, added, saved.
- Explicit clear persists as no junction rows.
- Inline create persists label then transaction relation; failure creating label leaves transaction unsaved or allows retry without losing form state.
- Keyboard and 360px checks.

---

## 13. Phase 9 — Sidebar grouping and icon

### Objective
Expose Labels in navigation with requested Organization/Insights structure without regressing existing routes.

### Files
- `frontend/src/components/layout/Sidebar.vue`.
- `frontend/src/components/icons/LabelsIcon.vue`.
- `frontend/src/locales/en.ts`, `ja.ts`, `id.ts`.

### Structure
- Primary: Dashboard, Transactions, Calendar.
- Analytics expandable group.
- Organization heading and items: Accounts, Categories, Labels, Scheduled.
- Insights heading and items: Reports, Shared Access.
- Settings existing position/behavior retained.

### Implementation constraints
1. Add icon Vue component using installed `Tag` from `@lucide/vue` after verifying export; follow existing icon component style.
2. Keep icon size/color behavior aligned with sidebar's `.nav-icon` and semantic token system.
3. Section titles `nav.sectionOrganization`, `nav.sectionInsights` translated.
4. Decide whether bottom nav arrays remain two arrays or become explicit section arrays. Minimum change preferred, but headings must be semantically grouped and repeated correctly.
5. Use heading elements or accessible `<span>` with section semantics; do not add fake clickable controls.
6. Set `aria-current="page"` via router-link if Vue Router isn't already doing so; preserve active state. Avoid `startsWith('/labels')` matching accidental deeper paths unless route intended.
7. Mobile sidebar: selecting item closes sidebar as existing links do.
8. Avoid duplicating route items in separate arrays.
9. Ensure long translated Japanese/Indonesian titles fit 360px and do not wrap awkwardly/cut off.

### Validation
- Navigate each section item on desktop/mobile.
- All links active correctly; analytics expansion remains working.
- Heading contrast and spacing use design tokens.

---

## 14. Phase 10 — i18n

### Objective
Provide complete localized UI for every new label control, list state, navigation heading, validation, and confirmation.

### Files
- `frontend/src/locales/en.ts`.
- `frontend/src/locales/ja.ts`.
- `frontend/src/locales/id.ts`.

### Proposed key inventory (final names may match existing conventions)

`nav`:
- `labels`
- `sectionOrganization`
- `sectionInsights`

`labels`:
- `title`
- `label`
- `addLabel`
- `editLabel`
- `deleteLabel`
- `name`
- `color`
- `icon`
- `description`
- `sortOrder` only if exposed
- `searchPlaceholder`
- `searchOrCreate`
- `createNew` interpolation `{{name}}`
- `noLabels`
- `noResults`
- `deleteConfirm` interpolation `{{name}}`
- `deleteImpact` only if linked count is available and truthful
- `legacyLabelNotice` only if UI surfaces it
- `invalidColor`
- `nameRequired`
- `nameTooLong`
- `createFailed`
- `updateFailed`
- `deleteFailed`
- `reorderFailed`
- `duplicateName`
- `removeLabel` interpolation `{{name}}`
- `moveUp` / `moveDown` for keyboard-accessible reorder
- `loading`
- `createInline`
- `noMatchingLabels`
- `selectedCount` if announced to screen reader

`common` additions only if missing:
- `create`
- `clear`
- `moveUp`/`moveDown` should preferably live in feature group if meaning is specific.

Rules:
1. Inventory actual keys used by components; no unused speculative translation keys.
2. Ensure same nested keys in all three locale objects.
3. Japanese translation uses natural product language; Indonesian is natural and concise; no placeholder remains untranslated except interpolation tokens.
4. Include long-text testing; no hard-coded HTML attributes that display English to users.
5. Existing English-only strings in legacy Category components are out of scope unless directly needed by this feature.

### Validation
- Script/check to compare new subtree keys across files if no i18n test exists.
- `npm run build` catches key typing only if setup supports it; manually inspect all `t('labels...')`/`t('nav...')` references.

---

## 15. Phase 11 — Optional analytics investigation

### Objective
Assess, but do not silently expand scope into, spending analytics by label.

### Actions
1. Inspect `frontend/src/widgets/analytics/SpendingByLabelWidget.vue` if present; plan says it may exist, actual glob inspection must confirm.
2. Check if component is used or orphaned; check its data source/API endpoint and whether live backend serves it.
3. Check whether chart expects category labels, string names, IDs, or legacy `transactions.labels`.
4. If endpoint/contract is not live, do not wire to nonexistent routes. Record gap.
5. Only update widget if it is already reachable and minimal ID migration can be made without broad analytics endpoint work.
6. No new chart, endpoint, migration, or analytics backend in this implementation unless explicitly approved after discovery.

### Exit criteria
- Explicitly mark analytics integration as included, updated, or deferred with evidence.

---

## 16. Phase 12 — Tests and verification

### Backend required commands
Run from `backend/`:

```bash
gofmt -w <changed Go files>
go build ./...
go test ./...
```

Tests should cover database migration, repository CRUD, relationship atomicity, transaction read path, HTTP input/status/auth. Before creating tests, inspect `.gitignore`: current project notes that `backend/**/*_test.go` is ignored. Do not modify `.gitignore` unless explicitly approved; if files are ignored, use an allowed test file location if architecture permits or report that test files would be untracked. Do not claim coverage from ignored files.

### Frontend required commands
Run from `frontend/`:

```bash
npm run build
npm run test:run
```

- Build includes `vue-tsc`; mandatory.
- Check whether vitest tests actually exist and what the configured `test:run` command does; no test folder was noted previously.
- Add tests only where the existing harness supports them; do not add a new testing dependency or scaffold solely for a nominal test unless approved.
- If no automated frontend tests exist, perform and report manual test plan; do not say frontend tests passed if no tests ran.

### Manual UI checks required by `AGENTS.md`
Check at minimum:
- Light theme and dark theme, controlled by app's settings mechanism.
- 360px viewport width.
- Long Japanese and Indonesian copy.
- Labels list with 0 labels, 1 label, many labels, long name, long description.
- Custom color validity and color contrast; chip text remains readable.
- Search no result and exact match.
- Create/edit/delete/reorder.
- Keyboard-only operation, focus visible, screen reader labels.
- Transaction create with no labels and multiple labels.
- Transaction edit preserving IDs and clearing all labels.
- Legacy transaction string labels remain visible and unchanged.
- Inline create duplicate conflict and network failure.
- Modal overflow and sidebar group scrolling on mobile.
- Hide-amounts setting ON: this feature contains no money values, but verify transaction form context remains unaffected.
- PWA navigation does not introduce stale API caching; no custom caching added.

### Integration/API verification
Against isolated local test DB only:
1. Create label.
2. List label; verify sort order and JSON shape.
3. Create transaction with two label IDs.
4. Read transaction and verify both `label_ids`, plus legacy `labels` shape.
5. Update transaction omitting label_ids; verify existing relation unchanged.
6. Update with empty `label_ids`; verify relation cleared.
7. Update with invalid ID; verify atomic failure and old relation preserved.
8. Delete label; verify its links removed, transaction remains.
9. Reorder labels and reload to verify persistence.
10. Unauthenticated request to labels API returns auth rejection.
11. Ensure no access credentials, secrets, data path, or tunnel URL are emitted.

### Exit criteria
- All required commands report their real status.
- Known unrelated warnings are separated from failures.
- Manual checks either completed or clearly listed as not performed.

---

## 17. Detailed file plan

### New backend files

| File | Scope |
|---|---|
| `backend/internal/db/migrations/<safe-version>_labels.sql` | Add labels table/indexes, not sample rows |
| `backend/internal/db/migrations/<safe-version>_transaction_labels.sql` | Add many-to-many junction, constraints and needed indexes |
| `backend/internal/domain/label.go` | Label entity, create/update inputs, repository interface; partial untracked file exists to inspect |
| `backend/internal/repository/sqlite/labels.go` | SQL implementation; partial untracked file exists to inspect |
| `backend/internal/delivery/http/labels.go` | Validated CRUD/reorder handlers |

### Modified backend files

| File | Scope |
|---|---|
| `backend/internal/domain/transaction.go` | `LabelIDs`; preserve legacy Labels |
| `backend/internal/domain/repository_transaction.go` | Change only if necessary; avoid unused interface expansion |
| `backend/internal/repository/sqlite/transactions.go` | Atomic junction writes, IDs returned in all read methods, legacy column preserved |
| `backend/cmd/server/main.go` | Repository/handler wiring; mount labels behind auth middleware |
| Transaction HTTP DTO/handler files | Decode `label_ids` with omitted-vs-empty presence semantics; inspect exact paths |
| Migration/repository tests | Add only if valid test locations are tracked/allowed |

### New frontend files

| File | Scope |
|---|---|
| `frontend/src/types/label.ts` | Types for model/filter/requests |
| `frontend/src/api/labels.ts` | API using shared axios client |
| `frontend/src/stores/labels.ts` | Pinia store CRUD/load/reorder |
| `frontend/src/widgets/labels/LabelCard.vue` | Display/action card |
| `frontend/src/widgets/labels/LabelFormModal.vue` | Accessible create/edit form |
| `frontend/src/widgets/labels/LabelFilters.vue` | Search control |
| `frontend/src/widgets/labels/LabelList.vue` | Empty/loading/list/reorder |
| `frontend/src/widgets/labels/index.ts` | Barrel exports |
| `frontend/src/views/LabelsView.vue` | Page composition/CRUD state |
| `frontend/src/components/common/LabelPicker.vue` | Accessible multi-select and inline create |
| `frontend/src/components/icons/LabelsIcon.vue` | Lucide label icon wrapper |

### Modified frontend files

| File | Scope |
|---|---|
| `frontend/src/types/transaction.ts` | `label_ids` while preserving `labels` contract |
| `frontend/src/types/index.ts` | Export label types |
| `frontend/src/stores/transactions.ts` | Forward IDs unchanged; preserve current API behavior |
| `frontend/src/widgets/transactions/TransactionFormModal.vue` | Picker integration; preserve legacy values |
| Transaction parent view (verify actual path) | Fetch/pass available label list if picker uses prop |
| `frontend/src/components/layout/Sidebar.vue` | Organization/Insights groups and Labels link |
| `frontend/src/router/index.ts` | Authenticated lazy `/labels` route |
| `frontend/src/locales/en.ts` | All labels/nav translations |
| `frontend/src/locales/ja.ts` | All labels/nav translations |
| `frontend/src/locales/id.ts` | All labels/nav translations |

Do not touch unrelated analytics, global design tokens, auth for existing routes, README, or `.gitignore` without a concrete requirement and approval.

---

## 18. Work sequencing and dependency graph

1. **Audit / decisions** → blocks all implementation.
2. **Migrations** → blocks repository behavior and integration tests.
3. **Domain contract** → blocks repository, handler, API types.
4. **SQLite repositories** → blocks handler end-to-end semantics.
5. **HTTP handler + auth wiring** → blocks frontend integration.
6. **Frontend types/API/store** → blocks LabelsView and picker.
7. **Management widgets + LabelsView + route** → blocks management flow acceptance.
8. **Picker + transaction backend DTO** → blocks transaction flow acceptance.
9. **Sidebar + icon + translations** → required navigation and localization.
10. **Optional analytics audit** → independent; defer without blocking core feature.
11. **Automated + manual verification** → blocks completion claim.

Parallel work is allowed only for independent frontend types/UI after API contracts and server semantics are fixed. Avoid parallel edits to `Sidebar.vue`, transaction form/types, migration files, or i18n objects without coordination.

---

## 19. Acceptance criteria (binary)

### Backend/data
- [ ] Safe unique migration version selected; no applied migration edited; fresh DB migrates.
- [ ] Upgrade DB retains legacy `transactions.labels` contents.
- [ ] Labels CRUD returns documented JSON/statuses.
- [ ] Duplicate/invalid labels return safe client error.
- [ ] Reorder is stable, atomic, validates IDs.
- [ ] Junction enforces valid IDs and no duplicate links.
- [ ] Create transaction with label IDs is atomic.
- [ ] Update omitted IDs preserves labels; empty array clears; IDs replace set.
- [ ] Invalid label IDs do not partially mutate transaction or link set.
- [ ] Get/List/Search transaction paths return `label_ids` consistently.
- [ ] Legacy `labels` remains available and unchanged unless explicitly set by legacy caller.
- [ ] Label delete removes junction links but not transactions.
- [ ] Soft-delete behavior is defined and tested.
- [ ] Labels routes require authentication.
- [ ] `go build ./...` and `go test ./...` pass from backend.

### Frontend
- [ ] `/labels` route resolves, is auth-protected in router, and appears in sidebar.
- [ ] Sidebar has Organization and Insights headings and correct items.
- [ ] Label CRUD UI works and refreshes state.
- [ ] Search works for supported client/server behavior.
- [ ] Reorder persists, has keyboard alternative, and filtered lists do not corrupt hidden ordering.
- [ ] Color picker rejects invalid values and label name stays legible.
- [ ] Picker selects multiple IDs, removes selections, searches, inline-creates and auto-selects.
- [ ] Create transaction persists selected IDs.
- [ ] Edit transaction loads current IDs; add/remove/clear persists correctly.
- [ ] Legacy names stay displayable and are not confused with IDs.
- [ ] Every user-visible new string exists in en/ja/id.
- [ ] All new UI follows DESIGN.md and `AGENTS.md` style/accessibility rules.
- [ ] `npm run build` passes.
- [ ] `npm run test:run` result accurately reported; no false claim if tests are absent.

### Manual quality
- [ ] Light and dark theme checked.
- [ ] 360px width checked.
- [ ] Long Japanese and Indonesian copy checked.
- [ ] Keyboard/focus/accessibility behavior checked.
- [ ] Hide-amounts ON checked in transaction context.
- [ ] No new stale financial data caching introduced.

---

## 20. Risks and mitigations

| Risk | Effect | Mitigation |
|---|---|---|
| Duplicate migration version or runner ordering | Fresh DB misses table | Audit runner and names; unique migration IDs; test fresh and upgrade paths |
| Migration version already applied in user DB | New schema never executes | Never edit applied migration; select next version; treat existing data carefully |
| Legacy comma-delimited values are ambiguous | Data corruption/backfill errors | Keep legacy field, do not auto-split/backfill without deterministic requirements |
| Update JSON cannot distinguish omitted/empty list | Unable to clear links or accidental clears | Use presence-aware request DTO; integration test all 3 states |
| Transaction scalar/link updates split | Partial financial metadata writes | One database transaction for both |
| Invalid ID fails after deleting prior links | Labels silently lost | FK failure inside same transaction and rollback test |
| N+1 label fetch in transaction lists | Slow mobile UI | Batch-fetch relation IDs or validated aggregate query |
| Extra query while rows open | SQLite pool deadlock/locking | Collect IDs, close rows, then fetch relations |
| SQL sort/filter injection | Database compromise | Whitelist columns/directions; parameterize values |
| Duplicate label names / casing confusion | UX/API inconsistencies | Explicit case policy + unique constraint conflict mapping |
| Custom colors have poor contrast | Inaccessible chips | Ensure name text remains visible; contrast checks; never color-only identity |
| Drag-only reorder | Keyboard/touch exclusion | Add move up/down controls, focus handling, native drag only as enhancement |
| Search + reorder payload drops hidden labels | Corrupt order | Disable reorder while filtered or reorder full canonical list |
| Inline create network failure | Form state loss or duplicate labels | Keep selection state, show localized error, retry flow |
| Stale label options | Selection of deleted ID | Refresh options when modal opens and after mutations; server FK validation |
| Route auth gap | Public data exposure | New labels route behind auth; explicitly report existing data routes still open |
| i18n incomplete | UI English fallback or missing text | Key inventory audit across all locale files |
| New component violates CSS rule | Inconsistent UI/build issue | Style block starts with `@reference`; semantic selectors; DESIGN tokens only |
| Test files ignored by git | Tests not delivered | Inspect `.gitignore`; report; do not change ignore policy without approval |
| Existing frontend tests absent | No automated UI behavior proof | Build + manual checklist; don't fabricate test pass |
| Optional analytics endpoint is not mounted | Broken chart | Audit and defer if no live path |
| Current partial work contains latent bug | Existing transaction update corruption | Inspect diff; correct ID args and verify label-only update before proceeding |

---

## 21. Implementation reporting format

At completion, summarize in concise Indonesian:
1. **Berubah:** major files/features, no exhaustive file dump unless requested.
2. **Verifikasi:** exact commands and pass/fail; mention manual checks actually performed.
3. **Belum diverifikasi:** browser visual checks/tests unavailable, if applicable.
4. **Gap diketahui:** auth gap on existing routes, migration/version caveats, ignored test files, or analytics deferred.
5. Never claim full acceptance if any checkbox remains unchecked.

---

## 22. Stop conditions requiring clarification

Pause and ask before proceeding if:
- Migration version cannot be selected safely because prior schema/version history is unknown.
- Product requires automatic backfill of legacy comma-delimited values but semantics are ambiguous.
- Product requires label filtering, usage counts, or analytics beyond existing live APIs.
- Existing auth middleware cannot be applied to labels route without unexpectedly changing other route behavior.
- The existing user worktree has overlapping changes that cannot be safely preserved.
- Label names are expected to be case-insensitive but database unique semantics remain case-sensitive.
- User asks to migrate/inspect their real database or expose live tunnel endpoints.

Otherwise use the minimal safe interpretation in this plan and document it.
