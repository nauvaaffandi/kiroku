import { BackgroundTask, BackgroundTaskResult } from '@capgo/capacitor-background-task'
import { checkAndProcessNewMonth } from '@/repository'

// Define task at module scope so the OS can locate it when starting background work.
BackgroundTask.defineTask('kiroku_monthly_dirs', async () => {
  try {
    const result = await checkAndProcessNewMonth()
    console.log('[BackgroundTask] kiroku_monthly_dirs done', result)
    return BackgroundTaskResult.Success
  } catch (e) {
    console.error('[BackgroundTask] kiroku_monthly_dirs failed', e)
    return BackgroundTaskResult.Failed
  }
})

export async function registerMonthlyDirsBackgroundTask(): Promise<void> {
  await BackgroundTask.registerTaskAsync('kiroku_monthly_dirs', {
    minimumInterval: 60,
  })
  console.log('[BackgroundTask] registered kiroku_monthly_dirs')
}
