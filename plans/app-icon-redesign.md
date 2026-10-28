# Plan Implementasi Redesign App Icon Kiroku

## 1. Ringkasan

Redesign ikon aplikasi frontend Kiroku agar terasa lebih khas melalui bentuk geometris yang unik, tetapi tetap sederhana, tenang, dan relevan untuk aplikasi finance tracker.

Target visual utama:

- Bukan ikon grafik saham generik.
- Bukan ilustrasi kompleks.
- Bukan dekorasi berlebihan.
- Menggunakan bentuk yang mudah dikenali pada ukuran kecil.
- Tetap selaras dengan identitas Kiroku: quiet ledger, paper-and-ink, modern restrained.
- Tetap terbaca pada launcher mobile, favicon browser, shortcut PWA, dan maskable icon.

Plan ini hanya mencakup perencanaan. Tidak ada implementasi final, commit, atau penghapusan perubahan pengguna pada tahap plan.

---

## 2. Tujuan

### 2.1 Tujuan produk

Memberikan Kiroku identitas visual yang lebih kuat saat aplikasi dibuka dari home screen atau dikenali melalui tab browser.

Ikon harus menyampaikan tiga hal secara singkat:

1. Ini adalah aplikasi pencatatan atau ledger.
2. Ini berkaitan dengan uang atau kesehatan finansial.
3. Ini adalah produk yang rapi, tenang, dan tidak agresif seperti aplikasi trading.

### 2.2 Tujuan desain

Membuat ikon berbasis shape unik yang:

- Memiliki siluet berbeda dari ikon finance biasa.
- Tetap menggunakan jumlah elemen yang sedikit.
- Mengutamakan bentuk luar dan kontras, bukan detail kecil.
- Tetap cocok pada light mode dan dark mode aplikasi.
- Tidak bergantung pada teks, angka, simbol mata uang, atau bahasa tertentu.
- Tidak menggunakan emoji sebagai ikon UI.
- Tidak menciptakan visual yang mengesankan investasi berisiko atau spekulasi.

### 2.3 Tujuan teknis

Menghasilkan aset yang valid dan konsisten untuk seluruh penggunaan berikut:

- Favicon browser.
- PWA manifest icon.
- Maskable PWA icon.
- PWA shortcut icon untuk expense, income, dan transfer jika aset shortcut memang digunakan.
- Preview lokal pada ukuran kecil.
- Build frontend yang sukses.

---

## 3. Kondisi repository yang ditemukan

### 3.1 Konfigurasi PWA

File utama:

- `frontend/vite.config.ts`

Konfigurasi PWA saat ini:

- Menggunakan `vite-plugin-pwa`.
- `registerType` bernilai `autoUpdate`.
- Strategi service worker `generateSW`.
- Manifest memiliki nama `Kiroku`.
- Manifest merujuk banyak PNG pada `/icons/`:
  - `/icons/icon-72x72.png`
  - `/icons/icon-96x96.png`
  - `/icons/icon-128x128.png`
  - `/icons/icon-144x144.png`
  - `/icons/icon-152x152.png`
  - `/icons/icon-192x192.png`
  - `/icons/icon-384x384.png`
  - `/icons/icon-512x512.png`
- Manifest shortcut merujuk:
  - `/icons/shortcut-expense.png`
  - `/icons/shortcut-income.png`
  - `/icons/shortcut-transfer.png`

Semua path yang dirujuk harus diverifikasi keberadaannya setelah implementasi.

### 3.2 HTML entry point

File:

- `frontend/index.html`

Kondisi saat ini:

- Favicon merujuk `/favicon.svg`.
- Manifest merujuk `/manifest.webmanifest`.
- Metadata theme color masih menggunakan `#4a90d9`.

Plan implementasi harus menentukan apakah favicon diganti agar menunjuk aset brand baru. Jangan meninggalkan favicon lama jika tujuan redesign adalah mengganti identitas app icon secara menyeluruh.

### 3.3 Aset publik

Folder:

- `frontend/public/`
- `frontend/public/icons/`

Saat inspeksi, aset icon PNG yang dirujuk manifest belum tersedia pada folder tersebut. Terdapat file:

- `frontend/public/icon.svg`

File tersebut dianggap sebagai percobaan sementara, bukan otomatis sebagai desain final. Sebelum implementasi final, file ini harus direview dan diputuskan: dipertahankan sebagai aset final, diperbaiki, atau dihapus.

