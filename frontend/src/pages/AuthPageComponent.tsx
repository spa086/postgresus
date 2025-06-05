import { Spin } from 'antd';
import { useEffect, useState } from 'react';

import { userApi } from '../entity/users';
import { SignInComponent } from '../features/users';
import { SignUpComponent } from '../features/users';
import { AuthNavbarComponent } from '../features/users';

export function AuthPageComponent() {
  const [isAnyUserExists, setIsAnyUserExists] = useState(false);
  const [isLoading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);

    userApi
      .isAnyUserExists()
      .then((isAnyUserExists) => {
        setIsAnyUserExists(isAnyUserExists);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return (
    <div>
      {isLoading ? (
        <div className="flex h-screen w-screen items-center justify-center">
          <Spin spinning={isLoading} />
        </div>
      ) : (
        <div>
          <div>
            <AuthNavbarComponent />

            <div className="mt-[20vh] flex justify-center">
              {isAnyUserExists ? <SignInComponent /> : <SignUpComponent />}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
