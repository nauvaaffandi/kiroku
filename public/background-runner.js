// Background runner task for Kiroku monthly transaction directories.
// Runs headless (no DOM, no web view). Attempts to create
// data/transaction/{year}/{month}/{day}.json files for every day of
// the current month when a new month is detected.
/* global addEventListener, Capacitor, console */

const EVENT = 'kiroku_monthly_dirs'
const TRANSACTION_BASE = 'transaction'
const MARKER_PATH = '.monthly-dirs-marker.json'

function log(action, detail) {
  try {
    console.log(`[BackgroundRunner] ${action}`, detail ?? '')
  } catch (e) {
    // console may not be available
  }
}

function getFilesystem() {
  try {
    const cap = (typeof Capacitor !== 'undefined' && Capacitor) || (typeof window !== 'undefined' && window.Capacitor) || {}
    const plugins = cap.Plugins || {}
    return plugins.Filesystem || null
  } catch (e) {
    return null
  }
}

async function readMarker(Filesystem) {
  const Directory = Filesystem && Filesystem.Directory ? Filesystem.Directory : null
  const opts = Directory ? { path: MARKER_PATH, directory: Directory.Documents } : { path: MARKER_PATH }
  const raw = await Filesystem.readFile(opts)
  const obj = JSON.parse(raw.data || '{}')
  return obj.lastProcessedYearMonth || null
}

async function writeMarker(Filesystem, yearMonth) {
  const Directory = Filesystem && Filesystem.Directory ? Filesystem.Directory : null
  const opts = Directory
    ? { path: MARKER_PATH, data: JSON.stringify({ lastProcessedYearMonth: yearMonth }), directory: Directory.Documents }
    : { path: MARKER_PATH, data: JSON.stringify({ lastProcessedYearMonth: yearMonth }) }
  await Filesystem.writeFile(opts)
  log('writeMarker', { yearMonth })
}

async function ensureMonthDays(Filesystem, year, month) {
  const monthDir = `${TRANSACTION_BASE}/${year}/${String(month).padStart(2, '0')}`

  log('ensureMonthDays start', { year, month, monthDir })

  const Directory = Filesystem.Directory
  await Filesystem.mkdir({ path: `${TRANSACTION_BASE}/${year}`, directory: Directory.Documents, recursive: true })
  await Filesystem.mkdir({ path: monthDir, directory: Directory.Documents, recursive: true })

  // last day of month
  const d = new Date(year, month, 0)
  const lastDay = d.getDate()

  let daysCreated = 0
  let daysSkipped = 0

  for (let day = 1; day <= lastDay; day++) {
    const dayStr = String(day).padStart(2, '0')
    const filePath = `${monthDir}/${dayStr}.json`
    let exists = false
    try {
      await Filesystem.stat({ path: filePath, directory: Directory.Documents })
      exists = true
    } catch (e) {
      exists = false
    }
    if (exists) {
      daysSkipped++
      continue
    }
    await Filesystem.writeFile({
      path: filePath,
      data: JSON.stringify([]),
      directory: Directory.Documents,
    })
    daysCreated++
    log('createDayFile', filePath)
  }

  log('ensureMonthDays done', { year, month, lastDay, daysCreated, daysSkipped })
  return { year, month, daysCreated, daysSkipped }
}

async function handleTask() {
  log('handleTask: start')
  const Filesystem = getFilesystem()
  if (!Filesystem) {
    log('handleTask: Filesystem plugin not available in headless runner; skipping dir creation')
    resolveTask()
    return
  }

  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() + 1
  const yearMonth = `${year}-${String(month).padStart(2, '0')}`

  log('handleTask: current', { yearMonth })

  let last
  try {
    last = await readMarker(Filesystem)
  } catch (e) {
    log('handleTask: readMarker failed, treating as first run', { error: String(e) })
    last = null
  }
  log('handleTask: last marker', { last })

  if (last === yearMonth) {
    log('handleTask: month already processed, skip')
    resolveTask()
    return
  }

  try {
    await ensureMonthDays(Filesystem, year, month)
    await writeMarker(Filesystem, yearMonth)
    log('handleTask: done', { yearMonth })
  } catch (e) {
    log('handleTask: error creating month dirs', { error: String(e) })
  }

  resolveTask()
}

let resolveTask = () => {}

addEventListener(EVENT, (_resolve, _reject, _args) => {
  resolveTask = _resolve
  handleTask().catch((err) => {
    log('handleTask: unhandled error', { error: String(err) })
    try {
      _reject(err)
    } catch (e) {
      // ignore
    }
  })
})
