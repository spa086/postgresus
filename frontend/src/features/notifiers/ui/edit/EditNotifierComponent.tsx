import { Button, Input, Select } from 'antd';
import { useEffect, useState } from 'react';

import {
  type Notifier,
  NotifierType,
  WebhookMethod,
  notifierApi,
} from '../../../../entity/notifiers';
import { getNotifierLogoFromType } from '../../../../entity/notifiers/models/getNotifierLogoFromType';
import { ToastHelper } from '../../../../shared/toast';
import { EditEmailNotifierComponent } from './notifiers/EditEmailNotifierComponent';
import { EditTelegramNotifierComponent } from './notifiers/EditTelegramNotifierComponent';
import { EditWebhookNotifierComponent } from './notifiers/EditWebhookNotifierComponent';

interface Props {
  isShowClose: boolean;
  onClose: () => void;

  isShowName: boolean;

  editingNotifier?: Notifier;
  onChanged: (notifier: Notifier) => void;
}

export function EditNotifierComponent({
  isShowClose,
  onClose,
  isShowName,
  editingNotifier,
  onChanged,
}: Props) {
  const [notifier, setNotifier] = useState<Notifier | undefined>();
  const [isUnsaved, setIsUnsaved] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const [isSendingTestNotification, setIsSendingTestNotification] = useState(false);
  const [isTestNotificationSuccess, setIsTestNotificationSuccess] = useState(false);

  const save = async () => {
    if (!notifier) return;

    setIsSaving(true);

    try {
      await notifierApi.saveNotifier(notifier);
      onChanged(notifier);
      setIsUnsaved(false);
    } catch (e) {
      alert((e as Error).message);
    }

    setIsSaving(false);
  };

  const sendTestNotification = async () => {
    if (!notifier) return;

    setIsSendingTestNotification(true);

    try {
      await notifierApi.sendTestNotificationDirect(notifier);
      setIsTestNotificationSuccess(true);
      ToastHelper.showToast({
        title: 'Test notification sent!',
        description: 'Test notification sent successfully',
      });
    } catch (e) {
      alert((e as Error).message);
    }

    setIsSendingTestNotification(false);
  };

  const setNotifierType = (type: NotifierType) => {
    if (!notifier) return;

    notifier.emailNotifier = undefined;
    notifier.telegramNotifier = undefined;

    if (type === NotifierType.TELEGRAM) {
      notifier.telegramNotifier = {
        botToken: '',
        targetChatId: '',
      };
    }

    if (type === NotifierType.EMAIL) {
      notifier.emailNotifier = {
        targetEmail: '',
        smtpHost: '',
        smtpPort: 0,
        smtpUser: '',
        smtpPassword: '',
      };
    }

    if (type === NotifierType.WEBHOOK) {
      notifier.webhookNotifier = {
        webhookUrl: '',
        webhookMethod: WebhookMethod.POST,
      };
    }

    setNotifier(
      JSON.parse(
        JSON.stringify({
          ...notifier,
          notifierType: type,
        }),
      ),
    );
  };

  useEffect(() => {
    setIsUnsaved(false);
    setNotifier(
      editingNotifier
        ? JSON.parse(JSON.stringify(editingNotifier))
        : {
            id: undefined as unknown as string,
            name: '',
            notifierType: NotifierType.TELEGRAM,
            telegramNotifier: {
              botToken: '',
              targetChatId: '',
            },
          },
    );
  }, [editingNotifier]);

  const isAllDataFilled = () => {
    if (!notifier) return false;

    if (!notifier.name) return false;

    if (notifier.notifierType === NotifierType.TELEGRAM) {
      return notifier.telegramNotifier?.botToken && notifier.telegramNotifier?.targetChatId;
    }

    if (notifier.notifierType === NotifierType.EMAIL) {
      return (
        notifier.emailNotifier?.targetEmail &&
        notifier.emailNotifier?.smtpHost &&
        notifier.emailNotifier?.smtpPort &&
        notifier.emailNotifier?.smtpUser &&
        notifier.emailNotifier?.smtpPassword
      );
    }

    if (notifier.notifierType === NotifierType.WEBHOOK) {
      return notifier.webhookNotifier?.webhookUrl;
    }

    return false;
  };

  if (!notifier) return <div />;

  return (
    <div>
      {isShowName && (
        <div className="mb-1 flex items-center">
          <div className="min-w-[110px]">Name</div>

          <Input
            value={notifier?.name || ''}
            onChange={(e) => {
              setNotifier({ ...notifier, name: e.target.value });
              setIsUnsaved(true);
            }}
            size="small"
            className="w-full max-w-[250px]"
            placeholder="Chat with me"
          />
        </div>
      )}

      <div className="mb-1 flex items-center">
        <div className="min-w-[110px]">Type</div>

        <Select
          value={notifier?.notifierType}
          options={[
            { label: 'Telegram', value: NotifierType.TELEGRAM },
            { label: 'Email', value: NotifierType.EMAIL },
            { label: 'Webhook', value: NotifierType.WEBHOOK },
          ]}
          onChange={(value) => {
            setNotifierType(value);
            setIsUnsaved(true);
          }}
          size="small"
          className="w-full max-w-[250px]"
        />

        <img src={getNotifierLogoFromType(notifier?.notifierType)} className="ml-2 h-4 w-4" />
      </div>

      <div className="mt-5" />

      <div>
        {notifier?.notifierType === NotifierType.TELEGRAM && (
          <EditTelegramNotifierComponent
            notifier={notifier}
            setNotifier={setNotifier}
            setIsUnsaved={setIsUnsaved}
          />
        )}

        {notifier?.notifierType === NotifierType.EMAIL && (
          <EditEmailNotifierComponent
            notifier={notifier}
            setNotifier={setNotifier}
            setIsUnsaved={setIsUnsaved}
          />
        )}

        {notifier?.notifierType === NotifierType.WEBHOOK && (
          <EditWebhookNotifierComponent
            notifier={notifier}
            setNotifier={setNotifier}
            setIsUnsaved={setIsUnsaved}
          />
        )}
      </div>

      <div className="mt-3 flex">
        {isUnsaved && !isTestNotificationSuccess ? (
          <Button
            className="mr-1"
            disabled={isSendingTestNotification || !isAllDataFilled()}
            loading={isSendingTestNotification}
            type="primary"
            onClick={sendTestNotification}
          >
            Send test notification
          </Button>
        ) : (
          <div />
        )}

        {isUnsaved && isTestNotificationSuccess ? (
          <Button
            className="mr-1"
            disabled={isSaving || !isAllDataFilled()}
            loading={isSaving}
            type="primary"
            onClick={save}
          >
            Save
          </Button>
        ) : (
          <div />
        )}

        {isShowClose ? (
          <Button
            className="mr-1"
            disabled={isSaving}
            type="primary"
            danger
            ghost
            onClick={onClose}
          >
            Cancel
          </Button>
        ) : (
          <div />
        )}
      </div>
    </div>
  );
}
