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

// Helper function to map language code to HTML lang attribute
const mapToHtmlLang = (language: string): string => {
  if (language.startsWith('zh-CN') || language === 'zh-Hans' || language === 'zh_CN') return 'zh-CN';
  if (language.startsWith('zh-TW') || language === 'zh-Hant' || language === 'zh_TW') return 'zh-TW';
  if (language.startsWith('ja')) return 'ja';
  if (language.startsWith('ko')) return 'ko';
  if (language.startsWith('ru')) return 'ru';
  return 'en';
};

const App: React.FC = () => {
  const { t, i18n } = useTranslation();
  
  // Set page title and HTML lang attribute based on current language
  useEffect(() => {
    document.title = t('app.title');
    document.documentElement.lang = mapToHtmlLang(i18n.language);
  }, [t, i18n.language]);

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
