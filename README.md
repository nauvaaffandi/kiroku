# 記録 (Kiroku)

Offline-first personal finance tracker. Single user, many accounts and currencies. Mobile-first PWA / Capacitor Android app with local SQLite only — **no backend**.

[![Vue 3](https://img.shields.io/badge/Vue-3-%234FC08D?style=for-the-badge&logo=vue.js)](https://vuejs.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-v4-%2338BDF8?style=for-the-badge&logo=tailwindcss)](https://tailwindcss.com)
[![SQLite](https://img.shields.io/badge/SQLite-local-%23003B57?style=for-the-badge&logo=sqlite)](https://sqlite.org)
[![Capacitor Android](https://img.shields.io/badge/Capacitor-Android-%2325D366?style=for-the-badge&logo=capacitor)](https://capacitorjs.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

---

## Table of Contents

1. [What is Kiroku?](#what-is-kiroku)
2. [Stack](#stack)
3. [Architecture](#architecture)
4. [Settings](#settings)
5. [OCR & AI Provider Settings](#ocr--ai-provider-settings)
6. [API Keys & At-Rest Encryption](#api-keys--at-rest-encryption)
7. [Commands](#commands)
8. [Android / Capacitor Build](#android--capacitor-build)
9. [Environment Variables (`.env` rule)](#environment-variables-env-rule)
10. [Known Gaps & Limitations](#known-gaps--limitations)
11. [Contributing & License](#contributing--license)

---

## What is Kiroku?

**Kiroku (記録 — "Record" or "Log" in Japanese)** is a privacy-first, offline-capable personal finance tracker. Your financial records never leave your device: data lives in a local SQLite database, not on a server.

- Single user, many accounts and multiple currencies.
- Mobile-first PWA / Capacitor Android app. App shell is a centered column (max 480px).
- Offline-first by design: no network request is required to record a transaction or view accounts.
- No third-party analytics, no cloud telemetry, no mandatory SaaS dependencies.
- Every user-facing string is internationalized across English (`en`), Japanese (`ja`), and Indonesian (`id`).

---

## Stack

All commands run from the repo root (there is no `frontend/` or `backend/` directory).

| Layer | Technology |
|---|---|
| UI framework | Vue 3 (`<script setup lang="ts">`) |
| State | Pinia stores (13 stores) |
| Routing | Vue Router |
| Styling | Tailwind CSS v4, consumed only via `<style scoped>` blocks with `@reference "tailwindcss";` and `@apply` |
| Build / dev server | Vite |
| Type check | `vue-tsc` (part of `npm run build`) |
| Tests | vitest (`tests/**`) |
| Local storage | Capacitor + `@capacitor-community/sqlite`, SQLite database `kiroku` with `no-encryption` (`src/storage/repository.ts`) |
| Icons | `@lucide/vue` |
| Dates | `dayjs` |
| Charts | `echarts` / `vue-echarts` |
| Export | `jszip` |
| Maps | `leaflet` |
| HTTP client | `axios` (`baseURL: ''` — no remote backend) |
| Native shell | `@capacitor/android` (Capacitor Android app) |
| i18n | Custom Pinia-based store (`src/locales/`, `useI18n()`), locales `en.ts`, `ja.ts`, `id.ts` |

See `package.json` for exact versions.

---

## Architecture

There is **no backend**. Data reads and writes go through Pinia stores → `src/storage/repository.ts` (local SQLite via Capacitor). `src/api/client.ts` creates a shared Axios instance with `baseURL: ''`. Do not hard-code a host and do not expect `/api/v1` endpoints.

Data layer (`src/storage/repository.ts`):
- `LocalDatabase` opens the `kiroku` SQLite database (`no-encryption`), runs migrations from `src/storage/schema.ts`, and ensures an `entities` table.
- Entity kinds: `account`, `category`, `label`, `notification`, `transaction`, `template`, `scheduled`, `budget`, `goal`, `setting`.
- Each entity row stores `id`, `kind`, `payload` (JSON), `updated_at`, and `deleted` (soft-delete flag).
- Queries filter `deleted = 0` unless the feature is trash or restore.

State:
- Pinia stores in `src/stores/` own reactive state; components call stores, stores call `src/storage/repository.ts` (local storage) or the (offline-only) API client.
- Settings live as one `setting` row (`id: 'settings'`, kind `setting`) whose `payload` is JSON shaped by `Settings` in `src/types/settings.ts`.

Code conventions (from `AGENTS.md`):
- UI strings go through i18n (`en`, `ja`, `id`); every new key added to all three locale files.
- Tailwind only via `<style scoped>` `@reference "tailwindcss";` + `@apply`; semantic class names in templates.
- Files at most 300 lines; split reusable functions/files when exceeding ~350 lines.
- Every `get` / `save` / `update` / `delete` logs via `console.log()` (input payload, response data, errors). Never log secrets, PINs, or raw auth tokens.
- `.env` is never opened or committed; read `.env.example` only (see section 9).

---

## Settings

Managed by `useSettingsStore` (`src/stores/settings.ts`) with defaults in `defaultSettings`. Fields:

| Section | Field | Type | Default |
|---|---|---|---|
| Theme | `theme` | `'light' \| 'dark' \| 'system'` | `'system'` |
| | `hide_amounts` | boolean | `false` |
| Currency / Language | `language` | `'en' \| 'id' \| 'ja'` | `'ja'` |
| | `base_currency` | string | `'USD'` |
| Date / number | `date_format` | `'YYYY-MM-DD' \| 'DD/MM/YYYY' \| 'MM/DD/YYYY'` | `'YYYY-MM-DD'` |
| | `number_format` | `'ja-JP' \| 'en-US' \| 'id-ID'` | `'ja-JP'` |
| | `week_start` | `0–6` (0 = Sunday) | `1` |
| Security | `totp_enabled` | boolean | `false` |
| OCR & AI | `ocr_enabled` | boolean | `false` |
| | `ocr_api_key` | string (secure storage) | `''` |
| | `ai_provider_enabled` | boolean | `false` |
| | `ai_provider_url` | string | `''` |
| | `ai_model` | string | `''` |
| | `ai_api_key` | string (secure storage) | `''` |

Store API: `fetchSettings()`, `updateSettings(Partial<Settings>)`, `getSetting(key)`, `setSetting(key, value)`, `createBackup()`, `clearAllData()`. Computed getters: `theme`, `language`, `baseCurrency`, `dateFormat`, `numberFormat`, `weekStart`, `isTotpEnabled`, `isHideAmounts`, `isOcrEnabled`, `isAiProviderEnabled`.

Theme is applied via `data-theme` on `<html>` through `useSettingsStore` / `useUIStore`; never manipulate `data-theme` directly elsewhere. Amounts are routed through `src/utils/amountDisplay.ts` so the hide-amounts setting stays effective; dates use `dayjs` respecting `date_format` and `week_start`.

Settings UI sections (`src/components/settings/`): `SettingsThemeSection`, `SettingsCurrencyLanguageSection`, `SettingsDataSection`, `SettingsOCRAISection`. Widgets at `src/widgets/settings` (`SettingSection`, `SettingItem`).

---

## OCR & AI Provider Settings

Two setting groups rendered by `SettingsOCRAISection.vue` (`src/components/settings/SettingsOCRAISection.vue`), each with its own toggle:

**OCR & Receipt Scanning** (`settings.ocrTitle` / i18n keys `settings.ocr*`):
1. `ocr_enabled` — toggle ("Enable OCR receipt scanning")
2. `ocr_api_key` — password input ("OCR API key"), disabled when OCR off

**AI Provider** (`settings.aiTitle` / i18n keys `settings.ai*`):
3. `ai_provider_enabled` — toggle ("Enable AI provider")
4. `ai_provider_url` — URL input ("AI provider URL"), disabled when AI off
5. `ai_model` — text input ("AI model"), disabled when AI off
6. `ai_api_key` — password input ("AI API key"), disabled when AI off

Warning banner (`settings.aiLocalOnlyWarning`, key `aiLocalOnlyWarning`):

> Stored locally only. No backend configured. OCR/AI calls not wired.

All values are stored as JSON `payload` in the local SQLite `setting` row; API key fields (`ocr_api_key`, `ai_api_key`) are routed to secure storage (section 6), so SQLite stores `''` for them after migration.

Offline-only note: no network request is wired for OCR or AI calls. `ocr_api_key` and `ai_api_key` are stored locally only. Actual OCR/AI provider integration requires a backend/network layer that does not exist in this repository.

---

## API Keys & At-Rest Encryption

`SECURE_KEYS = ['ocr_api_key', 'ai_api_key']` (`src/composables/useSecureStorage.ts`).

- `ocr_api_key` and `ai_api_key` are stored via `capacitor-secure-storage-plugin` (Android EncryptedSharedPreferences / Keystore-backed), **not** plaintext in SQLite.
- `migratePlaintextKeys()` moves any existing plaintext API keys from the SQLite `setting` payload into secure storage, then redacts them in SQLite to `''`. Keys whose secure write failed are intentionally left in SQLite plaintext so they are not lost.
- On load, secure storage wins over the (possibly redacted/empty) SQLite value.
- Backups (`createBackup()`) include every entity kind including `setting`, but **exclude API keys** — after a restore, the user must re-enter `ocr_api_key` and `ai_api_key`.

**Production prerequisite:** `capacitor-secure-storage-plugin` is declared in `package.json` (`^0.10.0`) and IS installed in `node_modules` at `^0.10.0`.

- `vite.config.ts` aliases `capacitor-secure-storage-plugin` → `src/stubs/secure-storage-plugin.ts` only in non-production modes (dev / test / CI); production builds resolve the real package from `node_modules`.
- The stub (`src/stubs/secure-storage-plugin.ts`) is an in-memory `Map`. It is **not at-rest encrypted** and is **dev/test/CI only**.
- Production Android at-rest encryption requires `npx cap sync android` + `./gradlew assembleDebug`. Never ship a production build that resolves to the stub.

`capacitor.config.ts` declares `plugins: { SecureStoragePlugin: {} }`.

---

## Commands

Run all commands from the repo root:

```bash
# Dev server (Vite, host: true, port 5173)
npm run dev

# Production build: vue-tsc type-check + vite build
npm run build

# Preview production build locally
npm run preview

# Unit tests (vitest)
npm run test:run
```

`npm run build` is also the type-check (`vue-tsc && vite build`).

---

## Android / Capacitor Build

Native shell: Capacitor Android (`@capacitor/android` ^6.2.3).

`capacitor.config.ts`:
- `appId: 'app.kiroku.finance'`
- `appName: 'Kiroku'`
- `webDir: 'dist'`
- `plugins: { SecureStoragePlugin: {} }`
- `server.url: 'http://192.168.0.160:5173'`, `cleartext: true` — local Android testing only (cleartext dev server)

Build flow:

```bash
# Build web assets
npm run build

# Sync web assets into the Android native project
npx cap sync android

# Build the debug APK
./gradlew assembleDebug
```

CI: `.github/workflows/build-app.yml` (Node 22, Java 21, Android SDK):
1. `npm ci --ignore-scripts`
2. `npm run build`
3. Copy `dist/*` to `android/app/src/main/assets/public/`
4. `chmod +x gradlew && ./gradlew assembleDebug` (in `android/`)
5. Upload `kiroku.apk` artifact

---

## Environment Variables (`.env` rule)

**Rule:** Never open or commit `.env`. Read `.env.example` only. Never print, log, or commit secrets, tokens, PINs, `data/`, `*.db*`, or tunnel URLs.

Frontend variables (from `.env.example`):

| Variable | Purpose |
|---|---|
| `VITE_BACKEND_URL` | Backend API base URL fallback (reference only; no backend exists) |
| `VITE_VAPID_PUBLIC_KEY` | Web Push VAPID public key (base64url) |
| `VITE_UPSTASH_URL` | Upstash Redis REST URL (reference only; deleted backend) |
| `VITE_FRONTEND_UPSTASH_REDIS_REST_TOKEN` | Upstash Redis REST token (reference only) |
| `VITE_USE_ERUDA` | Enable Eruda dev tools in mobile browsers (`true/false`) |
| `OCR_SPACE_API_KEY` | OCR.space API key |
| `KIROKU_MOCK_OCR_ENABLED` | Enable mock OCR (`true/false`) |

Backend variables in `.env.example` (PORT, JWT_SECRET, ARGON2_*, CLOUDFLARE_TUNNEL_ENABLED, etc.) are **reference only** — the Go backend is deleted; they are not read by the frontend. No `config.yaml` exists.

Runtime: `axios` baseURL is `''`. There is no remote backend — do not hard-code a host, do not expect `/api/v1` endpoints.

---

## Known Gaps & Limitations

- **Secure storage: real package installed, stub is dev/test/CI only.** `capacitor-secure-storage-plugin` is installed at `^0.10.0`. The vite alias resolves to the in-memory stub only in non-production modes (dev / test / CI); production builds resolve the real package from `node_modules`. Stub is not at-rest encrypted — never ship to production. Production Android at-rest encryption requires `npx cap sync android` + `./gradlew assembleDebug` (native rebuild).
- **Backups exclude API keys.** `createBackup()` zips all entity kinds including `setting`, but API keys live in secure storage and are not included; re-enter them after restore.
- **Tests folder.** vitest is configured for `tests/**`; the folder exists with test suites under `components/`, `composables/`, `i18n/`, `settings/`, `shared-access/`, `stores/`. Run `npm run test:run` to execute.
- **No `start.sh`.** Previously referenced PM2 + cloudflared wrapper; it does not exist. Do not reference it.
- **No `go.yml` CI.** Only `.github/workflows/build-app.yml` (Android APK) is active.
- **`.env.example` is partly decorative.** It contains commented-out backend reference vars (PORT, JWT_SECRET, ARGON2_*, 9Router, OCR, Telegram, Cloudflare, Upstash, Backup, VAPID, logging) that are not read from the environment; no `config.yaml` exists.
- **Cleartext dev server.** `capacitor.config.ts` `server.url` uses cleartext HTTP on a private LAN address for local Android testing only.

---

## Contributing & License

Contributions, bug reports, and pull requests are welcome. All changes must adhere to `AGENTS.md` (i18n `en`/`ja`/`id`, Tailwind via `<style scoped>` `@reference "tailwindcss";` + `@apply` semantic classes, offline-first SQLite only, no backend, `.env` never opened or committed, `console.log()` on data operations, max 300 lines per file, small diffs).

Distributed under the **MIT License**. See `LICENSE` for more information.
