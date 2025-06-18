import type { WebhookMethod } from './WebhookMethod';

export interface WebhookNotifier {
  webhookUrl: string;
  webhookMethod: WebhookMethod;
}
