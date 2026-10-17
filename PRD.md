# PRD — Kiroku: Personal Finance Tracker

Versi 1.0 (draft) · disusun dari analisis kode monorepo `kiroku` · 15 Oktober 2026

## 1. Ringkasan Eksekutif

Kiroku adalah aplikasi pencatatan keuangan pribadi **self-hosted, single-user, multi-akun dan multi-mata-uang**, dengan taksonomi kategori rumah tangga Jepang (食費, 交通費, 住居費, dst.) dan mata uang dasar JPY. Aplikasi berjalan sebagai monorepo: backend Go + SQLite dan frontend Vue 3 PWA. Seluruh data tersimpan di mesin milik pengguna; frontend dapat di-host terpisah (Vercel) dan menjangkau backend lewat Cloudflare Tunnel.

Dokumen ini mendefinisikan **produk target** Kiroku v1.0 dan membandingkannya dengan **kondisi kode saat ini** (Bagian 11), sehingga bisa dipakai sebagai acuan tunggal untuk stabilisasi, penyelesaian fitur, dan penugasan ke coding agent.

## 2. Latar Belakang dan Masalah

- Aplikasi keuangan komersial menyimpan data di cloud pihak ketiga, tidak mendukung kategori gaya Jepang, dan jarang memadukan multi-currency dengan laporan cashflow ke depan.
- Pencatatan manual cepat membosankan; pengguna butuh input super cepat (FAB, template, transaksi berulang, scan struk).
- Pengguna ingin satu sumber kebenaran yang bisa di-backup, diekspor, dan dibagikan sebagian (read-only) tanpa memberi akses penuh.

## 3. Tujuan, Non-Tujuan, dan Metrik

### 3.1 Tujuan

- G1 — Mencatat income, expense, dan transfer antar akun dalam hitungan detik di ponsel (PWA).
- G2 — Saldo, laporan, dan analitik selalu konsisten lintas akun dan mata uang (dikonversi ke mata uang dasar).
- G3 — Otomatisasi: transaksi berulang (RRULE), template, dan scan struk (OCR + AI) menghasilkan draf transaksi yang tinggal dikonfirmasi.
- G4 — Data aman dan portabel: backup terenkripsi, restore, ekspor/impor CSV.
- G5 — Dapat dijalankan penuh oleh satu orang di perangkat sendiri (termasuk Termux/PM2) dengan satu perintah.

### 3.2 Non-Tujuan (v1.0)

- Multi-user/multi-tenant dengan registrasi publik (v1.0 hanya satu pengguna).
- Integrasi bank/open banking dan sinkronisasi otomatis mutasi rekening.
- Aplikasi native iOS/Android (cukup PWA).
- Fitur pajak, invoice, atau akuntansi double-entry penuh.

### 3.3 Metrik Sukses

- Waktu mencatat expense sederhana ≤ 10 detik dan ≤ 4 ketukan dari dashboard.
- p95 respons API baca < 300 ms pada ≥ 50.000 transaksi di perangkat kelas ponsel/SBC.
- Akurasi draf scan struk: nominal dan tanggal benar pada ≥ 85% struk uji.
- 0 selisih antara saldo akun, total laporan, dan total analitik pada dataset uji (reconciliation test lulus).
- Backup harian sukses ≥ 99% dan restore teruji pada uji bulanan.

## 4. Persona dan Use Case

**Persona utama — pemilik Kiroku.** Individu teknis yang mengelola rekening/dompet di beberapa mata uang (mis. JPY dan IDR), ingin kontrol penuh atas data, dan memakai ponsel sebagai perangkat utama.

**Persona sekunder — penerima akses bersama.** Pasangan/akuntan yang hanya perlu melihat satu akun lewat tautan read-only berbatas waktu.

Use case utama:

- UC1: Catat pengeluaran di toko lalu foto struk; konfirmasi draf hasil OCR.
- UC2: Atur gaji bulanan dan tagihan sewa sebagai transaksi berulang.
- UC3: Transfer antar akun berbeda mata uang dan lihat dampaknya ke saldo.
- UC4: Tinjau pengeluaran per kategori/label/sifat (must/have/want) dan bandingkan dengan anggaran.
- UC5: Lihat proyeksi cashflow dan saldo ke depan.
- UC6: Bagikan riwayat satu akun ke orang lain, lalu cabut aksesnya.
- UC7: Backup, pindah mesin, restore.

## 5. Arsitektur dan Stack

### 5.1 Struktur Monorepo

