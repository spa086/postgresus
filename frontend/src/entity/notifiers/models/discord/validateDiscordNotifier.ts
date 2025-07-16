import type { DiscordNotifier } from './DiscordNotifier';

export const validateDiscordNotifier = (notifier: DiscordNotifier): boolean => {
  if (!notifier.channelWebhookUrl) {
    return false;
  }

  return true;
};
