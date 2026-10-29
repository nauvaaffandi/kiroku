# Plan: Python E2E Qualification Suite untuk Kiroku

## 0. Metadata

- **Tujuan:** membangun qualification suite Python yang menguji backend dan frontend Kiroku dari sudut pandang pengguna nyata, bukan hanya unit test.
- **Output utama:** struktur `tests/` yang rapi, test runner Python, fixture database/server terisolasi, coverage matrix seluruh endpoint live, use-case per endpoint, browser flow frontend, negative/security checks, report, dan CI-ready commands.
- **Tahap pertama:** membuat fondasi test suite, kontrak test, fixture, helper, serta smoke test.
- **Tahap kedua:** mengimplementasikan E2E semua endpoint live dan semua use-case penting per endpoint.
- **Bahasa:** Python 3.
- **Framework rekomendasi:** `pytest` sebagai runner, `httpx` untuk HTTP API, `pytest-playwright` + Playwright untuk browser/frontend, `jsonschema` opsional untuk schema assertions, `pytest-xdist` hanya setelah suite stabil.
- **Batas file:** satu file test maksimal sekitar 1.000 baris; target normal 150–500 baris per file. Jika mendekati 800–900 baris, pecah berdasarkan resource/use-case.
- **Sumber kebenaran route:** wiring aktif di `backend/cmd/server/main.go`, bukan seluruh package yang memiliki `Routes()`.
- **Catatan penting:** endpoint yang memiliki handler tetapi belum dimount harus masuk kategori `unmounted`, bukan dipalsukan sebagai endpoint live.

---

## 1. Hasil akhir yang dianggap selesai

Suite dianggap selesai hanya jika seluruh kondisi berikut terpenuhi:

1. `tests/` memiliki struktur modular dan setiap file test tetap di bawah batas 1.000 baris.
2. Test dapat dijalankan dengan satu command dari root repository.
3. Backend test berjalan terhadap server nyata dengan database sementara yang kosong dan terisolasi.
4. Frontend test berjalan terhadap frontend nyata melalui browser Playwright.
5. Semua route live dari `main.go` memiliki minimal:
   - happy path;
   - invalid input;
   - missing resource atau not-found;
   - method/path boundary yang relevan;
   - auth behavior sesuai kondisi aktual;
   - assertion status, content type, response shape, dan persistence.
6. Semua use-case frontend yang memanggil API live memiliki browser flow atau API qualification yang sesuai.
7. Endpoint frontend yang saat ini tidak tersedia ditandai `xfail`/blocked dengan alasan yang jelas, bukan dianggap lulus.
8. Test tidak memakai data production, `backend/data/`, `.env`, tunnel URL, credential nyata, atau token nyata.
9. Test tidak bergantung pada urutan eksekusi global kecuali fixture memang mendefinisikan dependency dengan eksplisit.
10. Test failure memberikan endpoint, payload aman, status aktual, dan diagnosis yang mudah dicari.
11. Backend verification tetap dijalankan dari `backend/`:
    - `go build ./...`
    - `go test ./...`
12. Frontend verification tetap dijalankan dari `frontend/`:
    - `npm run build`
    - `npm run test:run`
13. Python suite memiliki smoke command dan full command yang terdokumentasi.

---

## 2. Kondisi codebase yang sudah dikonfirmasi

### 2.1 Backend live

Entrypoint: `backend/cmd/server/main.go`.

Route live:

- `GET /healthz`
- `/api/docs/*`
- `/shared/*`
- `/api/v1/auth/*`
- `/api/v1/accounts/*`
- `/api/v1/categories/*`
- `/api/v1/labels/*`
- `/api/v1/transactions/*`
- `/api/v1/templates/*`
- `/api/v1/scheduled/*`
- `/api/v1/reports/*`
- `/api/v1/settings/*`
- `/api/v1/fx/*`

Auth middleware saat ini hanya jelas dipasang pada `/api/v1/labels`. Route accounts, categories, transactions, templates, scheduled, reports, settings, FX, dan `/shared` saat ini dimount tanpa `AuthMiddleware`. Test harus menangkap perilaku aktual ini dan memberi label sebagai security gap; test tidak boleh menganggap route tersebut aman hanya karena frontend menggunakannya setelah login.

### 2.2 Backend unmounted

Package berikut memiliki handler atau pemanggil frontend, tetapi tidak terbukti dimount oleh `main.go`:

- analytics;
- backup;
- receipts;
- import/export;
- beberapa package legacy domain.

Endpoint tersebut masuk `tests/contract/unmounted_endpoints.py` atau test inventory, bukan smoke pass. Jika nanti dimount, pindahkan coverage-nya ke suite live setelah wiring diverifikasi.

### 2.3 Frontend

Router: `frontend/src/router/index.ts`.

Public route:

- `/login`
- `/shared/:token`

Route lain membutuhkan auth guard.

