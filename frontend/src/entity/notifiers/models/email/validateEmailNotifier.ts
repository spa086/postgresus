import type { EmailNotifier } from './EmailNotifier';

export const validateEmailNotifier = (notifier: EmailNotifier): boolean => {
  if (!notifier.targetEmail) {
    return false;
  }

  if (!notifier.smtpHost) {
    return false;
  }

  if (!notifier.smtpPort) {
    return false;
  }

  if (!notifier.smtpUser) {
    return false;
  }

  if (!notifier.smtpPassword) {
    return false;
  }

  return true;
};