- `backend/` — API Go: `cmd/server` (entry point), `internal/<modul>` (accounts, analytics, auth, backup, categories, config, db, fx, impexp, logger, models, receipts, reports, scheduled, settings, shared, templates, transactions, tunnel), `data/` (SQLite), `docs/` (Swagger).
- `frontend/` — Vue 3 + Vite + Pinia + Tailwind 4 + ECharts + Leaflet; `src/` berisi `views`, `widgets`, `components`, `stores`, `api`, `router`, `composables`, `locales`, `types`.
- Root: `package.json` (skrip gabungan), `start.sh` (PM2), `.agent/workflows` (4 workflow finalisasi untuk agent), `CLAUDE.md` (direktif agent).

### 5.2 Teknologi

- Backend: Go, chi router, SQLite via `modernc.org/sqlite` (tanpa CGO, WAL, foreign_keys ON), koanf (YAML + env `KIROKU_*`), zerolog, JWT HS256, Argon2, rrule-go, Swagger (swaggo).
- Frontend: Vue 3, vue-router, Pinia, axios, dayjs, ECharts, Leaflet, vite-plugin-pwa (Workbox), Vitest.
- Layanan eksternal: OCR.space (OCR), 9Router (gateway LLM kompatibel OpenAI untuk normalisasi struk), exchangerate.host (kurs), Telegram Bot (persetujuan ganti password), Upstash Redis (penemuan URL tunnel), Cloudflare Tunnel.

### 5.3 Deployment

- Dev: `npm run dev` (backend :8080, frontend :5173 dengan proxy `/api`) atau `./start.sh --dev` (PM2 + air).
- Prod: `./start.sh --prod` membangun binary statis dan menjalankannya lewat PM2; backend mengelola proses `cloudflared`, menulis URL tunnel ke file dan Upstash; frontend dibangun dan di-host di Vercel, lalu membaca URL backend dari Upstash saat runtime.
- Konfigurasi: `config.yaml` + env `KIROKU_*`; rahasia (JWT, API key, token Telegram, passphrase backup) hanya lewat env, tidak di-commit.

## 6. Model Data

Semua nominal disimpan sebagai **integer satuan terkecil mata uang** (JPY = yen, USD = sen). ID berupa string UUID (kecuali seed kategori). Timestamp disimpan UTC.

| Entitas | Field kunci | Aturan |
|---|---|---|
| accounts | id, name, currency (3 huruf kapital), type (cash/bank/ewallet/investment/liability), icon, color, archived | Tidak dihapus keras bila punya transaksi; gunakan archived |
| categories | id, parent_id, name, icon, color, kind (must/have/want), sort_order, budget_amount, budget_period (weekly/monthly/yearly) | Hierarki; tidak bisa dihapus bila ada transaksi aktif; budget_amount dipakai untuk anggaran dan goals |
| transactions | id, type (income/expense/transfer), date_time, amount, currency, from/to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status (cleared/uncleared/reconciled), place_name/lat/lng, attachments, template_id, deleted_at | Soft delete; FTS5 pada note, payee_payer, labels |
| templates | seperti transaksi tanpa tanggal + name | Dapat di-apply menjadi transaksi |
| scheduled | rrule, next_run, max_occurrences, occurrences_done, transaction_template_json, enabled | RRULE divalidasi saat simpan |
| fx_rates | base, quote, rate, date, source | PK (base, quote, date) |
| settings | key, value | theme, language (default ja), base_currency, date_format, number_format, week_start, totp_enabled |
| users | id, email, password_hash, totp_secret | v1.0 hanya satu baris |
| backups | id, filename, size, retention_tier (daily/weekly/monthly) | Retensi 7 harian / 4 mingguan / 6 bulanan |
| shares | id, account_id, token, expires_at, last_accessed_at | Token unik, dapat dicabut |
| receipt_jobs, receipt_job_events | status, stage, ocr_text, ai_raw_output, transaction_draft, transaction_id, error_message; log event per tahap | Dipakai alur scan struk |

### Aturan Bisnis Inti