### 3.4 Dependency

`frontend/package.json` saat inspeksi memiliki `sharp` di `devDependencies`. Dependency tersebut tidak boleh dipertahankan hanya karena eksperimen konversi ikon. Keputusan final:

- Jika tidak dipakai oleh script build atau generator yang dipertahankan, hapus.
- Jika dipakai, pastikan instalasi dan eksekusinya kompatibel dengan environment proyek.
- Jangan menambah dependency native baru hanya untuk menghasilkan beberapa aset PNG jika pendekatan SVG cukup.

### 3.5 Perubahan tidak terkait

Git status saat inspeksi menunjukkan perubahan/untracked backend terkait label:

- `backend/internal/db/migrations/015_labels.sql`
- `backend/internal/db/migrations/016_transaction_labels.sql`
- `backend/internal/domain/label.go`
- `backend/internal/repository/sqlite/labels.go`

Perubahan tersebut di luar scope redesign ikon dan tidak boleh disentuh, dihapus, di-stage, atau di-commit.

---

## 4. Batasan scope

### 4.1 Termasuk

- Desain visual ikon Kiroku.
- Aset SVG master.
- Aset PNG PWA jika manifest tetap memerlukannya.
- Aset shortcut PWA.
- Favicon reference.
- Manifest icon reference.
- Theme color metadata jika diperlukan agar selaras dengan desain.
- Script generator hanya jika benar-benar dibutuhkan dan dapat dipelihara.
- Verifikasi file aset.
- Verifikasi build frontend.
- Verifikasi tampilan ukuran kecil secara lokal atau melalui browser jika tersedia.

### 4.2 Tidak termasuk

- Mengubah logo atau brand di seluruh halaman aplikasi.
- Mengubah layout dashboard.
- Mengubah warna global aplikasi.
- Mengubah typography.
- Mengubah navigasi.
- Mengubah komponen UI.
- Mengubah backend.
- Mengubah service worker caching.
- Menambah fitur baru.
- Menambah icon library.
- Mengganti framework atau bundler.
- Mengubah bahasa/i18n.
- Mengubah shortcut URL.
- Mengubah nama aplikasi.

---

## 5. Arah desain yang dipilih

### 5.1 Konsep utama: folded ledger mark

Gunakan bentuk utama berupa lembar ledger atau kartu pencatatan yang dilipat secara geometris. Lipatan tersebut membentuk aksen visual yang dapat dibaca sebagai:

- halaman catatan,
- bentuk dompet/kartu,
- aliran saldo,
- atau simbol pertumbuhan yang tidak literal.

Konsep ini lebih khas daripada grafik batang biasa, tetapi tetap relevan dengan finance tracker.

### 5.2 Struktur bentuk

Ikon final idealnya terdiri dari:

1. Satu siluet dasar berbentuk rounded geometric tile atau lembar ledger.
2. Satu potongan/lipatan diagonal yang menjadi ciri khas.
3. Satu jalur atau blok internal sederhana yang menyiratkan pencatatan/aliran.
4. Maksimal satu aksen warna sekunder.

Jumlah elemen harus dibatasi agar bentuk tetap terbaca pada ukuran 16–32 px.

### 5.3 Karakter bentuk

Bentuk harus:

- Asimetris ringan agar terasa unik.
- Tidak terlalu ornamental.
- Memakai sudut dan radius yang konsisten.
- Tidak memakai outline tipis yang hilang pada ukuran kecil.
- Tidak memiliki detail mikro.
- Tidak terlihat seperti logo perusahaan investasi.
- Tidak terlihat seperti ikon spreadsheet generik.

### 5.4 Palet

Gunakan token desain yang sudah ada di `DESIGN.md`.

Warna utama yang disarankan:

- Light accent: `#2b67ad`.
- Dark accent: `#7db7ee` jika membutuhkan versi gelap.
- Income/positive accent: `#0d7a50` untuk light mode.
- Dark income accent: `#34d399` untuk dark mode.
- Surface light: `#fcfbf9`.
- Canvas light: `#f6f5f2`.
- Surface dark: `#18181b`.
- Canvas dark: `#121212`.

Ikon launcher sebaiknya memiliki background stabil dan kontras tinggi. Jangan memakai CSS variable dalam file SVG/PNG karena aset akan digunakan di luar DOM aplikasi.

### 5.5 Larangan visual

Jangan menggunakan:

