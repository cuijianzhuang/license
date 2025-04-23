import React, { useState, useEffect } from 'react';
import { Layout, Menu, Button, Drawer } from 'antd';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
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
import LanguageSelector from '../components/LanguageSelector';
import { server } from '../api';

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

const LogoWrapper = styled.div`
  display: flex;
  align-items: center;
  margin-right: 48px;
  cursor: pointer;
`;

const LogoImage = styled.img`
  height: 32px;
  margin-right: 8px;
`;

const LogoText = styled.span`
  color: #1890ff;
  font-size: 20px;
  font-weight: bold;
  display: flex;
  align-items: baseline;
`;

const VersionText = styled.span`
  font-size: 11px;
  font-weight: 400;
  color: #7a9bcf;
  margin-left: 4px;
  opacity: 0.85;
`;

const StyledMenu = styled(Menu)`
  flex: 1;
  border-bottom: none;
`;

const MobileMenuButton = styled(Button)<{ $isMobile: boolean }>`
  display: ${props => props.$isMobile ? 'block' : 'none'};
  margin-right: 16px;
`;

const DesktopMenu = styled.div<{ $isMobile: boolean }>`
  display: ${props => props.$isMobile ? 'none' : 'block'};
  flex: 1;
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

const HeaderControls = styled.div`
  display: flex;
  align-items: center;
  margin-left: auto;
`;

const MainLayout: React.FC = () => {
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const [isMobile, setIsMobile] = useState(responsive.isMobile());
  const { t } = useTranslation();
  const [version, setVersion] = useState<string>("");

  useEffect(() => {
    // Fetch server version
    const fetchVersion = async () => {
      try {
        const versionData = await server.getVersion();
        if (versionData) {
          setVersion(versionData);
        }
      } catch (error) {
        console.error('Failed to fetch server version:', error);
      }
    };

    fetchVersion();
  }, []);

  useEffect(() => {
    setIsMobile(responsive.isMobile());
    
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
      label: t('nav.home'),
    },
    {
      key: '/jetbrains',
      icon: <CodeOutlined />,
      label: t('nav.jetbrains'),
    },
    {
      key: '/jrebel',
      icon: <AppstoreOutlined />,
      label: t('nav.jrebel'),
    },
    {
      key: '/gitlab',
      icon: <BranchesOutlined />,
      label: t('nav.gitlab'),
    },
    {
      key: '/finalshell',
      icon: <DesktopOutlined />,
      label: t('nav.finalshell'),
    },
    {
      key: '/mobaxterm',
      icon: <CodeSandboxOutlined />,
      label: t('nav.mobaxterm'),
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
          $isMobile={isMobile}
          type="text"
          icon={<MenuUnfoldOutlined />}
          onClick={() => setMobileOpen(true)}
        />
        <LogoWrapper onClick={() => navigate('/')}>
          <LogoImage src="/logo.svg" alt="License" />
          <LogoText>
            License
            {version && <VersionText>v{version}</VersionText>}
          </LogoText>
        </LogoWrapper>
        <DesktopMenu $isMobile={isMobile}>
          <StyledMenu 
            mode="horizontal" 
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => onMenuClick(key)}
          />
        </DesktopMenu>
        <HeaderControls>
          <LanguageSelector />
        </HeaderControls>
      </StyledHeader>
      
      <Drawer
        title={
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <img src="/logo.svg" alt="License" style={{ height: '28px', marginRight: '8px' }} />
            <span style={{ display: 'flex', alignItems: 'baseline' }}>
              License
              {version && <span style={{ fontSize: '11px', marginLeft: '4px', color: '#7a9bcf', opacity: 0.85 }}>v{version}</span>}
            </span>
          </div>
        }
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