- BR1 — expense wajib `from_account_id`; income wajib `to_account_id`; transfer wajib keduanya dan harus berbeda; amount > 0.
- BR2 — Saldo akun = saldo awal (jika ada) + income masuk − expense keluar ± transfer masuk/keluar; transaksi soft-deleted tidak dihitung.
- BR3 — Transfer lintas mata uang memakai kurs tanggal transaksi (atau kurs terdekat sebelumnya); kurs yang dipakai harus tercatat agar laporan dapat direproduksi.
- BR4 — Semua agregat lintas mata uang dikonversi ke `base_currency`; bila kurs tidak tersedia, tampilkan peringatan, jangan diam-diam memakai kurs 1.
- BR5 — Kategori `kind` (must/have/want) menjadi dimensi "sifat pengeluaran" di analitik.
- BR6 — Status `reconciled` membuat transaksi terkunci terhadap perubahan nominal/akun kecuali dibuka lagi secara eksplisit.

## 7. Kebutuhan Fungsional

Prioritas: **P0** wajib untuk rilis, **P1** penting, **P2** nilai tambah.

### 7.1 Autentikasi dan Keamanan (AUTH)

- AUTH-1 (P0) Login dengan password tunggal; sukses menyetel cookie httpOnly access (15 menit) dan refresh (default 7 hari), `SameSite` dan `Secure` sesuai lingkungan.
- AUTH-2 (P0) Endpoint `me` mengembalikan identitas pengguna sesi aktif; frontend memakainya untuk guard route dan overlay password.
- AUTH-3 (P0) Refresh token memperbarui access token secara transparan; frontend tidak logout pada 401 pertama sebelum mencoba refresh.
- AUTH-4 (P0) Rate limiting login per IP (mis. 5 percobaan/15 menit) dengan respons 429.
- AUTH-5 (P0) Tidak ada password default tetap; setup awal memaksa pengguna menetapkan password (atau memakai password dari env pada first-run) dan JWT secret wajib ≥ 32 karakter, aplikasi menolak start bila memakai nilai default.
- AUTH-6 (P1) Ganti password wajib verifikasi password lama dan persetujuan Telegram (YES/NO, timeout 60 detik) dengan kode permintaan unik; proses tidak boleh memblokir goroutine HTTP lebih dari batas wajar.
- AUTH-7 (P1) TOTP opsional (setup via QR, verifikasi saat login) memakai modul `totp` yang sudah ada.
- AUTH-8 (P1) Semua rute di bawah `/api/v1` kecuali login/logout/health dilindungi middleware; rute publik hanya `/shared/public/{token}`.

### 7.2 Akun (ACC)

- ACC-1 (P0) CRUD akun dengan tipe, mata uang, ikon, warna; arsip/pulihkan.
- ACC-2 (P0) Ringkasan saldo per akun dan total per mata uang, konsisten dengan BR2.
- ACC-3 (P1) Saldo awal per akun dan rekonsiliasi (tandai transaksi reconciled, bandingkan dengan saldo rekening nyata).
- ACC-4 (P1) Akun bertipe `liability` dihitung sebagai pengurang kekayaan bersih.

### 7.3 Kategori (CAT)

- CAT-1 (P0) Pohon kategori dua tingkat dengan seed taksonomi Jepang saat DB kosong.
- CAT-2 (P0) CRUD, reorder (drag), ikon/warna, atribut kind.
- CAT-3 (P0) Anggaran per kategori (nominal + periode) yang menjadi sumber halaman Budget dan Goals.
- CAT-4 (P1) Hapus kategori dengan opsi pindahkan transaksi ke kategori lain (merge).

### 7.4 Transaksi (TXN)

- TXN-1 (P0) Buat/ubah/hapus (soft) transaksi income, expense, transfer sesuai BR1.
- TXN-2 (P0) Daftar dengan paginasi, filter (tanggal, akun, kategori, tipe, status, label, nominal), dan urutan.
- TXN-3 (P0) Pencarian full-text (FTS5) pada catatan, payee/payer, label; autocomplete payee/label.
- TXN-4 (P0) Operasi batch (hapus, ubah kategori, ubah status).
- TXN-5 (P1) Field lokasi dengan peta Leaflet, kalkulator di form, garansi (`warranty_until`), lampiran.
- TXN-6 (P1) Tempat sampah: lihat dan pulihkan transaksi yang dihapus.
- TXN-7 (P0) FAB dan pintasan PWA (Add Expense/Income/Transfer) membuka form dengan tipe terisi; rute yang dituju shortcut harus ada.

### 7.5 Template dan Transaksi Terjadwal (TPL, SCH)

