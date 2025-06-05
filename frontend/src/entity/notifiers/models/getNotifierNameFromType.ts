import { NotifierType } from './NotifierType';

export const getNotifierNameFromType = (type: NotifierType) => {
  switch (type) {
    case NotifierType.EMAIL:
      return 'Email';
    case NotifierType.TELEGRAM:
      return 'Telegram';
    default:
      return '';
  }
};
