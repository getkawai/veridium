// Entry component
import { App } from 'antd';
import type { MessageInstance } from 'antd/es/message/interface';
import type { ModalStaticFunctions } from 'antd/es/modal/confirm';
import type { NotificationInstance } from 'antd/es/notification/interface';
import { memo } from 'react';

let message: MessageInstance | null = null;
let notification: NotificationInstance | null = null;
let modal: Omit<ModalStaticFunctions, 'warn'> | null = null;

/**
 * AntdStaticMethods component must be mounted at app root (once).
 * It captures AntD static instances and stores them in module-scope variables.
 */
export default memo(() => {
  const staticFunction = App.useApp();
  message = staticFunction.message;
  modal = staticFunction.modal;
  notification = staticFunction.notification;
  return null;
});

export { message, modal, notification };
/**
 * Accessors — callers should use these instead of importing variables directly.
 * They throw a helpful error if called before the component is mounted.
 */
export function getMessage(): MessageInstance {
  if (!message) {
    throw new Error(
      'Antd message instance is not initialized. Ensure <AntdStaticMethods /> is mounted under your App root before using getMessage().'
    );
  }
  return message;
}

export function getModal(): Omit<ModalStaticFunctions, 'warn'> {
  if (!modal) {
    throw new Error(
      'Antd modal instance is not initialized. Ensure <AntdStaticMethods /> is mounted under your App root before using getModal().'
    );
  }
  return modal;
}

export function getNotification(): NotificationInstance {
  if (!notification) {
    throw new Error(
      'Antd notification instance is not initialized. Ensure <AntdStaticMethods /> is mounted under your App root before using getNotification().'
    );
  }
  return notification;
}