API client: `frontend/src/api/client.ts`, cookie-based dengan `withCredentials: true`, base URL runtime, dan interceptor 401 yang memanggil logout.

Frontend menggunakan Vue 3, Pinia, Vue Router, Vite, dan Vitest. Folder `frontend/tests/` belum ada. Python suite berfungsi sebagai qualification/E2E lintas proses; Vitest tetap untuk component/store-level test jika nanti ditambahkan.

---

## 3. Struktur direktori yang direncanakan

Buat struktur berikut di root repository:

```text
tests/
├── README.md
├── pyproject.toml
├── pytest.ini
├── conftest.py
├── requirements-dev.txt
├── config/
│   ├── __init__.py
│   ├── environments.py
│   └── routes.py
├── support/
│   ├── __init__.py
│   ├── api_client.py
│   ├── browser.py
│   ├── assertions.py
│   ├── data_factory.py
│   ├── process.py
│   ├── database.py
│   ├── diagnostics.py
│   └── types.py
├── fixtures/
│   ├── __init__.py
│   ├── auth_fixtures.py
│   ├── account_fixtures.py
│   ├── category_fixtures.py
│   ├── transaction_fixtures.py
│   ├── template_fixtures.py
│   ├── scheduled_fixtures.py
│   └── shared_fixtures.py
├── contract/
│   ├── test_health.py
│   ├── test_options_and_cors.py
│   ├── test_auth_contract.py
│   ├── test_route_inventory.py
│   ├── test_unmounted_endpoints.py
│   └── schemas/
│       ├── auth.py
│       ├── accounts.py
│       ├── categories.py
│       ├── labels.py
│       ├── transactions.py
│       ├── templates.py
│       ├── scheduled.py
│       ├── reports.py
│       ├── settings.py
│       ├── fx.py
│       └── shared.py
├── e2e/
│   ├── auth/
│   │   ├── test_login_session.py
│   │   ├── test_logout_and_expiry.py
│   │   ├── test_pin_flows.py
│   │   └── test_telegram_status.py
│   ├── accounts/
│   │   ├── test_account_read.py
│   │   ├── test_account_create_update.py
│   │   └── test_account_delete_archive.py
│   ├── categories/
│   │   ├── test_category_read_tree.py
│   │   ├── test_category_create_update.py
│   │   └── test_category_reorder_delete.py
│   ├── labels/
│   │   ├── test_label_crud.py
│   │   └── test_label_reorder.py
│   ├── transactions/
│   │   ├── test_transaction_crud.py
│   │   ├── test_transaction_types.py
│   │   ├── test_transaction_filters.py
│   │   ├── test_transaction_search.py
│   │   ├── test_transaction_batch.py
│   │   └── test_transaction_summary.py
│   ├── templates/
│   │   ├── test_template_crud.py
│   │   └── test_template_apply.py
│   ├── scheduled/
│   │   ├── test_scheduled_crud.py
│   │   └── test_scheduled_toggle.py
│   ├── reports/
│   │   ├── test_spending_report.py
│   │   └── test_trend_report.py
│   ├── settings/
│   │   ├── test_settings_read_write.py
│   │   └── test_currency_settings.py
│   ├── fx/
│   │   ├── test_fx_rates.py
│   │   └── test_fx_convert.py
│   ├── shared/
│   │   ├── test_share_management.py
│   │   └── test_public_shared_view.py
│   └── frontend/
│       ├── test_login_and_route_guard.py
│       ├── test_dashboard_and_accounts.py
│       ├── test_categories_and_transactions.py
│       ├── test_templates_and_scheduled.py
│       ├── test_reports_and_analytics.py
│       ├── test_settings_and_preferences.py
│       ├── test_shared_access.py
│       └── test_responsive_accessibility.py
├── qualification/
│   ├── test_cross_resource_lifecycle.py
│   ├── test_data_isolation.py
│   ├── test_idempotency_and_replay.py
│   ├── test_error_consistency.py
│   └── test_frontend_backend_alignment.py
└── reports/
    └── .gitkeep
```

### Struktur rules

- `contract/`: route shape dan protocol, bukan workflow panjang.
- `e2e/`: satu domain/use-case group per file.
- `qualification/`: behavior lintas endpoint atau lintas UI/API.
- `support/`: helper reusable, tanpa assertion bisnis yang spesifik ke satu resource.
- `fixtures/`: object/data setup; jangan menyimpan state antar test secara implicit.
- `config/`: route inventory dan environment config, bukan secret.
- `reports/`: output runtime di-ignore Git; jangan commit response berisi data user.
- Tidak membuat file test raksasa. Pisahkan ketika file mendekati 800 baris.
- Nama test harus menyebut behavior, misalnya `test_create_expense_updates_account_balance`.

---

## 4. Tooling Python yang dipilih

### 4.1 `pytest`

Runner utama. Alasan: fixture kuat, parametrization, markers, plugin ecosystem, readable failure output.

Markers yang diperlukan:

