import { InfoCircleOutlined } from '@ant-design/icons';
import { Input, Tooltip } from 'antd';
import { useState } from 'react';

import type { Notifier } from '../../../../../entity/notifiers';

interface Props {
  notifier: Notifier;
  setNotifier: (notifier: Notifier) => void;
  setIsUnsaved: (isUnsaved: boolean) => void;
}

export function EditTelegramNotifierComponent({ notifier, setNotifier, setIsUnsaved }: Props) {
  const [isShowHowToGetChatId, setIsShowHowToGetChatId] = useState(false);

  return (
    <>
      <div className="flex items-center">
        <div className="min-w-[110px]">Bot token</div>

        <div className="w-[250px]">
          <Input
            value={notifier?.telegramNotifier?.botToken || ''}
            onChange={(e) => {
              if (!notifier?.telegramNotifier) return;
              setNotifier({
                ...notifier,
                telegramNotifier: {
                  ...notifier.telegramNotifier,
                  botToken: e.target.value.trim(),
                },
              });
              setIsUnsaved(true);
            }}
            size="small"
            className="w-full"
            placeholder="1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ"
          />
        </div>
      </div>

      <div className="mb-1 ml-[110px]">
        <a
          className="text-xs !text-blue-600"
          href="https://www.siteguarding.com/en/how-to-get-telegram-bot-api-token"
          target="_blank"
          rel="noreferrer"
        >
          How to get Telegram bot API token?
        </a>
      </div>

      <div className="mb-1 flex items-center">
        <div className="min-w-[110px]">Target chat ID</div>

        <div className="w-[250px]">
          <Input
            value={notifier?.telegramNotifier?.targetChatId || ''}
            onChange={(e) => {
              if (!notifier?.telegramNotifier) return;

              setNotifier({
                ...notifier,
                telegramNotifier: {
                  ...notifier.telegramNotifier,
                  targetChatId: e.target.value.trim(),
                },
              });
              setIsUnsaved(true);
            }}
            size="small"
            className="w-full"
            placeholder="-1001234567890"
          />
        </div>

        <Tooltip
          className="cursor-pointer"
          title="The chat where you want to receive the message (it can be your private chat or a group)"
        >
          <InfoCircleOutlined className="ml-2" style={{ color: 'gray' }} />
        </Tooltip>
      </div>

      <div className="ml-[110px] max-w-[250px]">
        {!isShowHowToGetChatId ? (
          <div
            className="mt-1 cursor-pointer text-xs text-blue-600"
            onClick={() => setIsShowHowToGetChatId(true)}
          >
            How to get Telegram chat ID?
          </div>
        ) : (
          <div className="mt-1 text-xs text-gray-500">
            To get your chat ID, message{' '}
            <a href="https://t.me/getmyid_bot" target="_blank" rel="noreferrer">
              @getmyid_bot
            </a>{' '}
            in Telegram. <u>Make sure you started chat with the bot</u>
            <br />
            <br />
            If you want to get chat ID of a group, add your bot with{' '}
            <a href="https://t.me/getmyid_bot" target="_blank" rel="noreferrer">
              @getmyid_bot
            </a>{' '}
            to the group and write /start (you will see chat ID)
          </div>
        )}
      </div>
    </>
  );
}
