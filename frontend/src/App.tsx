import { useEffect, useState } from 'react';
import { BrowserRouter, Route } from 'react-router';
import { Routes } from 'react-router';

import { userApi } from './entity/users';
import { AuthPageComponent } from './pages/AuthPageComponent';
import { MainScreenComponent } from './widgets/main/MainScreenComponent';

function App() {
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    const isAuthorized = userApi.isAuthorized();
    setIsAuthorized(isAuthorized);

    userApi.addAuthListener(() => {
      setIsAuthorized(userApi.isAuthorized());
    });
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={!isAuthorized ? <AuthPageComponent /> : <MainScreenComponent />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
