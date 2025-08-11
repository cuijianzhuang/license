import {NotificationInstance} from 'antd/es/notification/interface';

/**
 * 复制文本到剪贴板
 * @param text 要复制的文本
 * @param notification Ant Design的通知组件实例
 * @param successMessage 复制成功的消息
 * @param failMessage 复制失败的消息
 * @param duration 通知显示的持续时间（秒）
 * @returns Promise<boolean> 复制是否成功
 */
export const copyToClipboard = async (
  text: string,
  notification: NotificationInstance,
  successMessage: string = 'Copied to clipboard',
  failMessage: string = 'Failed to copy',
  duration: number = 3
): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text);
    
    notification.success({
      message: successMessage,
      duration,
    });
    
    return true;
  } catch (error) {
    console.error('Copy failed:', error);
    
    notification.error({
      message: failMessage,
      duration,
    });
    
    return false;
  }
};

/**
 * 管理组件中的复制状态
 * @param key 用于标识复制项的键
 * @param copyingState 当前复制状态对象
 * @param setCopyingState 设置复制状态的函数
 * @param resetDelay 复制状态重置延迟（毫秒）
 */
export const manageCopyState = (
  key: string,
  copyingState: Record<string, boolean>,
  setCopyingState: React.Dispatch<React.SetStateAction<Record<string, boolean>>>,
  resetDelay: number = 2000
): void => {
  // 先设置按钮状态为复制中
  setCopyingState({ ...copyingState, [key]: true });
  
  // 在指定时间后重置复制状态
  setTimeout(() => {
    setCopyingState({ ...copyingState, [key]: false });
  }, resetDelay);
};

/**
 * 复制并管理复制状态的组合函数
 * @param key 用于标识复制项的键
 * @param text 要复制的文本
 * @param copyingState 当前复制状态对象
 * @param setCopyingState 设置复制状态的函数
 * @param notification Ant Design的通知组件实例
 * @param successMessage 复制成功的消息
 * @param failMessage 复制失败的消息
 * @param duration 通知显示的持续时间（秒）
 * @param resetDelay 复制状态重置延迟（毫秒）
 */
export const copyAndManageState = async (
  key: string,
  text: string,
  copyingState: Record<string, boolean>,
  setCopyingState: React.Dispatch<React.SetStateAction<Record<string, boolean>>>,
  notification: NotificationInstance,
  successMessage: string = 'Copied to clipboard',
  failMessage: string = 'Failed to copy',
  duration: number = 3,
  resetDelay: number = 2000
): Promise<void> => {
  // 先设置按钮状态为复制中
  setCopyingState({ ...copyingState, [key]: true });
  
  try {
    await navigator.clipboard.writeText(text);
    
    notification.success({
      message: successMessage,
      duration,
    });
  } catch (error) {
    console.error('Copy failed:', error);
    
    notification.error({
      message: failMessage,
      duration,
    });
  }
  
  // 在指定时间后重置复制状态
  setTimeout(() => {
    setCopyingState({ ...copyingState, [key]: false });
  }, resetDelay);
}; 