- `smoke`: server hidup dan jalur minimum.
- `contract`: status/shape/protocol.
- `api`: HTTP E2E backend.
- `browser`: frontend browser flow.
- `security`: auth, isolation, CORS, unauthorized behavior.
- `slow`: browser atau workflow panjang.
- `external`: Telegram/OCR/FX provider yang membutuhkan dependency luar.
- `unmounted`: inventory only; default tidak dianggap live pass.

### 4.2 `httpx`

Client HTTP sync atau async. Mulai dari sync agar fixture mudah dipahami; gunakan `AsyncClient` hanya bila ada kebutuhan nyata. Client harus menangani cookie jar, base URL, timeout, JSON parsing, dan diagnostic-safe logging.

### 4.3 Playwright

Gunakan `pytest-playwright` untuk frontend browser. Browser test harus menjalankan frontend nyata dan backend nyata, bukan mock axios. Default browser Chromium headless; opsi `--headed` untuk debugging.

### 4.4 `jsonschema` atau assertion Python

Mulai dengan assertion Python typed/helper untuk schema ringan. Tambahkan `jsonschema` hanya bila schema berulang sudah banyak dan dependency memang diperlukan. Jangan membuat schema formal spekulatif sebelum response aktual dikonfirmasi.

### 4.5 Dependency policy

- Jangan menambah dependency production.
- Tambahkan hanya dependency dev test.
- Pin major/minor secara konsisten.
- Pastikan Python version minimum ditulis eksplisit.
- Siapkan fallback bila Playwright browser belum ter-install.

---

## 5. Environment dan lifecycle test

### 5.1 Environment terisolasi

Test harus mendukung:

- `KIROKU_TEST_BACKEND_URL`, default `http://127.0.0.1:8080`;
- `KIROKU_TEST_FRONTEND_URL`, default `http://127.0.0.1:5173`;
- `KIROKU_TEST_START_SERVER=true|false`;
- `KIROKU_TEST_START_FRONTEND=true|false`;
- `KIROKU_TEST_DB_PATH` hanya bila memang backend mendukung config path secara aman;
- `KIROKU_TEST_KEEP_PROCESSES=true|false` untuk debugging;
- `PYTEST_ADDOPTS` untuk marker/verbosity.

Jangan membaca `.env`. Jangan default ke tunnel. Jangan memakai `backend/data/`.

### 5.2 Backend process

Ada dua mode:

1. **External mode:** test menganggap backend sudah berjalan; cocok untuk CI yang menjalankan service terpisah.
2. **Managed mode:** fixture menyalakan backend dengan config test dan database temporary directory.

Managed mode harus:

- membuat temporary directory unik;
- menggunakan database baru untuk run;
- men-disable tunnel;
- menunggu `GET /healthz` hingga timeout;
- menyimpan stdout/stderr hanya di artifact failure yang sudah disaring;
- terminate process dan child process saat selesai;
- tidak mengirim SIGKILL sebelum graceful shutdown kecuali timeout.

Sebelum implementasi managed mode, verifikasi mekanisme config backend benar-benar menerima database path dan config test. Jangan mengarang nama environment variable yang belum dibaca oleh `config.Load()`.

### 5.3 Frontend process

Frontend harus dijalankan dengan base API yang menunjuk backend test. Verifikasi bagaimana `backendUrl` store memilih URL; jangan mengandalkan hard-coded host bila runtime config memiliki mekanisme sendiri.

Fixture browser harus:

- menunggu frontend siap;
- membuat context baru per test atau per flow sesuai kebutuhan;
- memakai viewport desktop default dan viewport mobile 360px untuk marker responsive;
- menangkap console error dan failed network requests;
- menyimpan screenshot/trace hanya saat failure atau ketika diminta.

### 5.4 Test data isolation

Prioritas:

1. database baru per session/run;
2. bila tidak memungkinkan, namespace data dengan prefix unik;
3. cleanup resource yang dibuat test;
4. jangan menghapus semua data dari database bersama.

Data factory harus menghasilkan ID/name unik tanpa secret. Time value harus fixed atau timezone-aware agar test tidak flaky.

### 5.5 First login behavior

AGENTS.md mencatat first login dapat auto-create user dengan PIN default `123456`. Test harus menjadikan ini sebagai bootstrap behavior yang diverifikasi di environment test saja. Jangan menulis PIN ke log, report, screenshot, atau response artifact. Ambil dari fixture constant internal test bila kontrak ini memang masih aktual; beri nama generic dan redaction.

---

## 6. Support layer yang harus dibuat

### 6.1 `support/api_client.py`

Responsibility:

- request method wrapper;
- cookie jar per user/session;
- optional Authorization bearer support untuk coverage auth;
- JSON decode aman;
- expected status assertion;
- request ID capture;
- redaction payload/header;
- helper `get_json`, `post_json`, `put_json`, `delete`;
- tidak mencetak cookie, PIN, bearer token, tunnel URL, atau database path.

