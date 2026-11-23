// Barrel export for src/repository layer

// Infrastructure (core)
export { filesystem } from './core/filesystem';
export { cacheRegistry, MemoryCache, CacheRegistryImpl } from './core/cache';
export { writeQueue, type QueuedOperation } from './queue/write-queue';
export { atomicWrite, type AtomicWriteResult } from './core/atomic-write';
export { cleanup } from './core/cleanup';
export { manifestManager, type Manifest } from './core/manifest';
export { journal } from './core/journal';

// Base CRUD
export { BaseRepository } from './crud/base';
import type { EntityKind, FileEntry, Filter, RepositoryConfig, WriteOperation } from './types';
export type { EntityKind, FileEntry, Filter, RepositoryConfig, WriteOperation };

// Entity repositories
export { accountRepository, type AccountRepository } from './entities/account.repository';
export { categoryRepository, type CategoryRepository } from './entities/category.repository';
export { labelRepository, type LabelRepository } from './entities/label.repository';
export { notificationRepository, type NotificationRepository } from './entities/notification.repository';
export { scheduledRepository, type ScheduledRepository } from './entities/scheduled.repository';
export { templateRepository, type TemplateRepository } from './entities/template.repository';
export { budgetRepository, type BudgetRepository } from './entities/budget.repository';
export { goalRepository, type GoalRepository } from './entities/goal.repository';
export { settingRepository, type SettingRepository } from './entities/setting.repository';
export { transactionRepository, type TransactionRepository } from './entities/transaction.repository';