- TPL-1 (P0) CRUD template dan aksi "apply" yang membuat transaksi dengan tanggal hari ini (dapat diubah).
- SCH-1 (P0) Buat jadwal berbasis RRULE + template transaksi; preset (harian/mingguan/bulanan/tahunan); batas jumlah kemunculan; toggle aktif.
- SCH-2 (P0) **Eksekutor**: job latar di backend yang tiap menit/jam mengambil jadwal jatuh tempo, membuat transaksi, menaikkan `occurrences_done`, menghitung `next_run`, dan menonaktifkan jadwal bila `max_occurrences` tercapai. Harus idempoten (tidak membuat ganda saat restart) dan mengejar kemunculan yang terlewat saat backend mati.
- SCH-3 (P1) Pratinjau 5 kemunculan berikutnya sebelum menyimpan; riwayat transaksi yang dihasilkan.

### 7.6 Scan Struk (RCP)

- RCP-1 (P0) Unggah foto struk (dikompresi di klien), backend membuat job dan memprosesnya asinkron: OCR → normalisasi AI → draf transaksi.
- RCP-2 (P0) OCR memakai OCR.space dengan retry 3x; bahasa OCR dapat dikonfigurasi (minimal `eng` dan `jpn`) karena struk utama berbahasa Jepang.
- RCP-3 (P0) Normalisasi lewat 9Router dengan prompt yang menyertakan daftar kategori dan akun valid, skema ketat, dan validasi hasil (ID kategori/akun harus ada; ulangi bila tidak valid).
- RCP-4 (P0) Pelacakan job (`GET /receipts/jobs/{id}`) dengan stage dan event agar UI menampilkan progres; konfirmasi menghasilkan transaksi dan menautkan `transaction_id`.
- RCP-5 (P1) Mode mock OCR untuk uji E2E; pembersihan gambar/teks mentah setelah N hari demi privasi.
- RCP-6 (P2) Deteksi duplikat (nominal + tanggal + merchant mirip) sebelum konfirmasi.

### 7.7 Impor/Ekspor CSV (IMP)

- IMP-1 (P0) Ekspor transaksi ke CSV dengan filter.
- IMP-2 (P0) Impor dua langkah: unggah → pratinjau (pemetaan kolom, deteksi baris bermasalah, pencocokan fuzzy akun/kategori) → konfirmasi.
- IMP-3 (P1) Deteksi duplikat saat impor dan laporan hasil (dibuat/dilewati/gagal).

### 7.8 Laporan dan Analitik (RPT, ANL)

- RPT-1 (P0) Laporan pengeluaran per periode dan tren (harian/mingguan/bulanan) dengan pemilih rentang tanggal.
- ANL-1 (P0) **Balance**: tren saldo, saldo per akun, saldo per mata uang.
- ANL-2 (P0) **Cashflow**: proyeksi cashflow, tabel cashflow, waterfall, tren, perbandingan antar periode.
- ANL-3 (P0) **Spending**: per kategori, per label, per sifat (must/have/want).
- ANL-4 (P0) **Outlook**: rencana per kategori dan timeline dari transaksi terjadwal, prakiraan saldo.
- ANL-5 (P0) **Budget**: ringkasan dan tabel anggaran vs realisasi per kategori.
- ANL-6 (P0) **Goals**: progres, status (tercapai/on track/tertinggal), dan proyeksi per kategori.
- ANL-7 (P1) **Stats**: ringkasan, per kategori/mata uang/akun, top merchant, cashflow, net worth.
- ANL-8 (P1) Kalender transaksi bulanan dengan detail harian.
- ANL-9 (P0) Semua endpoint analitik wajib memakai kolom yang benar-benar ada di skema, memakai konversi kurs sesuai BR3–BR4, dan mengecualikan transaksi terhapus.

### 7.9 Kurs (FX)

- FX-1 (P0) Ambil kurs terbaru terhadap mata uang dasar saat tabel kosong dan terjadwal harian; waktu dapat dikonfigurasi lewat cron.
- FX-2 (P0) Endpoint kurs per tanggal dan konversi; bila sumber gagal, pakai kurs terakhir tersimpan dan tandai usianya.
- FX-3 (P2) Penimpaan kurs manual per transaksi.

### 7.10 Akses Bersama (SHR)

- SHR-1 (P1) Buat tautan read-only per akun dengan token acak kuat dan masa berlaku opsional; daftar tautan menampilkan token tersamarkan dan waktu akses terakhir.
- SHR-2 (P1) Cabut tautan; tautan kedaluwarsa mengembalikan 410.
- SHR-3 (P1) Tampilan publik hanya memuat saldo dan riwayat akun tersebut (tanpa lampiran, lokasi, atau data akun lain), dengan saldo dihitung memakai aturan yang sama dengan BR2 (termasuk transfer).
- SHR-4 (P1) Rate limit pada rute publik.