### 6.2 `support/assertions.py`

Minimal helper:

- `assert_status(response, expected)`;
- `assert_json_object(response)`;
- `assert_json_array(response)`;
- `assert_has_keys(value, keys)`;
- `assert_error_response(response)` yang menerima variasi `http.Error` dan JSON aktual;
- `assert_created_resource_id(payload)`;
- `assert_iso_datetime(value)`;
- `assert_no_secret(value)` untuk diagnostic output.

### 6.3 `support/data_factory.py`

Factory untuk:

- account bank/cash/card/investment/other;
- income/expense/transfer transaction;
- parent/child category;
- label valid/invalid;
- template;
- scheduled item;
- settings/currency;
- share expiration.

Semua factory harus mengembalikan payload yang eksplisit. Jangan menyembunyikan field penting di helper magic.

### 6.4 `support/process.py`

Process manager untuk backend/frontend. Requirements:

- command array, bukan shell string;
- timeout startup;
- health/readiness polling;
- process tree cleanup;
- sanitized logs;
- failure diagnostics.

### 6.5 `support/browser.py`

Wrapper Playwright:

- fresh context;
- login helper;
- route/console/network capture;
- screenshot on failure;
- local storage/cookie cleanup;
- semantic locator helper.

Selector priority:

1. accessible role/name;
2. label;
3. stable `data-testid` jika tersedia;
4. semantic text;
5. CSS selector hanya fallback.

Jangan mengubah production UI hanya untuk menambah selector tanpa kebutuhan yang jelas.

### 6.6 `support/diagnostics.py`

Failure output harus memuat:

- test node id;
- method + path;
- sanitized request body;
- status;
- sanitized response body;
- request ID;
- backend/frontend process status;
- screenshot/trace path jika ada.

Default output tetap ringkas. Detail penuh hanya saat failure.

---

## 7. Route inventory authoritative

Buat `config/routes.py` sebagai inventory terstruktur. Setiap entry memiliki:

- method;
- path template;
- domain;
- live/unmounted;
- auth expected menurut wiring aktual;
- request factory;
- response kind;
- happy path test module;
- negative cases;
- frontend caller jika ada;
- external dependency flag.

Inventory live berikut wajib di-cover.

### 7.1 Health/docs

- `GET /healthz`: 200, JSON content type, status `ok`.
- `GET /api/docs/`: reachable.
- `GET /api/docs/swagger.json`: valid Swagger document; jangan menganggap seluruh `paths` lengkap sebelum dokumentasi diperbaiki.
- `OPTIONS` route sample: 204 dan CORS behavior aktual.

### 7.2 Auth

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/pin/reset-request`
- `POST /api/v1/auth/pin/reset-verify`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/pin`
- `POST /api/v1/auth/pin/change-request`
- `POST /api/v1/auth/pin/verify`
- `POST /api/v1/auth/telegram/register-token`
- `GET /api/v1/auth/telegram/status`

### 7.3 Accounts

- `GET /api/v1/accounts/`
- `POST /api/v1/accounts/`
- `GET /api/v1/accounts/summary`
- `GET /api/v1/accounts/{id}`
- `PUT /api/v1/accounts/{id}`
- `DELETE /api/v1/accounts/{id}`

### 7.4 Categories

- `GET /api/v1/categories/`
- `GET /api/v1/categories/tree`
- `POST /api/v1/categories/`
- `PUT /api/v1/categories/reorder`
- `GET /api/v1/categories/{id}`
- `PUT /api/v1/categories/{id}`
- `DELETE /api/v1/categories/{id}`

### 7.5 Labels

- `GET /api/v1/labels/`
- `POST /api/v1/labels/`
- `PUT /api/v1/labels/reorder`
- `GET /api/v1/labels/{id}`
- `PUT /api/v1/labels/{id}`
- `DELETE /api/v1/labels/{id}`

### 7.6 Transactions

- `GET /api/v1/transactions/`
- `POST /api/v1/transactions/`
- `POST /api/v1/transactions/batch`
- `GET /api/v1/transactions/search`
- `GET /api/v1/transactions/autocomplete`
- `GET /api/v1/transactions/summary`
- `GET /api/v1/transactions/{id}`
- `PUT /api/v1/transactions/{id}`
- `DELETE /api/v1/transactions/{id}`

### 7.7 Templates

- `GET /api/v1/templates/`
- `POST /api/v1/templates/`
- `GET /api/v1/templates/{id}`
- `PUT /api/v1/templates/{id}`
- `DELETE /api/v1/templates/{id}`
- `POST /api/v1/templates/{id}/apply`

Catat bahwa `apply` saat ini hanya mengambil dan mengembalikan template; jangan mengharapkan transaksi baru sebelum behavior backend berubah.

### 7.8 Scheduled

- `GET /api/v1/scheduled/`
- `POST /api/v1/scheduled/`
- `GET /api/v1/scheduled/{id}`
- `PUT /api/v1/scheduled/{id}`
- `DELETE /api/v1/scheduled/{id}`
- `POST /api/v1/scheduled/{id}/toggle`

