import { useState } from 'react';

import { DatabasesComponent } from '../../features/databases/ui/DatabasesComponent';
import { NotifiersComponent } from '../../features/notifiers/ui/NotifiersComponent';
import { StoragesComponent } from '../../features/storages/StoragesComponent';
import { useScreenHeight } from '../../shared/hooks';

export const MainScreenComponent = () => {
  const screenHeight = useScreenHeight();

  const [selectedTab, setSelectedTab] = useState<'notifiers' | 'storages' | 'databases'>(
    'databases',
  );

  const contentHeight = screenHeight - 95;

  return (
    <div style={{ height: screenHeight }} className="bg-[#f5f5f5] p-3">
      {/* ===================== NAVBAR ===================== */}
      <div className="mb-3 flex h-[60px] items-center rounded bg-white p-3 shadow">
        <div className="flex items-center gap-3 hover:opacity-80">
          <a href="https://postgresus.com" target="_blank" rel="noreferrer">
            <img className="h-[35px] w-[35px]" src="/logo.svg" />
          </a>

          <div className="text-xl font-bold">
            <a
              href="https://postgresus.com"
              className="text-black"
              target="_blank"
              rel="noreferrer"
            >
              Postgresus
            </a>
          </div>
        </div>

        <div className="mr-3 ml-auto flex gap-5">
          <a
            className="hover:opacity-80"
            href="https://postgresus.com/community"
            target="_blank"
            rel="noreferrer"
          >
            Community
          </a>

          <a
            className="hover:opacity-80"
            href="https://github.com/postgresus/postgresus"
            target="_blank"
            rel="noreferrer"
          >
            GitHub
          </a>
        </div>
      </div>
      {/* ===================== END NAVBAR ===================== */}

      <div className="flex">
        <div
          className="max-w-[60px] min-w-[60px] rounded bg-white py-2 shadow"
          style={{ height: contentHeight }}
        >
          {[
            {
              text: 'Databases',
              name: 'databases',
              icon: '/icons/menu/database-gray.svg',
              selectedIcon: '/icons/menu/database-white.svg',
              onClick: () => setSelectedTab('databases'),
            },
            {
              text: 'Storages',
              name: 'storages',
              icon: '/icons/menu/storage-gray.svg',
              selectedIcon: '/icons/menu/storage-white.svg',
              onClick: () => setSelectedTab('storages'),
            },
            {
              text: 'Notifiers',
              name: 'notifiers',
              icon: '/icons/menu/notifier-gray.svg',
              selectedIcon: '/icons/menu/notifier-white.svg',
              onClick: () => setSelectedTab('notifiers'),
            },
          ].map((tab) => (
            <div key={tab.text} className="flex justify-center">
              <div
                className={`flex h-[50px] w-[50px] cursor-pointer items-center justify-center rounded ${selectedTab === tab.name ? 'bg-blue-600' : 'hover:bg-blue-50'}`}
                onClick={tab.onClick}
              >
                <div className="mb-1">
                  <div className="flex justify-center">
                    <img
                      src={selectedTab === tab.name ? tab.selectedIcon : tab.icon}
                      width={20}
                      alt={tab.text}
                      loading="lazy"
                    />
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {selectedTab === 'notifiers' && <NotifiersComponent contentHeight={contentHeight} />}
        {selectedTab === 'storages' && <StoragesComponent contentHeight={contentHeight} />}
        {selectedTab === 'databases' && <DatabasesComponent contentHeight={contentHeight} />}
      </div>
    </div>
  );
};
