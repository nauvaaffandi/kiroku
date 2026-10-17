<template>
  <Transition name="modal">
    <div v-if="showMapModal" class="modal-overlay" @click.self="close">
      <div class="modal map-modal" @click.stop>
        <div class="modal-header">
          <h2>Pick Location</h2>
          <button @click="close" class="modal-close" aria-label="Close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <div class="map-body">
          <div id="leaflet-map" class="map-container"></div>
          <p class="map-hint">Click on the map to set location coordinates.</p>
        </div>
        <div class="modal-actions">
          <button @click="close" class="btn-primary">Confirm Location</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps<{
  modelValue: boolean
  placeLat: number | null
  placeLng: number | null
  placeName: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'update:placeLat': [lat: number | null]
  'update:placeLng': [lng: number | null]
}>()

const showMapModal = ref(false)
let mapInstance: L.Map | null = null
let markerInstance: L.Marker | null = null

watch(() => props.modelValue, async (val) => {
  showMapModal.value = val
  if (val) {
    await nextTick()
    initMap()
  }
}, { immediate: true })

const close = () => {
  showMapModal.value = false
  emit('update:modelValue', false)
}

const initMap = () => {
  const el = document.getElementById('leaflet-map')
  if (!el || mapInstance) return

  const defaultLat = props.placeLat || 35.6762
  const defaultLng = props.placeLng || 139.6503

  mapInstance = L.map(el).setView([defaultLat, defaultLng], 13)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap'
  }).addTo(mapInstance)

  if (props.placeLat && props.placeLng) {
    markerInstance = L.marker([props.placeLat, props.placeLng]).addTo(mapInstance)
  }

  mapInstance.on('click', (e: L.LeafletMouseEvent) => {
    const { lat, lng } = e.latlng
    emit('update:placeLat', lat)
    emit('update:placeLng', lng)
    if (markerInstance) {
      markerInstance.setLatLng([lat, lng])
    } else {
      markerInstance = L.marker([lat, lng]).addTo(mapInstance!)
    }
  })
}
</script>

<style scoped>
.map-modal { max-width: 600px; width: 90vw; }
.map-body { padding: 1rem; }
.map-container { width: 100%; height: 350px; border-radius: var(--radius-sm); z-index: 1; }
.map-hint { font-size: 0.75rem; color: var(--color-neutral-500); margin-top: 0.5rem; text-align: center; }
</style>