### 7.9 Reports

- `GET /api/v1/reports/spending`
- `GET /api/v1/reports/trend`

Parameter combinations wajib dicoba sesuai handler: default, period/grouping, months, account filter, invalid values.

### 7.10 Settings

- `GET /api/v1/settings/`
- `PUT /api/v1/settings/`
- `GET /api/v1/settings/currencies`
- `POST /api/v1/settings/currencies`
- `GET /api/v1/settings/{key}`
- `PUT /api/v1/settings/{key}`

Uji route collision `/currencies` versus `/{key}`.

### 7.11 FX

- `GET /api/v1/fx/rates`
- `GET /api/v1/fx/convert`

Uji required query (`from`, `to`, `amount`), base/date, amount zero/negative/decimal, unavailable pair, dan fallback rate aktual tanpa menganggap fallback sebagai financial correctness.

### 7.12 Shared

- `POST /shared/`
- `GET /shared/`
- `DELETE /shared/{id}`
- `GET /shared/public/{token}`

Uji create/list/revoke/public access, expired token, invalid token, revoked token, serta ownership/isolation. Catat known gap: management route menggunakan `system` user ID dan tidak memakai auth middleware.

---

## 8. Use-case coverage per domain

### 8.1 Auth/session

Happy paths:

1. First login dengan test bootstrap PIN.
2. Session cookie tersimpan.
3. `GET /auth/me` mengembalikan user.
4. Refresh/reload frontend tetap authenticated.
5. Logout menghapus session.
6. Login ulang berhasil.

Negative/security:

1. PIN kosong.
2. PIN terlalu pendek/panjang.
3. PIN salah.
4. malformed JSON.
5. missing body.
6. expired/invalid cookie.
7. no-cookie request ke `/auth/me`.
8. bearer token path bila handler mendukungnya.
9. 401 interceptor frontend mengarahkan ke login.
10. reset PIN request/verify dengan token invalid atau expired.
11. change PIN request/verify lifecycle.
12. Telegram status/register token diberi marker external bila Telegram tidak configured.

### 8.2 Accounts

1. List awal.
2. Create setiap tipe account yang didukung.
3. Create dengan currency/name minimal.
4. Read by ID.
5. Update partial field.
6. Archive melalui `archived`.
7. Delete.
8. Get summary setelah transaction income/expense.
9. Invalid ID.
10. Missing required value.
11. Duplicate/conflicting data.
12. Verify account data tidak bocor antar test user bila backend auth sudah diperbaiki.

### 8.3 Categories

1. Seed categories tersedia pada database baru.
2. List.
3. Tree.
4. Create parent.
5. Create child.
6. Read/update.
7. Reorder multiple items.
8. Delete unused category.
9. Delete category yang dipakai transaction.
10. Invalid color/icon/name.
11. Invalid parent/self-parent/cycle bila handler/repository punya constraint.

### 8.4 Labels

1. List default sort.
2. Sort/order query.
3. Create valid label.
4. Name empty.
5. Name >100 runes.
6. Valid six-digit hex color.
7. Invalid color formats.
8. Read/update/delete.
9. Reorder.
10. No auth cookie harus 401.

### 8.5 Transactions

1. Create expense.
2. Create income.
3. Create transfer dengan from/to account.
4. Read by ID.
5. List default.
6. Update amount, note, category, labels, status.
7. Soft delete dan verify tidak muncul pada list normal.
8. Filter account/category/type/status/date range.
9. Full text/search query.
10. Autocomplete by payee/note.
11. Summary by category/account/type.
12. Batch create/update/delete sesuai contract aktual.
13. Invalid date.
14. Invalid type/status.
15. Missing account/category.
16. Zero/negative/decimal amount behavior.
17. Cross-currency data.
18. Attachments fields without external upload unless API supports it.
19. Verify account balances and report inputs after mutation.
20. Verify amount conversion behavior because DB column integer while domain uses float64.

### 8.6 Templates

1. Create template.
2. List/read.
3. Update.
4. Apply existing template.
5. Verify current apply response is template, not transaction.
6. Delete.
7. Invalid template ID.
8. Missing required fields.
9. Template references deleted account/category.

### 8.7 Scheduled transactions

1. Create.
2. List/read.
3. Update schedule.
4. Toggle enabled/disabled.
5. Delete.
6. Invalid recurrence/date.
7. Missing linked resource.
8. Verify toggle idempotency or actual state transition.

### 8.8 Reports

1. Spending default.
2. Spending period variants.
3. Spending group by variants.
4. Spending account filter.
5. Trend default months.
6. Trend custom months.
7. Trend account filter.
8. Empty dataset.
9. Invalid parameter.
10. Results change after controlled transaction fixture.
11. Multi-currency limitation/fallback recorded, not guessed.