### 7.11 Backup dan Restore (BKP)

- BKP-1 (P0) Backup manual dan terjadwal (cron `backup_schedule`) berupa `tar.gz` berisi manifest dan snapshot DB yang konsisten (gunakan SQLite backup/`VACUUM INTO`, bukan menyalin file saat WAL aktif).
- BKP-2 (P0) Enkripsi opsional AES-GCM dengan kunci turunan KDF (scrypt/argon2) dan salt acak; bukan hash SHA-256 polos.
- BKP-3 (P0) Retensi 7 harian, 4 mingguan, 6 bulanan; berkas lama dihapus otomatis.
- BKP-4 (P0) Unduh dan restore dengan validasi manifest/versi skema serta backup otomatis sebelum restore.
- BKP-5 (P1) Path penyimpanan berasal dari `data_dir`, bukan path relatif tetap.

### 7.12 Pengaturan, i18n, PWA (SET)

- SET-1 (P0) Pengaturan: tema (system/light/dark), bahasa (id/en/ja; default ja), mata uang dasar, format tanggal/angka, awal minggu.
- SET-2 (P0) Semua teks UI tersedia di tiga lokal dengan cakupan kunci sama; tes otomatis memeriksa kunci yang hilang.
- SET-3 (P0) PWA terpasang (manifest, service worker Workbox, prompt instal); API memakai NetworkFirst dengan cache 24 jam; data finansial di cache dihapus saat logout.
- SET-4 (P1) Mode offline-read: dashboard dan daftar transaksi terakhir tetap terbaca tanpa jaringan.
- SET-5 (P2) Haptics, skeleton loading, dan toast untuk umpan balik.

### 7.13 Operasional (OPS)

- OPS-1 (P0) `/healthz`, log terstruktur, graceful shutdown.
- OPS-2 (P0) Swagger UI di `/api/docs` selalu sinkron dengan rute nyata (nama cookie, skema).
- OPS-3 (P1) Manajer tunnel Cloudflare dengan auto-restart dan publikasi URL ke Upstash; frontend menangani kegagalan pencarian URL dengan pesan jelas.
- OPS-4 (P1) Satu perintah build prod yang menghasilkan binary statis dan bundel frontend.

## 8. Permukaan API (v1)

Basis `/api/v1`, autentikasi cookie, JSON. Kecuali dicatat, semua rute butuh sesi.

| Modul | Endpoint |
|---|---|
| Auth (publik: login/logout) | POST /auth/login, POST /auth/logout, GET /auth/me, POST /auth/password |
| Accounts | GET/POST /accounts, GET /accounts/summary, GET/PUT/DELETE /accounts/{id} |
| Categories | GET/POST /categories, GET /categories/tree, PUT /categories/reorder, GET/PUT/DELETE /categories/{id} |
| Transactions | GET/POST /transactions, POST /transactions/batch, GET /transactions/search, /autocomplete, /summary, GET/PUT/DELETE /transactions/{id} |
| Templates | GET/POST /templates, GET/PUT/DELETE /templates/{id}, POST /templates/{id}/apply |
| Scheduled | GET/POST /scheduled, POST /scheduled/preset, GET/PUT/DELETE /scheduled/{id}, POST /scheduled/{id}/toggle |
| Receipts | POST /receipts/scan, GET /receipts/jobs/{id}, POST /receipts/jobs/{id}/confirm |
| Import/Export | GET /export/csv, POST /import/csv, POST /import/csv/confirm |
| Reports | GET /reports/spending, GET /reports/trend |
| Analytics | /analytics/cashflow/projection, /calendar, /balance/{trend,accounts,currencies}, /outlook/{category,timeline,forecast}, /budget/{summary,vs-actual}, /goals/{progress,projection}, /stats/{summary,by-category,by-currency,by-account,top-merchants,cashflow,net-worth} |
| FX | GET /fx/rates, GET /fx/convert |
| Settings | GET/PUT /settings, GET/PUT /settings/{key} |
| Backup | GET/POST /backup, GET /backup/{id}/download, POST /backup/restore |
| Shares | POST/GET /shares-admin, DELETE /{id}; publik: GET /shared/public/{token} |

Konvensi: error berformat JSON `{error, code}` (bukan teks polos), paginasi `limit/offset`, tanggal ISO 8601, nominal integer satuan terkecil.

## 9. Pengalaman Pengguna

### 9.1 Navigasi

