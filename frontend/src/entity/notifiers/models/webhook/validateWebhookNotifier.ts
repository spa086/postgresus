import type { WebhookNotifier } from './WebhookNotifier';

export const validateWebhookNotifier = (notifier: WebhookNotifier): boolean => {
  if (!notifier.webhookUrl) {
    return false;
  }

  return true;
};