### 8.9 Settings/currencies

1. Read defaults.
2. Read known key.
3. Read unknown key.
4. Update one setting.
5. Update multiple settings.
6. Put key individually.
7. Add currency.
8. Duplicate currency.
9. Empty/invalid currency code.
10. Persistence after new request/session.
11. Frontend settings UI reflects value.
12. Clear-data frontend call is classified separately if backend route is not live.

### 8.10 FX

1. Rates default base.
2. Rates custom base/date.
3. Convert valid pair.
4. Convert decimal amount.
5. Missing from/to/amount.
6. Invalid amount.
7. Unknown currency.
8. unavailable provider/cache behavior.
9. Never assert a hard-coded external rate unless fixture freezes provider/database.

### 8.11 Shared access

1. Create share for account.
2. List shares.
3. Public token read.
4. Revoke.
5. Revoked token denied.
6. Expired token denied.
7. Invalid token denied.
8. Management route auth/isolation test.
9. Verify shared view exposes only selected account data.
10. Verify no unrelated account/transaction data leaks.

---

## 9. Frontend browser qualification

### 9.1 Login and router

- unauthenticated `/dashboard` redirects to `/login`;
- authenticated user cannot access guest-only login route or is redirected according to router behavior;
- login success redirects to intended destination/dashboard;
- invalid login displays user-facing error;
- reload preserves session;
- logout returns to login;
- expired 401 causes store logout and navigation.

### 9.2 Account flow

- dashboard/account page loads;
- create account form;
- edit account;
- archive/delete behavior;
- balance displayed;
- empty state;
- long account name at 360px.

### 9.3 Category/transaction flow

- category list/tree renders;
- create/edit category;
- transaction form opens;
- create expense/income/transfer;
- category/account selectors work;
- search/autocomplete works;
- edit/delete transaction;
- filter/date controls;
- hide-amount setting masks amounts.

### 9.4 Template/scheduled flow

- create/edit/delete template;
- apply template matches actual backend behavior;
- create/edit/toggle/delete scheduled item;
- loading/error/empty states.

### 9.5 Reports/analytics

- report pages load controlled fixture data;
- spending/trend charts render without console errors;
- empty reports render safely;
- account filter changes request/result;
- analytics pages calling unmounted routes are marked blocked until backend wiring exists.

### 9.6 Settings

- settings load/save;
- currency add;
- PIN flow UI;
- Telegram flow external marker;
- theme light/dark through actual store;
- language en/ja/id;
- hide amounts.

### 9.7 Shared view

- create share from UI;
- copy/open share link without exposing real tunnel URL in artifacts;
- public shared page loads token;
- revoked/expired token shows error state.

### 9.8 Responsive/accessibility

Run at 360px and default desktop:

- no horizontal overflow;
- buttons/inputs usable;
- icon-only buttons have accessible names;
- focus-visible visible;
- forms have labels;
- loading/errors are perceivable;
- income/expense not distinguished by color alone;
- long Japanese and Indonesian strings do not break layout;
- light/dark theme remains readable;
- hidden amounts remain hidden in visual and accessible text where intended.

---

## 10. Contract and negative test rules

Every endpoint test group should cover, where applicable:

1. valid request;
2. empty JSON object;
3. missing body;
4. malformed JSON;
5. wrong JSON type;
6. unknown fields;
7. missing path ID;
8. nonexistent path ID;
9. unsupported method;
10. query missing/invalid/duplicate;
11. content-type absent or wrong;
12. response content type;
13. response shape;
14. persistence after a second request;
15. idempotency for repeat mutation;
16. safe error body without internal stack/secret.

Do not demand behavior not present in source. If endpoint currently returns 500 for not-found, capture it as actual behavior and classify as defect; do not mask it by accepting any 4xx/5xx unless test is explicitly a characterization test.

---

## 11. Auth and security qualification

Create dedicated `security` tests, not scattered only in domain files.

### Required checks

- public auth endpoints reachable without session;
- protected labels rejected without session;
- currently unprotected data routes characterized explicitly;
- shared management routes characterized explicitly;
- user A cannot access user B data once auth isolation is implemented;
- cookies have expected flags, without logging values;
- invalid session rejected;
- CORS arbitrary-origin reflection is captured as known failure until fixed;
- credentialed cross-origin behavior is not silently accepted as secure;
- no secret appears in error body, logs, screenshot, trace, or report;
- response does not expose database paths or internal stack traces;
- public shared token only exposes permitted account data;
- OPTIONS does not accidentally bypass mutation behavior.

Security tests must distinguish:

- `expected_pass`: intended secure behavior;
- `characterization`: current behavior captured for regression awareness;
- `known_failure`: security defect currently present;
- `blocked`: cannot test until feature/wiring exists.

---

## 12. Unmounted endpoint qualification

Create an explicit inventory test that records these as unavailable until wiring is confirmed:

