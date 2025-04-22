import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ConfigProvider, App as AntApp } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';
import enUS from 'antd/lib/locale/en_US';
import { useTranslation } from 'react-i18next';
import MainLayout from './layouts/MainLayout';
import Home from './pages/Home';
import JetBrains from './pages/JetBrains';
import GitLab from './pages/GitLab';
import FinalShell from './pages/FinalShell';
import MobaXterm from './pages/MobaXterm';
import JRebel from './pages/JRebel';
import GlobalStyles from './styles/GlobalStyles';
import { theme } from './styles/theme';

const App: React.FC = () => {
  const { t, i18n } = useTranslation();
  
  // Set page title based on current language
  useEffect(() => {
    document.title = t('app.title');
  }, [t]);

  // Get antd locale based on current language
  const getAntdLocale = () => {
    return i18n.language.startsWith('zh') ? zhCN : enUS;
  };

  return (
    <ConfigProvider locale={getAntdLocale()} theme={theme}>
      <AntApp>
        <GlobalStyles />
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<MainLayout />}>
              <Route index element={<Home />} />
              <Route path="jetbrains" element={<JetBrains />} />
              <Route path="gitlab" element={<GitLab />} />
              <Route path="finalshell" element={<FinalShell />} />
              <Route path="mobaxterm" element={<MobaXterm />} />
              <Route path="jrebel" element={<JRebel />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </AntApp>
    </ConfigProvider>
  );
}

export default App;