- simbol `$` sebagai elemen utama,
- koin bertumpuk,
- grafik saham dengan panah naik yang klise,
- dompet realistis,
- kalkulator,
- bank building,
- ilustrasi 3D,
- gradient berat,
- shadow kompleks,
- teks kecil,
- angka saldo,
- emoji,
- banyak warna status sekaligus.

---

## 6. Spesifikasi aset

### 6.1 Master asset

Buat satu master asset vector, disarankan:

- `frontend/public/icon.svg`

Persyaratan:

- ViewBox persegi, idealnya `0 0 512 512`.
- Tidak memiliki external reference.
- Tidak menggunakan font.
- Tidak memakai filter berat.
- Tidak memakai image embedding.
- Bentuk utama harus tetap terbaca jika SVG dirasterisasi.
- Elemen penting berada dalam safe area.
- Kontras cukup pada background terang dan gelap.

### 6.2 Favicon

Favicon harus menggunakan master asset yang sama atau turunan langsung yang konsisten.

Pilihan implementasi, evaluasi berurutan:

1. Gunakan `icon.svg` sebagai favicon jika browser support dan hasilnya valid.
2. Jika kebutuhan browser lama atau toolchain mengharuskan PNG, hasilkan `favicon.png` dari master asset.
3. Jangan membuat desain favicon terpisah yang berbeda identitas tanpa alasan teknis.

Update `frontend/index.html` hanya jika diperlukan agar reference benar-benar menunjuk aset baru.

### 6.3 PWA icons

Jika `vite.config.ts` tetap menggunakan daftar ukuran PNG, sediakan seluruh ukuran yang dirujuk:

- 72×72.
- 96×96.
- 128×128.
- 144×144.
- 152×152.
- 192×192.
- 384×384.
- 512×512.

Setiap file harus:

- benar-benar valid PNG,
- berukuran sesuai nama file,
- berasal dari master asset yang sama,
- tidak kosong,
- tidak hanya berupa solid color placeholder,
- memiliki background yang cocok untuk maskable usage jika `purpose` mencakup `maskable`.

### 6.4 Maskable icon

Perhatikan safe zone maskable icon. Elemen penting harus berada di area tengah yang tidak terpotong oleh bentuk launcher platform.

Rencana verifikasi:

- Pastikan siluet luar tidak menjadi satu-satunya pembawa makna.
- Pastikan simbol internal tidak terlalu dekat edge.
- Pastikan icon tetap dapat dibaca jika platform memotong sudut.

Jika satu asset tidak dapat aman untuk `any maskable`, pisahkan purpose atau siapkan asset maskable khusus. Jangan mencantumkan `maskable` tanpa verifikasi visual.

### 6.5 Shortcut icons

Shortcut icons tetap harus mengikuti sistem visual yang sama.

Pilihan minimal:

- Gunakan master icon yang sama untuk ketiganya jika shortcut hanya membutuhkan branding.
- Jika shortcut membutuhkan pembeda, gunakan variasi kecil yang konsisten:
  - expense: aksen expense,
  - income: aksen income,
  - transfer: aksen accent atau neutral.

Jangan menambah ilustrasi baru yang tidak diperlukan. Shortcut icon bukan scope untuk tiga logo baru.

---

## 7. Keputusan teknis aset rasterisasi

### 7.1 Prinsip

Jangan menambah tool native hanya untuk menghasilkan PNG jika project dapat memakai SVG secara valid dan browser/platform target mendukungnya.

### 7.2 Jika PNG tetap diperlukan

Karena manifest saat ini merujuk PNG, pilih salah satu pendekatan berikut:

#### Opsi A — Generator yang sudah didukung environment

Gunakan tool yang sudah tersedia di environment developer/CI. Contoh harus dikonfirmasi dari repository atau environment, bukan diasumsikan.

Generator wajib:

- deterministic,
- menerima master SVG,
- menghasilkan ukuran yang tercantum,
- tidak mengubah master asset,
- tidak menyimpan file temporary ke repository,
- tidak membutuhkan dependency runtime aplikasi.

#### Opsi B — SVG-only manifest

Jika `vite-plugin-pwa` dan target browser mendukung SVG icon dengan baik, sederhanakan manifest agar menggunakan master SVG.

Sebelum memilih opsi ini, verifikasi:

- behavior `vite-plugin-pwa` pada manifest,
- dukungan platform install PWA yang ditargetkan,
- validitas `sizes: any`,
- dukungan `purpose: maskable` untuk SVG.

Jangan memilih opsi ini hanya karena mengurangi file tanpa memastikan hasil installable.

#### Opsi C — Commit aset PNG hasil generate

Jika tool generator tidak portable, hasilkan PNG sekali, commit hasilnya, dan jangan commit generator ad-hoc yang tidak akan dipakai lagi.

Aset committed harus diverifikasi secara binary dan dimension check.

### 7.3 Dependency cleanup

Setelah keputusan aset dibuat:

- Hapus `sharp` dari `package.json` jika tidak diperlukan.
- Sinkronkan lockfile.
- Jangan meninggalkan dependency eksperimen.
- Jangan menambah `canvas`, `jimp`, atau library native lain tanpa kebutuhan yang terbukti.

---

## 8. Rencana perubahan file

### 8.1 `frontend/public/icon.svg`

Tindakan:

- Review SVG sementara.
- Ganti dengan desain final folded ledger mark.
- Pastikan tidak ada komentar jika tidak diperlukan oleh codebase convention.
- Pastikan warna memakai nilai yang berasal dari DESIGN.md.
- Pastikan viewBox persegi.
- Pastikan elemen berada di safe area.

### 8.2 `frontend/public/icons/*`

Tindakan:

- Buat aset PNG yang memang dirujuk manifest, atau ubah manifest ke strategi SVG yang tervalidasi.
- Jangan membuat file yang tidak dirujuk.
- Jangan membuat placeholder.
- Pastikan semua ukuran konsisten.

### 8.3 `frontend/index.html`

Tindakan potensial:

- Update `rel="icon"` dari `/favicon.svg` ke master asset yang benar.
- Selaraskan `theme-color` dengan palette final jika metadata lama tidak sesuai.
- Jangan mengubah title atau viewport.

### 8.4 `frontend/vite.config.ts`

Tindakan potensial:

- Update daftar `icons` agar semua path benar.
- Update shortcut icons bila aset shortcut diubah.
- Pertahankan `NetworkFirst` API caching apa adanya.
- Jangan mengubah service worker caching sebagai bagian redesign icon.
- Jangan mengubah route shortcut.

### 8.5 `frontend/package.json`

Tindakan potensial:

- Hapus dependency generator yang tidak dipakai.
- Tambahkan script generator hanya jika diperlukan dan benar-benar akan dipakai lagi.
- Jangan menambah script tanpa kebutuhan maintenance.

### 8.6 `frontend/package-lock.json`

Tindakan:

- Hanya berubah jika dependency memang ditambah atau dihapus.
- Jangan melakukan upgrade massal dependency.

### 8.7 `plans/*`

File plan ini adalah satu-satunya dokumentasi yang dibuat pada tahap sekarang. Implementasi tidak boleh menambah README atau dokumentasi lain kecuali diminta.

---

## 9. Urutan implementasi yang wajib

### Tahap 1 — Baseline

1. Cek `git status --short`.
2. Pastikan perubahan backend label tidak disentuh.
3. Baca `AGENTS.md` dan `DESIGN.md`.
4. Cek semua reference icon di:
   - `frontend/index.html`,
   - `frontend/vite.config.ts`,
   - file konfigurasi PWA lain jika ada.
5. Cek folder `frontend/public/`.
6. Catat file yang benar-benar ada dan file yang hanya dirujuk.

Acceptance:

- Scope icon terpetakan.
- Perubahan backend tetap tidak tersentuh.
- Tidak ada asumsi bahwa file PNG tersedia sebelum diverifikasi.

### Tahap 2 — Pilih arah asset

1. Tetapkan folded ledger mark sebagai arah desain.
2. Buat master SVG sederhana.
3. Review siluet pada 512 px, 192 px, 96 px, 48 px, dan 32 px.
4. Pastikan bentuk internal tidak hilang.
5. Pastikan ikon tidak menyerupai trading app.

Acceptance:

- Ikon dapat dikenali tanpa teks.
- Ikon memiliki elemen shape khas.
- Detail maksimum tetap rendah.

### Tahap 3 — Integrasi favicon

1. Hubungkan `frontend/index.html` ke master asset.
2. Pastikan file reference tidak 404.
3. Pastikan hanya ada satu source of truth untuk identitas icon.

Acceptance:

