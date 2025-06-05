import type { EmailNotifier } from './EmailNotifier';
import type { NotifierType } from './NotifierType';
import type { TelegramNotifier } from './TelegramNotifier';

export interface Notifier {
  id: string;
  name: string;
  notifierType: NotifierType;
  lastSendError?: string;

  // specific notifier
  telegramNotifier?: TelegramNotifier;
  emailNotifier?: EmailNotifier;
}
