import { describe, test, expect } from 'vitest'
import { readFileSync, existsSync, readdirSync } from 'node:fs'
import { join, resolve } from 'node:path'

const ROOT = resolve(__dirname, '../../')
const SRC = resolve(ROOT, 'src')

function read(file: string): string {
  return readFileSync(join(ROOT, file), 'utf-8')
}

const routerSrc = read('src/router/index.ts')
const sidebarSrc = read('src/composables/useSidebarNav.ts')

const locales = ['en', 'ja', 'id'].map((lang) => ({
  lang,
  src: read(`src/locales/${lang}.ts`),
}))

const deletedFiles = [
  'src/widgets/shared/useShare.ts',
  'src/widgets/shared/CreateShareModal.vue',
  'src/widgets/shared/ShareCard.vue',
  'src/widgets/shared/ShareResultModal.vue',
  'src/widgets/shared/index.ts',
  'src/components/shared/SharedAccessList.vue',
  'src/components/shared/CreateShareModal.vue',
  'src/components/icons/SharedAccessIcon.vue',
  'src/views/shared/SharedAccessView.vue',
  'src/views/shared/SharedView.vue',
]

const deletedDirs = ['src/widgets/shared', 'src/components/shared', 'src/views/shared']

function listSourceFiles(dir: string): string[] {
  const result: string[] = []
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = join(dir, entry.name)
    if (entry.isDirectory()) {
      result.push(...listSourceFiles(full))
    } else if (entry.isFile() && /\.(ts|vue)$/.test(entry.name)) {
      result.push(full)
    }
  }
  return result
}

const srcFiles = listSourceFiles(SRC)

const sharedImportRegex = /from\s+['"]@\/(widgets|components|views)\/shared(['"]|\/)/

describe('shared access removal', () => {
  describe('router', () => {
    test('does not contain /shared-access route path', () => {
      expect(routerSrc).not.toMatch(/path:\s*['"]\/shared-access['"]/)
    })

    test('does not contain /shared/:token route path', () => {
      expect(routerSrc).not.toMatch(/path:\s*['"]\/shared\/:token['"]/)
    })

    test('does not import SharedAccessView', () => {
      expect(routerSrc).not.toContain('@/views/shared/SharedAccessView.vue')
    })

    test('does not import SharedView', () => {
      expect(routerSrc).not.toContain('@/views/shared/SharedView.vue')
    })
  })

  describe('sidebar nav', () => {
    test('does not contain shared-access path', () => {
      expect(sidebarSrc).not.toMatch(/\/shared-access/)
    })

    test('does not reference nav.sharedAccess label', () => {
      expect(sidebarSrc).not.toContain('nav.sharedAccess')
    })

    test('does not import SharedAccessIcon', () => {
      expect(sidebarSrc).not.toContain('@/components/icons/SharedAccessIcon.vue')
    })
  })

  describe('locale files', () => {
    test.each(locales)('lang $lang does not contain nav.sharedAccess key', ({ lang, src }) => {
      expect(src).not.toContain('sharedAccess:')
    })

    test.each(locales)('lang $lang does not contain shared block', ({ lang, src }) => {
      expect(src).not.toMatch(/shared\s*:\s*\{/)
    })
  })

  describe('deleted files', () => {
    test.each(deletedFiles)('file $path no longer exists on disk', (file) => {
      expect(existsSync(join(ROOT, file))).toBe(false)
    })

    test.each(deletedDirs)('directory $path no longer exists on disk', (dir) => {
      expect(existsSync(join(ROOT, dir))).toBe(false)
    })
  })

  describe('imports', () => {
    test('no src file imports from @/widgets/shared, @/components/shared, or @/views/shared', () => {
      for (const file of srcFiles) {
        const content = readFileSync(file, 'utf-8')
        expect(content).not.toMatch(sharedImportRegex)
      }
    })
  })
})
