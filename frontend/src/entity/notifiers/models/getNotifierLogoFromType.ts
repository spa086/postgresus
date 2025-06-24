import { NotifierType } from './NotifierType';

export const getNotifierLogoFromType = (type: NotifierType) => {
  switch (type) {
    case NotifierType.EMAIL:
      return '/icons/notifiers/email.svg';
    case NotifierType.TELEGRAM:
      return '/icons/notifiers/telegram.svg';
    case NotifierType.WEBHOOK:
      return '/icons/notifiers/webhook.svg';
    case NotifierType.SLACK:
      return '/icons/notifiers/slack.svg';
    default:
      return '';
  }
};
