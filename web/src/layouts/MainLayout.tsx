import React, { useState, useEffect } from 'react';
import { Layout, Menu, Button, Drawer } from 'antd';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { 
  MenuUnfoldOutlined,
  AppstoreOutlined,
  CodeOutlined,
  BranchesOutlined,
  CodeSandboxOutlined,
  HomeOutlined,
  DesktopOutlined
} from '@ant-design/icons';
import { responsive } from '../styles/theme';

const { Header, Content } = Layout;

const StyledLayout = styled(Layout)`
  min-height: 100vh;
  background: #f5f7fa;
`;

const StyledHeader = styled(Header)`
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 1;
  backdrop-filter: blur(8px);
`;

const Logo = styled.div`
  color: #1890ff;
  font-size: 20px;
  font-weight: bold;
  margin-right: 48px;
`;

const StyledMenu = styled(Menu)`
  flex: 1;
  border-bottom: none;
`;

const MobileMenuButton = styled(Button)`
  display: none;
  @media (max-width: 768px) {
    display: block;
    margin-right: 16px;
  }
`;

const DesktopMenu = styled.div`
  @media (max-width: 768px) {
    display: none;
  }
`;

const MainContent = styled(Content)`
  padding: 24px;
  margin: 16px;
  background: white;
  border-radius: 8px;
  
  @media (max-width: 768px) {
    margin: 8px;
    padding: 16px;
  }
`;

const MainLayout: React.FC = () => {
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const [isMobile, setIsMobile] = useState(responsive.isMobile());

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(responsive.isMobile());
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: '/jetbrains',
      icon: <CodeOutlined />,
      label: 'JetBrains',
    },
    {
      key: '/gitlab',
      icon: <BranchesOutlined />,
      label: 'GitLab',
    },
    {
      key: '/finalshell',
      icon: <DesktopOutlined />,
      label: 'FinalShell',
    },
    {
      key: '/mobaxterm',
      icon: <CodeSandboxOutlined />,
      label: 'MobaXterm',
    },
    {
      key: '/jrebel',
      icon: <AppstoreOutlined />,
      label: 'JRebel',
    },
  ];

  const onMenuClick = (path: string) => {
    navigate(path);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  return (
    <StyledLayout>
      <StyledHeader>
        <MobileMenuButton
          type="text"
          icon={<MenuUnfoldOutlined />}
          onClick={() => setMobileOpen(true)}
        />
        <Logo>License Generator</Logo>
        <DesktopMenu>
          <StyledMenu 
            mode="horizontal" 
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => onMenuClick(key)}
          />
        </DesktopMenu>
      </StyledHeader>
      
      <Drawer
        title="选择工具"
        placement="left"
        onClose={() => setMobileOpen(false)}
        open={mobileOpen}
        bodyStyle={{ padding: 0 }}
      >
        <Menu 
          mode="inline" 
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => onMenuClick(key)}
          style={{ border: 'none' }}
        />
      </Drawer>
      
      <MainContent>
        <Outlet />
      </MainContent>
    </StyledLayout>
  );
};

export default MainLayout; 