Dashboard, Transaksi (+ Scan Struk), Akun, Kategori, Terjadwal, Laporan, Analitik (Balance, Outlook, Cashflow, Spending, Budget, Goals), Kalender, Akses Bersama, Pengaturan, Tentang. Mobile memakai BottomNav + FAB; desktop memakai Sidebar. Tautan publik `/shared/:token` dapat dibuka tanpa login.

### 9.2 Prinsip UX

- Mobile-first, satu tangan; target sentuh ≥ 44 px.
- Input minim: default akun/kategori terakhir, kalkulator, template, autocomplete.
- Presisi "Japanese precision": angka tabular, format ¥ tanpa desimal, kategori dengan ikon konsisten.
- Status jelas: skeleton saat memuat, toast sukses/gagal, empty state yang mengajak aksi.
- Aksesibilitas: kontras AA, label form, navigasi keyboard di desktop.

### 9.3 Alur Kunci

1. **Catat cepat**: FAB → pilih tipe → nominal (kalkulator) → akun/kategori → simpan.
2. **Scan struk**: kamera/galeri → progres tahap (OCR, AI) → edit draf → konfirmasi.
3. **Jadwal**: pilih preset atau RRULE → template → pratinjau kemunculan → simpan.
4. **Bagikan akun**: pilih akun → masa berlaku → salin tautan → cabut kapan saja.

## 10. Kebutuhan Non-Fungsional

- **Keamanan**: Argon2id untuk password, JWT dengan secret ≥ 32 karakter dari env, cookie httpOnly/Secure/SameSite, CORS dibatasi origin frontend, rate limit login dan rute publik, header keamanan, tidak ada rahasia di repo atau log, token share dengan entropi ≥ 128 bit.
- **Privasi**: seluruh data berada di perangkat pengguna; layanan eksternal hanya menerima teks/gambar struk (OCR/LLM) dan tidak menerima data transaksi lain; ada opsi mematikan scan struk.
- **Kinerja**: indeks pada tanggal, akun, kategori, tipe; agregasi di SQL; bundel frontend per-route (lazy load); LCP < 2,5 s pada jaringan 4G.
- **Keandalan**: transaksi DB atomik untuk operasi multi-baris (batch, restore, scheduled), migrasi berurutan dan idempoten, backup konsisten saat WAL.
- **Kemudahan operasi**: berjalan di ARM/Termux tanpa CGO, satu binary, konfigurasi via env/YAML, log rotasi.
- **Kualitas**: cakupan tes unit backend ≥ 70% pada modul inti (transactions, accounts, analytics, scheduled, backup), tes komponen/store frontend, tes kontrak API terhadap Swagger, lint dan typecheck bersih di CI.
- **Internasionalisasi**: id, en, ja dengan kunci identik; format angka/tanggal mengikuti pengaturan.

## 11. Status Implementasi Saat Ini dan Kesenjangan

Bagian ini hasil membaca kode, migrasi, dan skema DB pada arsip. Hal yang ditandai **(terverifikasi)** dicek langsung ke skema/kode; sisanya **perlu verifikasi runtime**.

### 11.1 Sudah ada

- Backend modular lengkap per domain dengan rute CRUD, analitik, FX, backup, import/export, share, receipts, tunnel.
- Frontend dengan ±20 view, widget analitik, i18n tiga bahasa, PWA, komponen UI dasar.
- Skrip `start.sh` untuk PM2 (dev/prod), konfigurasi YAML + env, dokumentasi Swagger minimal.

### 11.2 Kesenjangan prioritas tinggi (blocker rilis)

