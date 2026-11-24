import dayjs from 'dayjs'
import { filesystem } from './filesystem'

const TRANSACTION_BASE = 'transaction'
const MARKER_PATH = '.monthly-dirs-marker.json'

function log(action: string, detail?: unknown) {
  console.log(`[ScheduleMonthDirs] ${action}`, detail ?? '')
}

export interface MonthDirsResult {
  year: number
  month: number
  daysCreated: number
  daysSkipped: number
}

function monthKey(year: number, month: number): string {
  return `${year}-${String(month).padStart(2, '0')}`
}

async function readMarker(): Promise<string | null> {
  try {
    const raw = await filesystem.readFile(MARKER_PATH)
    const obj = JSON.parse(raw) as { lastProcessedYearMonth?: string }
    return obj.lastProcessedYearMonth ?? null
  } catch {
    log('readMarker: marker not found, first run')
    return null
  }
}

async function writeMarker(yearMonth: string): Promise<void> {
  await filesystem.writeFile(MARKER_PATH, JSON.stringify({ lastProcessedYearMonth: yearMonth }))
  log('writeMarker', { yearMonth })
}

export async function ensureMonthDays(year: number, month: number): Promise<MonthDirsResult> {
  const monthDir = `${TRANSACTION_BASE}/${year}/${String(month).padStart(2, '0')}`

  log('ensureMonthDays start', { year, month, monthDir })

  await filesystem.ensureDir(`${TRANSACTION_BASE}/${year}`)
  await filesystem.ensureDir(monthDir)

  const lastDay = dayjs(`${year}-${String(month).padStart(2, '0')}-01`).endOf('month').date()
  let daysCreated = 0
  let daysSkipped = 0

  for (let d = 1; d <= lastDay; d++) {
    const day = String(d).padStart(2, '0')
    const filePath = `${monthDir}/${day}.json`
    const exists = await filesystem.fileExists(filePath)
    if (exists) {
      daysSkipped++
      continue
    }
    await filesystem.writeFile(filePath, JSON.stringify([]))
    daysCreated++
    log('createDayFile', filePath)
  }

  log('ensureMonthDays done', { year, month, lastDay, daysCreated, daysSkipped })
  return { year, month, daysCreated, daysSkipped }
}

export async function checkAndProcessNewMonth(): Promise<MonthDirsResult | null> {
  const now = dayjs()
  const year = now.year()
  const month = now.month() + 1
  const yearMonth = monthKey(year, month)

  log('checkAndProcessNewMonth', { now: now.format(), yearMonth })

  const last = await readMarker()
  log('readMarker last', { last })

  if (last === yearMonth) {
    log('month already processed, skip')
    return null
  }

  const result = await ensureMonthDays(year, month)
  await writeMarker(yearMonth)
  log('marker updated', { yearMonth })
  return result
}
