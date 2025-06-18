import type { Notifier } from '../../../../../entity/notifiers';

interface Props {
  notifier: Notifier;
}

export function ShowWebhookNotifierComponent({ notifier }: Props) {
  return (
    <>
      <div className="flex items-center">
        <div className="min-w-[110px]">Webhook URL</div>

        <div className="w-[250px]">{notifier?.webhookNotifier?.webhookUrl || '-'}</div>
      </div>

      <div className="mt-1 mb-1 flex items-center">
        <div className="min-w-[110px]">Method</div>
        <div>{notifier?.webhookNotifier?.webhookMethod || '-'}</div>
      </div>
    </>
  );
}
