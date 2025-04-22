import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ConfigProvider, App as AntApp } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';
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
  // 设置页面标题
  useEffect(() => {
    document.title = '软件许可证生成服务';
  }, []);

  return (
    <ConfigProvider locale={zhCN} theme={theme}>
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
