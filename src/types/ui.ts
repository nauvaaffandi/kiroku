// UI Types
export type ToastType = 'success' | 'error' | 'info' | 'warning'

export interface Toast {
  id: number
  type: ToastType
  message: string
}

export interface ToastMessage {
  id: string
  message: string
  type: ToastType
  duration?: number
}

export interface ModalState {
  isOpen: boolean
  title?: string
  data?: any
}

export interface ConfirmDialogState {
  isOpen: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  type?: 'danger' | 'warning' | 'info'
  onConfirm?: () => void | Promise<void>
  onCancel?: () => void
}