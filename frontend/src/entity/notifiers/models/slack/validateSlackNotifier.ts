import type { SlackNotifier } from './SlackNotifier';

export const validateSlackNotifier = (notifier: SlackNotifier): boolean => {
  if (!notifier.botToken) {
    return false;
  }

  if (!notifier.targetChatId) {
    return false;
  }

  return true;
};