1. **Migrasi 008 ganda (terverifikasi).** Ada dua berkas `008_*` (`backups_users_fts` dan `receipt_jobs`). Runner hanya menerapkan versi yang lebih besar dari versi terakhir, sehingga yang kedua dilewati; DB contoh berisi `schema_version` 1–9 dan **tidak punya tabel `receipt_jobs`/`receipt_job_events`**. Fitur scan struk gagal di instalasi baru. Perbaikan: ubah penomoran (mis. `010_receipt_jobs.sql`) dan sediakan migrasi yang aman untuk DB yang sudah berjalan.
2. **Skema vs query tidak cocok (terverifikasi pada goals dan shared).** Tabel `accounts`, `categories`, `transactions` tidak punya `user_id` dan transaksi memakai `date_time`/`from_account_id`/`to_account_id`, sedangkan query Goals memakai `c.user_id`, `t.user_id`, `t.date`, dan Shared memakai `accounts.user_id`, `transactions.account_id`. Query ini akan gagal. Karena v1.0 single-user, solusi sederhana: hapus filter `user_id` dan sesuaikan nama kolom; modul analitik lain perlu diaudit dengan pola yang sama.
3. **Rumus saldo tampak tidak lengkap.** Ringkasan akun menjumlahkan nominal hanya dari sisi `from_account_id`; income (`to_account_id`) dan transfer masuk tidak terlihat dihitung, sementara tampilan Shared menghitung semua non-income sebagai pengurang tanpa memisah transfer. Perlu satu fungsi saldo tunggal (BR2) yang dipakai semua modul.
4. **`/auth/me` dan `/auth/password` berada di luar middleware.** Router auth di-mount sebelum grup terproteksi sehingga `GetUserID` selalu kosong dan kedua rute mengembalikan 401; guard route frontend bergantung pada `me`. Perlu diverifikasi runtime dan diperbaiki dengan middleware opsional pada subrute tersebut.
5. **Transaksi terjadwal belum dieksekusi.** Hanya `next_run` yang dihitung saat membuat jadwal; tidak ditemukan job yang membuat transaksi saat jatuh tempo (lihat SCH-2).
6. **Keamanan autentikasi.** Password default `123456` dibuat otomatis saat login pertama; JWT secret fallback `changeme`; rate limiter dan TOTP sudah ditulis tetapi tidak terhubung; tidak ada endpoint refresh; frontend langsung logout pada 401; Swagger menyebut cookie `auth_token` sedangkan kode memakai `access_token`; `config.yaml` memuat secret dev.

### 11.3 Kesenjangan menengah

- Scheduler FX mengabaikan parameter cron (hard-coded 04:00); tidak ditemukan pemicu backup terjadwal walau `backup_schedule` ada.
- Backup memakai path relatif `data/...` (bukan `data_dir`), menyalin file DB langsung saat WAL aktif, dan kunci enkripsi hanya SHA-256 dari passphrase tanpa salt/KDF.
- OCR dikunci ke bahasa `eng`, padahal struk target berbahasa Jepang.
- Persetujuan Telegram memakai `exec curl`, polling sinkron hingga 60 detik di dalam handler HTTP, dan hanya cocok dengan string YES/NO persis; token bot sebaiknya divalidasi saat startup.
- Pintasan PWA mengarah ke `/transactions/new?...` yang tidak ada di router.
- Respons error berupa teks polos, belum seragam JSON.
- Kurs disimpan sebagai REAL dan fallback `COALESCE(rate, 1)` pada analitik bisa menghasilkan angka salah bila kurs hilang (melanggar BR4).
- Tidak ada satu pun berkas tes (`*_test.go`, `*.test.ts`) walau skrip `test` dan Vitest sudah dikonfigurasi.
- Dokumentasi: README merujuk `PRD.md`, `AGENTS.md`, dan `.tasks/` yang tidak ada di arsip; folder `docs/` kosong.

### 11.4 Higiene repositori (segera)

- Arsip memuat `.env`, `backend/.env`, `frontend/.env`, `cookies.txt`, `backend/data/kiroku.db`, log, dan `dist/`. Pastikan semuanya di `.gitignore` dan tidak pernah di-commit; **rotasi** semua kredensial yang pernah ada di berkas tersebut (JWT secret, API key OCR/9Router, token Telegram, token Upstash, passphrase backup).
- Bersihkan artefak agent (`.remember/`, `.code-hygiene/`) dari repo atau masukkan ke `.gitignore`.

## 12. Roadmap

### Fase 0 — Stabilisasi (blocker)

- Perbaiki migrasi ganda dan semua ketidakcocokan skema/query; satu fungsi saldo; verifikasi `me`/`password`.
- Hardening auth: hapus password default, wajib secret kuat, hubungkan rate limit, refresh token.
- Tulis tes kontrak dan reconciliation untuk transaksi, akun, saldo, laporan.
- Rotasi kredensial dan bersihkan repo.

### Fase 1 — Inti harian

- Eksekutor transaksi terjadwal, template, batch, tempat sampah, pintasan PWA yang valid.
- Scan struk end-to-end (OCR jpn+eng, validasi AI, konfirmasi, pembersihan data).
- Impor/ekspor CSV dengan deteksi duplikat.

### Fase 2 — Analitik dan kurs

- Audit seluruh endpoint analitik (konversi kurs benar, tanpa fallback diam-diam), budget/goals, outlook.
- Cron FX yang benar, kurs manual, peringatan kurs usang.

