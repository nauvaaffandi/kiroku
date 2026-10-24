# Explorer Report: Frontend UI Refactoring for Mobile-Only Architecture

## 1. Executive Summary
The Kiroku frontend is a Vue 3 application built with Vite, TypeScript, Pinia, Vue Router, and Tailwind CSS. The user's goal is to refactor all UI components to resolve existing glitches and strictly optimize the interface for **mobile-only** usage, ensuring all logic remains untouched.

## 2. Architecture & Tech Stack Overview
- **Framework:** Vue 3 (Composition API with `<script setup lang="ts">`)
- **State Management:** Pinia stores (`frontend/src/stores/`)
- **Styling:** Tailwind CSS + custom CSS variables / scoped styles (`frontend/src/assets/main.css`)
- **Routing:** Vue Router (`frontend/src/router/index.ts`)
- **Testing:** Vitest (`frontend/package.json`)

## 3. Key UI & Layout Components Identified
- `frontend/src/components/layout/AppShell.vue`: Manages desktop sidebar vs mobile bottom navigation, handling responsiveness based on window width (`< 768px`).
- `frontend/src/components/layout/Sidebar.vue`: Currently implements a desktop-style sidebar overlay drawer that can be toggled on mobile, but adds complexity and glitches.
- `frontend/src/components/layout/BottomNav.vue`: Main mobile bottom navigation.
- `frontend/src/components/layout/Header.vue`: Top header bar.
- **Widgets & Views:** Various widgets under `frontend/src/widgets/` and views under `frontend/src/views/` that inherit desktop max-widths (`max-width: 1200px;`, large padding, desktop table headers, etc.).

## 4. Identified Glitches & Mobile-Only Alignment Issues
1. **Desktop Responding Residue:** `AppShell.vue` and layout styles include desktop media queries (`@media (min-width: 768px)`) that shift layouts (`margin-left: 280px`), causing layout jumps and responsive glitches since the app is exclusively mobile-only.
2. **Sidebars and Overlays:** The desktop sidebar and complex drawer toggles are unnecessary when the app exclusively serves mobile devices.
3. **Padding and Max-Width Constraints:** `.app-main` sets `max-width: 1200px` and large padding (`1.5rem`), which cramps mobile viewports and causes overflow/glitches.
4. **Table / Card Responsiveness:** Many budget/analytics tables overflow horizontally without proper mobile card layouts.

## 5. Recommended Refactoring Plan
1. **Enforce Mobile-Only Layout (`AppShell.vue` & Layouts):**
   - Remove desktop media queries and desktop sidebar logic.
   - Force single-column full-width mobile view with fixed bottom navigation and top header.
2. **Clean up Styling & Glitches:**
   - Standardize padding/margins for mobile screens (`1rem` max).
   - Fix horizontal overflow issues across widgets and tables.
3. **Preserve Logic:**
   - Keep all Pinia stores, API clients, routes, and component props/emits untouched.
