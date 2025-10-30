import {
  NotificationService,
  NotificationOptions,
  NotificationCategory,
} from '@/bindings/github.com/wailsapp/wails/v3/pkg/services/notifications';

export interface ShowDesktopNotificationParams {
  id?: string;
  title: string;
  subtitle?: string;
  body?: string;
  categoryId?: string;
  data?: { [key: string]: any };
}

export interface DesktopNotificationResult {
  success: boolean;
  error?: string;
}

/**
 * 桌面通知服务
 */
export class DesktopNotificationService {
  /**
   * 请求通知权限
   * @returns 是否授权
   */
  async requestAuthorization(): Promise<boolean> {
    try {
      return await NotificationService.RequestNotificationAuthorization();
    } catch (error) {
      console.error('Failed to request notification authorization:', error);
      return false;
    }
  }

  /**
   * 检查通知权限
   * @returns 是否已授权
   */
  async checkAuthorization(): Promise<boolean> {
    try {
      return await NotificationService.CheckNotificationAuthorization();
    } catch (error) {
      console.error('Failed to check notification authorization:', error);
      return false;
    }
  }

  /**
   * 显示桌面通知
   * @param params 通知参数
   * @returns 通知结果
   */
  async showNotification(
    params: ShowDesktopNotificationParams,
  ): Promise<DesktopNotificationResult> {
    try {
      // Check authorization first
      const authorized = await this.checkAuthorization();
      if (!authorized) {
        const granted = await this.requestAuthorization();
        if (!granted) {
          return {
            success: false,
            error: 'Notification permission denied',
          };
        }
      }

      // Create notification options
      const options = new NotificationOptions({
        id: params.id || `notification-${Date.now()}`,
        title: params.title,
        subtitle: params.subtitle,
        body: params.body,
        categoryId: params.categoryId,
        data: params.data,
      });

      // Send notification
      await NotificationService.SendNotification(options);

      return {
        success: true,
      };
    } catch (error) {
      console.error('Failed to show notification:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * 显示带操作按钮的通知
   * @param params 通知参数
   * @returns 通知结果
   */
  async showNotificationWithActions(
    params: ShowDesktopNotificationParams,
  ): Promise<DesktopNotificationResult> {
    try {
      // Check authorization first
      const authorized = await this.checkAuthorization();
      if (!authorized) {
        const granted = await this.requestAuthorization();
        if (!granted) {
          return {
            success: false,
            error: 'Notification permission denied',
          };
        }
      }

      // Create notification options
      const options = new NotificationOptions({
        id: params.id || `notification-${Date.now()}`,
        title: params.title,
        subtitle: params.subtitle,
        body: params.body,
        categoryId: params.categoryId,
        data: params.data,
      });

      // Send notification with actions
      await NotificationService.SendNotificationWithActions(options);

      return {
        success: true,
      };
    } catch (error) {
      console.error('Failed to show notification with actions:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * 注册通知类别（用于定义操作按钮）
   * @param category 通知类别
   */
  async registerCategory(category: NotificationCategory): Promise<void> {
    try {
      await NotificationService.RegisterNotificationCategory(category);
    } catch (error) {
      console.error('Failed to register notification category:', error);
      throw error;
    }
  }

  /**
   * 移除特定通知
   * @param identifier 通知标识符
   */
  async removeNotification(identifier: string): Promise<void> {
    try {
      await NotificationService.RemoveNotification(identifier);
    } catch (error) {
      console.error('Failed to remove notification:', error);
      throw error;
    }
  }

  /**
   * 移除所有已送达的通知
   */
  async removeAllDeliveredNotifications(): Promise<void> {
    try {
      await NotificationService.RemoveAllDeliveredNotifications();
    } catch (error) {
      console.error('Failed to remove all delivered notifications:', error);
      throw error;
    }
  }

  /**
   * 移除所有待处理的通知
   */
  async removeAllPendingNotifications(): Promise<void> {
    try {
      await NotificationService.RemoveAllPendingNotifications();
    } catch (error) {
      console.error('Failed to remove all pending notifications:', error);
      throw error;
    }
  }
}

export const desktopNotificationService = new DesktopNotificationService();