### Fase 3 — Keamanan data dan operasi

- Backup konsisten terjadwal, KDF, retensi otomatis, restore teruji.
- Akses bersama (rate limit, tampilan terbatas), TOTP.
- CI (lint, typecheck, test), Swagger sinkron, dokumentasi operasional.

### Fase 4 — Polesan

- Mode offline-read, aksesibilitas, performa pada dataset besar, deteksi duplikat struk.

## 13. Kriteria Penerimaan Rilis v1.0

- Instalasi baru dari nol (tanpa DB lama) menjalankan semua migrasi dan seluruh fitur, termasuk scan struk, tanpa error.
- Dataset uji 10.000 transaksi multi-currency: saldo akun = jumlah laporan = jumlah analitik (selisih 0).
- Jadwal bulanan menghasilkan tepat satu transaksi per bulan, termasuk setelah backend dimatikan lalu dinyalakan kembali.
- Login diblokir setelah batas percobaan; tidak ada kredensial default; `me` dan `password` bekerja dengan sesi valid.
- Backup terenkripsi dapat di-restore ke mesin bersih dan menghasilkan data identik.
- Tautan share kedaluwarsa/dicabut tidak lagi dapat diakses; tampilan publik tidak membocorkan data akun lain.
- `npm run lint`, `npm run test`, dan build prod lulus tanpa peringatan kritis.

## 14. Risiko, Asumsi, dan Pertanyaan Terbuka

### Risiko

- Ketergantungan layanan gratis (OCR.space, exchangerate.host) dapat berubah kebijakan/kuota; sediakan penyedia cadangan dan degradasi yang jelas.
- SQLite pada perangkat dengan penyimpanan lambat/terbatas (Termux) berisiko korup bila mati mendadak; WAL + backup konsisten mengurangi risiko.
- Mengirim gambar struk ke layanan eksternal menimbulkan risiko privasi; perlu persetujuan eksplisit dan opsi menonaktifkan.
- Tunnel Cloudflare gratis tidak menjamin URL stabil; penemuan URL via Upstash menambah titik kegagalan.

### Asumsi

- Satu pengguna, satu instance; tidak ada kebutuhan multi-tenant pada v1.0.
- Mata uang dasar JPY (dapat diubah di pengaturan); kurs harian cukup, bukan real-time.
- Perangkat host Linux/Termux dengan PM2 dan akses internet berkala.

### Pertanyaan Terbuka

- Apakah saldo awal akun perlu menjadi field tersendiri, atau cukup transaksi "penyesuaian saldo"?
- Apakah multi-user (pasangan/keluarga) masuk roadmap pasca-v1.0 sehingga skema perlu `user_id` sekarang?
- Penyedia OCR/LLM utama dan cadangan apa yang diinginkan jika OCR.space atau 9Router tidak tersedia?
- Apakah akses bersama perlu mendukung rentang tanggal terbatas atau hanya seluruh riwayat akun?
- Apakah frontend tetap di Vercel, atau backend sebaiknya juga menyajikan bundel statis agar tidak butuh Upstash?

## Lampiran A — Variabel Konfigurasi

- Backend: `PORT`/`addr`, `data_dir`, `db_path`, `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `ARGON2_*`, `FX_BASE_CURRENCY`, `backup_schedule`, `backup_retention_*`, `BACKUP_PASSPHRASE`.
- Scan struk: `OCR_SPACE_API_KEY`, `KIROKU_MOCK_OCR_ENABLED`, `NINE_ROUTER_URL`, `NINE_ROUTER_API_KEY`, `NINE_ROUTER_MODEL_COMBO_NAME`.
- Notifikasi/ops: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`, `KIROKU_DEV`.
- Frontend: `VITE_API_BASE_URL`, `VITE_BACKEND_URL`, `VITE_UPSTASH_URL`, `VITE_UPSTASH_TOKEN`.
- Catatan: nama variabel env di `.env.example` dan loader koanf (`KIROKU_` + nama bertingkat) belum seragam; seragamkan dan dokumentasikan satu konvensi.

## Lampiran B — Glosarium

- **Must / Have / Want** — sifat pengeluaran: wajib, perlu, keinginan.
- **RRULE** — standar iCalendar untuk aturan pengulangan.
- **FTS5** — fitur pencarian teks penuh SQLite.
- **Soft delete** — penandaan `deleted_at` tanpa menghapus baris.
- **Reconciled** — transaksi yang sudah dicocokkan dengan mutasi rekening nyata.
