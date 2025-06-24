import type { TelegramNotifier } from './TelegramNotifier';

export const validateTelegramNotifier = (notifier: TelegramNotifier): boolean => {
  if (!notifier.botToken) {
    return false;
  }

  if (!notifier.targetChatId) {
    return false;
  }

  return true;
};
