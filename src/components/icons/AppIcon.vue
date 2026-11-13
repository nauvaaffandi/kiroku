<template>
  <svg
    class="app-icon"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <defs>
      <linearGradient id="kiroku-gradient" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="var(--accent)" />
        <stop offset="100%" stop-color="var(--accent-hover)" />
      </linearGradient>
    </defs>
    
    <!-- Base ledger/book shape -->
    <path
      class="ledger-base"
      d="M6 4 L26 4 C28.2091 4 30 5.79086 30 8 L30 24 C30 26.2091 28.2091 28 26 28 L6 28 C3.79086 28 2 26.2091 2 24 L2 8 C2 5.79086 3.79086 4 6 4 Z"
      fill="url(#kiroku-gradient)"
      stroke="var(--accent)"
      stroke-width="0.5"
    />
    
    <!-- Inner pages/depth -->
    <path
      class="ledger-pages"
      d="M7 5 L25 5 C25.5523 5 26 5.44772 26 6 L26 23 C26 23.5523 25.5523 24 25 24 L7 24 C6.44772 24 6 23.5523 6 23 L6 6 C6 5.44772 6.44772 5 7 5 Z"
      fill="var(--bg-card)"
      opacity="0.9"
    />
    
    <!-- Stylized "K" / 記 mark - Japanese brush stroke style -->
    <g class="ki-mark" stroke="var(--accent)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none">
      <!-- Vertical stem (like 言 radical) -->
      <path d="M11 9 L11 23" />
      <!-- Top horizontal (like ノ stroke) -->
      <path d="M11 9 L16 9" />
      <!-- Diagonal brush stroke (like 記's right side) -->
      <path d="M16 9 L21 15" />
      <!-- Bottom curve (completion stroke) -->
      <path d="M11 23 Q16 25 21 23" />
    </g>
    
    <!-- Small accent dots - representing entries/records -->
    <circle class="accent-dot" cx="9" cy="12" r="1.5" fill="var(--accent)" opacity="0.6" />
    <circle class="accent-dot" cx="9" cy="17" r="1" fill="var(--accent)" opacity="0.4" />
    <circle class="accent-dot" cx="23" cy="10" r="1" fill="var(--accent)" opacity="0.3" />
    <circle class="accent-dot" cx="24" cy="19" r="1.5" fill="var(--accent)" opacity="0.5" />
    
    <!-- Binding rings on the left (ledger detail) -->
    <g class="binding-rings" stroke="var(--accent)" stroke-width="1" fill="none" opacity="0.5">
      <circle cx="5.5" cy="10" r="1.5" />
      <circle cx="5.5" cy="16" r="1.5" />
      <circle cx="5.5" cy="22" r="1.5" />
    </g>
  </svg>
</template>

<script setup lang="ts">
// AppIcon - Kiroku (記録) custom brand icon
// Represents: ledger/book + stylized "K" / 記 brush stroke + entry dots
// Uses CSS variables for theming (--accent, --accent-hover, --bg-card)
</script>

<style scoped>
@reference "tailwindcss";

.app-icon {
  @apply inline-block h-full w-full shrink-0;
  width: 1em;
  height: 1em;
}

/* Ensure gradient uses currentColor context */
.ledger-base {
  stroke: currentColor;
}

.ki-mark {
  stroke: currentColor;
}

.binding-rings {
  stroke: currentColor;
}

.accent-dot {
  fill: currentColor;
}

/* Size variants via class */
.app-icon--sm {
  width: 1rem;
  height: 1rem;
}

.app-icon--md {
  width: 1.5rem;
  height: 1.5rem;
}

.app-icon--lg {
  width: 2rem;
  height: 2rem;
}

.app-icon--xl {
  width: 3rem;
  height: 3rem;
}

/* Monochrome variant for toolbar/header use */
.app-icon--mono .ledger-base,
.app-icon--mono .ledger-pages,
.app-icon--mono .ki-mark,
.app-icon--mono .binding-rings,
.app-icon--mono .accent-dot {
  stroke: currentColor;
  fill: currentColor;
}

.app-icon--mono .ledger-base {
  fill: none;
  stroke-width: 1.5;
}

.app-icon--mono .ledger-pages {
  fill: none;
  stroke: currentColor;
  stroke-width: 0.5;
}
</style>