- analytics endpoints such as calendar, cashflow projection, balance, outlook, budget, goals, stats;
- backup endpoints;
- receipts scan/jobs/confirm;
- import/export;
- frontend calls to `/api/v1/backup/*`, `/api/v1/receipts/*`, `/api/v1/analytics/calendar` when no live mount exists.

The test should:

1. read the static inventory;
2. mark `live=false`;
3. optionally probe only in a separate opt-in test;
4. fail if a supposedly live endpoint silently disappears;
5. never convert an expected 404 into a pass for the full E2E suite;
6. be updated only after `main.go` wiring and auth behavior are verified.

---

## 13. Test execution phases

### Phase 0 — Recon and contract freeze

Deliverables:

- route inventory;
- endpoint source map;
- live/unmounted classification;
- request/response examples derived from source;
- external dependency list;
- known security gap list;
- decision on server startup mode.

Verification:

- compare inventory against `main.go`;
- compare payload field names against domain structs and handlers;
- no `.env`, DB, or tunnel data read.

### Phase 1 — Foundation

Implement:

- Python manifests;
- pytest config;
- root `conftest.py`;
- environment config;
- API client;
- assertions;
- diagnostics;
- process lifecycle;
- first health smoke test;
- first browser readiness test.

Verification:

```bash
python -m pytest -m smoke -q
```

### Phase 2 — Auth and contract baseline

Implement:

- login/session fixture;
- logout fixture;
- health/docs;
- auth contract;
- labels auth baseline;
- route inventory test;
- CORS/OPTIONS characterization.

Verification:

```bash
python -m pytest -m "smoke or contract" -q
```

### Phase 3 — CRUD API domains

Implement in this order:

1. accounts;
2. categories;
3. labels;
4. transactions;
5. templates;
6. scheduled;
7. settings;
8. FX;
9. reports;
10. shared.

Reason: transactions depend on accounts/categories/labels; reports depend on transaction fixtures; shared depends on account data.

Verification after each domain:

```bash
python -m pytest tests/e2e/<domain> tests/contract -q
```

### Phase 4 — Cross-resource qualification

Implement:

- transaction changes account balances;
- transaction changes report output;
- category/label references;
- template/scheduled linked resources;
- soft-delete visibility;
- settings persistence;
- shared data boundary;
- replay/idempotency.

### Phase 5 — Frontend browser flows

Implement browser tests after API flows are stable. Each flow uses controlled backend fixtures, not guessed seed state.

Verification:

```bash
python -m pytest -m browser -q
```

### Phase 6 — Responsive/accessibility/security

Implement:

- 360px checks;
- light/dark;
- en/ja/id;
- hidden amounts;
- keyboard/focus checks;
- auth isolation;
- CORS characterization;
- no-secret diagnostics.

### Phase 7 — Hardening and CI

Implement:

- retries only for process readiness, never for assertion failures;
- deterministic clock/data;
- parallelization only for isolated DB contexts;
- artifact retention on failure;
- CI command;
- report summary;
- flaky test quarantine marker with owner/reason/expiry.

---

## 14. Commands yang harus disediakan

Target commands dari root:

```bash
python -m pytest -m smoke -q
python -m pytest -m contract -q
python -m pytest -m api -q
python -m pytest -m browser -q
python -m pytest -m security -q
python -m pytest -q
```

Jika managed process belum aman, sediakan external-server variant:

```bash
KIROKU_TEST_START_SERVER=false KIROKU_TEST_START_FRONTEND=false python -m pytest -q
```

Debug browser:

```bash
python -m pytest -m browser --headed --slowmo 100
```

Jangan menjalankan seluruh suite otomatis terhadap tunnel atau database development.

---

## 15. CI plan

CI job minimal:

1. Python install.
2. Install test dependencies.
3. Install Playwright Chromium.
4. Build backend.
5. Build frontend.
6. Start isolated backend with tunnel disabled.
7. Start frontend configured toward test backend.
8. Run smoke/contract/API tests.
9. Run browser tests.
10. Upload sanitized failure artifacts.
11. Run existing Go and frontend test commands.

CI harus memperbaiki stale workflow yang saat ini memakai Go 1.20 dan working directory yang tidak cocok dengan module `backend/`; perubahan workflow adalah task terpisah kecuali diminta dalam implementasi plan.

---

## 16. Definition of test quality

Test tidak dianggap berkualitas hanya karena status 200. Setiap test mutation minimal harus memastikan:

- response benar;
- entity dapat dibaca kembali;
- perubahan terlihat pada list/summary/report yang relevan;
- invalid input tidak merusak state;
- data test tidak bocor ke test lain;
- error tidak mengandung secret/internal detail;
- behavior aktual dicatat bila berbeda dari intended behavior.

Test browser tidak dianggap cukup hanya karena halaman terbuka. Setiap flow harus memverifikasi:

- network request yang diharapkan;
- UI state setelah response;
- error/loading/empty state minimal satu kali per domain;
- persistence setelah reload bila behavior memang persistent;
- accessible name dan keyboard path untuk kontrol utama.