- Browser favicon merujuk file yang ada.
- Tidak ada reference lama yang masih aktif secara tidak sengaja.

### Tahap 4 — Integrasi PWA

1. Tentukan SVG-only atau PNG generated assets.
2. Jika PNG dipakai, hasilkan seluruh ukuran yang dirujuk.
3. Validasi dimensi dan MIME file.
4. Validasi shortcut icons.
5. Review `purpose` dan maskable safe area.

Acceptance:

- Semua icon path manifest ada.
- Semua file valid.
- Tidak ada placeholder.
- Installability tidak rusak.

### Tahap 5 — Dependency dan repository hygiene

1. Hapus dependency yang hanya dipakai eksperimen jika tidak diperlukan.
2. Hapus script generator ad-hoc jika tidak menjadi bagian workflow final.
3. Jangan menyentuh perubahan backend user.
4. Cek `git diff --stat`.
5. Cek file untracked yang relevan.

Acceptance:

- Diff hanya menyentuh icon/PWA/frontend dependency yang diperlukan.
- Tidak ada binary atau generated file yang tidak dirujuk.
- Tidak ada secret atau file sensitif.

### Tahap 6 — Verifikasi build

Dari `frontend/`, jalankan:

```bash
npm run build
```

Perintah ini wajib karena menjalankan `vue-tsc` dan Vite build.

Acceptance:

- Type check sukses.
- Vite build sukses.
- PWA plugin tidak gagal karena asset missing.
- Tidak ada warning fatal terkait manifest icon.

### Tahap 7 — Verifikasi visual

Jika browser harness tersedia:

1. Jalankan dev server.
2. Buka aplikasi.
3. Periksa favicon pada tab.
4. Periksa light mode.
5. Periksa dark mode.
6. Periksa lebar 360 px.
7. Periksa app icon melalui manifest jika browser tool mendukung.
8. Periksa maskable crop secara manual atau dengan preview tool.

Jika browser harness tidak tersedia:

- Validasi SVG secara statis.
- Validasi file PNG secara binary/dimension.
- Gunakan screenshot atau tool lokal yang benar-benar tersedia.
- Catat keterbatasan tersebut dalam laporan akhir.

Acceptance:

- Bentuk tetap terbaca pada ukuran kecil.
- Kontras memadai.
- Tidak ada clipping.
- Tidak ada perubahan UI di luar scope.

---

## 10. Verifikasi teknis detail

### 10.1 Reference audit

Cari semua reference berikut:

- `favicon.svg`
- `icon.svg`
- `/icons/`
- `shortcut-expense`
- `shortcut-income`
- `shortcut-transfer`
- `manifest.webmanifest`

Pastikan setiap reference memiliki target nyata atau dihasilkan oleh build plugin secara resmi.

### 10.2 SVG audit

Periksa:

- XML valid.
- `viewBox` tersedia.
- Tidak ada external URL.
- Tidak ada script.
- Tidak ada embedded secret.
- Tidak ada font dependency.
- Tidak ada filter yang menyebabkan rasterization gagal.
- Ukuran file tetap kecil.

### 10.3 PNG audit

Jika PNG digunakan, periksa:

- File dapat dibuka oleh image parser yang tersedia.
- Width dan height sesuai nama.
- Color mode konsisten.
- Alpha channel sesuai kebutuhan.
- File tidak berisi placeholder solid-color.

### 10.4 Manifest audit

Periksa:

- `name` tetap `Kiroku`.
- `short_name` tetap `Kiroku`.
- `scope` tetap `/`.
- `start_url` tetap `/`.
- `display` tetap `standalone`.
- Semua icon path valid.
- Shortcut path valid.
- `purpose` tidak mengklaim maskable tanpa safe-area check.

### 10.5 Build audit

Perintah utama:

```bash
npm run build
```

Perintah tambahan hanya jika diperlukan oleh perubahan:

```bash
npm run test:run
```

Tidak wajib menambahkan test frontend baru karena perubahan berupa asset/configuration, kecuali ditemukan logic generator yang perlu diuji.

---

## 11. Acceptance criteria final

Implementasi dianggap selesai hanya jika semua syarat berikut terpenuhi:

### Visual

- [ ] Ikon memiliki bentuk geometris unik.
- [ ] Ikon tetap simple.
- [ ] Makna ledger/finance terasa tanpa simbol uang literal.
- [ ] Ikon tidak terlihat seperti grafik trading generik.
- [ ] Ikon selaras dengan Quiet Ledger pada `DESIGN.md`.
- [ ] Ikon terbaca pada ukuran kecil.
- [ ] Ikon aman untuk crop maskable.
- [ ] Kontras baik pada konteks launcher terang dan gelap.

