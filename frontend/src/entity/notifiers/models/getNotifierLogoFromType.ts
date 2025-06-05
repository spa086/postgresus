import { NotifierType } from './NotifierType';

export const getNotifierLogoFromType = (type: NotifierType) => {
  switch (type) {
    case NotifierType.EMAIL:
      return '/icons/notifiers/email.svg';
    case NotifierType.TELEGRAM:
      return '/icons/notifiers/telegram.svg';
    default:
      return '';
  }
};