---

## 17. Known gaps yang harus ditulis di report

- Data routes selain labels belum memakai auth middleware.
- Shared management memakai hard-coded `system` user ID dan tidak terlindungi auth.
- CORS merefleksikan Origin apa pun dengan credentials.
- Default first-login PIN `123456` masih aktif menurut project notes.
- Backup, receipts, analytics, import/export belum dimount.
- Duplicate migration number `008` dapat memengaruhi receipt jobs table.
- Frontend memiliki callers untuk endpoint yang belum tersedia.
- Frontend Vitest belum memiliki `frontend/tests/`.
- CI Go workflow stale.
- FX external data tidak deterministik tanpa fixture/provider control.
- Telegram/OCR/backup membutuhkan environment-specific integration dan harus diberi marker.

Known gap bukan alasan untuk menghapus test. Buat characterization/blocked test dengan status yang jelas.

---

## 18. Urutan implementasi file yang disarankan

1. `tests/README.md`
2. `tests/requirements-dev.txt`
3. `tests/pyproject.toml`
4. `tests/pytest.ini`
5. `tests/config/environments.py`
6. `tests/config/routes.py`
7. `tests/support/types.py`
8. `tests/support/api_client.py`
9. `tests/support/assertions.py`
10. `tests/support/diagnostics.py`
11. `tests/support/process.py`
12. `tests/conftest.py`
13. `tests/contract/test_health.py`
14. `tests/contract/test_auth_contract.py`
15. auth fixtures and auth E2E
16. account/category/label fixtures and tests
17. transaction fixtures and tests
18. remaining domain tests
19. qualification tests
20. browser wrapper and frontend tests
21. responsive/security hardening
22. final README and command verification

Setiap tahap harus menghasilkan suite yang bisa dijalankan. Jangan membuat seluruh folder kosong lalu menunda verifikasi sampai akhir.

---

## 19. Guardrails implementasi

- Jangan menambah endpoint agar test pass.
- Jangan mengubah production behavior hanya untuk menyesuaikan assertion.
- Jika kontrak belum jelas, baca handler/domain/repository lebih dulu.
- Jangan menggunakan sleep tetap; gunakan polling readiness atau wait condition.
- Jangan retry assertion.
- Jangan hard-code database ID jika factory dapat mengambil ID hasil create.
- Jangan membagikan cookie/token dalam error output.
- Jangan menyimpan response production ke repository.
- Jangan memakai test order sebagai implicit dependency.
- Jangan menulis lebih dari sekitar 1.000 baris dalam satu test file.
- Jangan menguji unmounted endpoint sebagai live pass.
- Jangan menandai known security gap sebagai secure hanya karena request berhasil.
- Jangan menambah komentar kode kecuali diperlukan; dokumentasi plan ini bukan komentar source code.

---

## 20. Final verification checklist

### Repository

- [ ] `tests/` hanya berisi test/support/config yang diperlukan.
- [ ] Tidak ada `.env`, database, token, PIN, tunnel URL, atau data runtime.
- [ ] Semua file test di bawah 1.000 baris.
- [ ] `git diff --check` bersih.

### Python

- [ ] Dependency install berhasil.
- [ ] `python -m pytest -m smoke -q` pass.
- [ ] Contract suite pass atau failure diklasifikasikan.
- [ ] API suite pass atau known defect tercatat.
- [ ] Browser suite pass atau blocked dependency tercatat.
- [ ] Failure artifact tersanitasi.

### Backend/frontend

- [ ] `cd backend && go build ./...` pass.
- [ ] `cd backend && go test ./...` pass.
- [ ] `cd frontend && npm run build` pass.
- [ ] `cd frontend && npm run test:run` pass atau existing absence dicatat.

### Coverage

- [ ] Semua live route inventory memiliki test module.
- [ ] Semua mutation memiliki read-back/persistence assertion.
- [ ] Semua frontend caller live memiliki browser/API qualification.
- [ ] Unmounted caller tercatat blocked.
- [ ] Auth/security gaps terlihat di report.
- [ ] 360px, light/dark, en/ja/id, hide-amounts diuji.

---

## 21. Deliverable implementation berikutnya

Implementasi setelah plan ini harus dimulai dari Phase 0–1, bukan langsung menulis ratusan test. Urutan minimum:

1. freeze route inventory;
2. buat Python tooling;
3. buktikan server lifecycle dengan health smoke;
4. buktikan login/session fixture;
5. implement CRUD domain satu per satu;
6. implement cross-resource qualification;
7. implement browser flows;
8. harden security/responsive checks;
9. jalankan seluruh verification checklist;
10. laporkan endpoint yang pass, fail karena bug, blocked karena unmounted, dan skipped karena external dependency.

Plan ini sengaja tidak menyamarkan gap production. E2E suite harus menjadi alat qualification dan regression detector, bukan sekadar generator status hijau.