### Teknis

- [ ] Master SVG valid.
- [ ] Favicon reference valid.
- [ ] Semua manifest icon tersedia.
- [ ] Semua shortcut icon tersedia jika tetap dirujuk.
- [ ] Dimensi PNG sesuai jika PNG dipakai.
- [ ] Tidak ada broken asset path.
- [ ] Tidak ada dependency generator yang tidak dipakai.
- [ ] Tidak ada perubahan backend label ikut tersentuh.
- [ ] `npm run build` sukses.

### Scope

- [ ] Tidak ada redesign UI lain.
- [ ] Tidak ada perubahan API/backend.
- [ ] Tidak ada perubahan service worker caching.
- [ ] Tidak ada perubahan i18n.
- [ ] Tidak ada secret atau file sensitif masuk diff.

---

## 12. Risiko dan mitigasi

### Risiko 1 — Ikon bagus pada 512 px tetapi buruk pada 32 px

Mitigasi:

- Review pada banyak ukuran.
- Hilangkan detail kecil.
- Gunakan siluet kuat.

### Risiko 2 — PNG manifest missing

Mitigasi:

- Audit path sebelum build.
- Pilih SVG-only secara sadar atau generate semua ukuran.
- Jangan mengandalkan file yang belum ada.

### Risiko 3 — Maskable crop memotong elemen penting

Mitigasi:

- Letakkan shape utama dalam safe area.
- Lakukan preview crop.
- Pisahkan asset maskable jika diperlukan.

### Risiko 4 — Dependency native gagal di Termux/CI

Mitigasi:

- Hindari dependency native baru.
- Gunakan tool yang sudah tersedia.
- Commit hasil PNG hanya jika generator tidak portable.

### Risiko 5 — Ikon tidak selaras dengan DESIGN.md

Mitigasi:

- Gunakan palette DESIGN.md.
- Hindari gradient dan shadow yang tidak ada dalam design language.
- Pertahankan restrained visual.

### Risiko 6 — Perubahan backend user tertimpa

Mitigasi:

- Jangan memakai reset, checkout, clean, atau perintah destruktif.
- Stage hanya file icon/frontend jika nanti diminta commit.
- Review status dan diff sebelum selesai.

### Risiko 7 — Favicon baru tidak tampil karena cache

Mitigasi:

- Ubah asset filename atau query/reference bila diperlukan.
- Verifikasi dengan hard reload.
- Catat cache limitation jika tidak dapat diuji.

---

## 13. Hal yang sengaja tidak dilakukan

- Tidak membuat logo teks `Kiroku`.
- Tidak membuat banyak variasi logo.
- Tidak menambah library icon.
- Tidak mengubah seluruh sistem warna aplikasi.
- Tidak mengubah icon komponen UI.
- Tidak mengubah shortcut behavior.
- Tidak memperbaiki masalah backend yang terlihat di git status.
- Tidak memperbaiki masalah dependency lain yang tidak terkait.
- Tidak mengubah service worker caching.
- Tidak menganggap file SVG sementara sebagai final tanpa review.

---

## 14. Output implementasi yang diharapkan

Setelah plan dijalankan, hasil akhir idealnya berupa:

1. Satu master icon SVG yang valid.
2. Favicon yang menunjuk master icon final.
3. PWA icon assets valid untuk semua ukuran yang dirujuk.
4. Shortcut assets valid atau reference manifest yang disederhanakan secara tervalidasi.
5. Konfigurasi PWA tanpa broken path.
6. Tidak ada dependency eksperimen tersisa.
7. `npm run build` sukses.
8. Git diff kecil dan terbatas pada scope icon/frontend.

---

## 15. Urutan kerja singkat untuk implementer

```text
Audit current references
  -> freeze unrelated backend changes
  -> finalize folded ledger mark
  -> create master SVG
  -> choose validated SVG/PNG strategy
  -> generate or reference assets
  -> update favicon and manifest
  -> remove unused generator dependencies
  -> run asset audit
  -> run npm run build
  -> review visual sizes and theme contexts
  -> review git diff
```

Urutan tersebut tidak boleh dilompati pada bagian asset audit dan build verification.
