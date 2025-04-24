import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Layout, Menu, Button, Drawer, Tooltip } from 'antd';
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
  DesktopOutlined,
  UpCircleOutlined
} from '@ant-design/icons';
import { responsive } from '../styles/theme';
import LanguageSelector from '../components/LanguageSelector';
import { server } from '../api';
import { keyframes } from 'styled-components';

const { Header, Content } = Layout;

const pulse = keyframes`
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
`;

const float = keyframes`
  0% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-3px);
  }
  100% {
    transform: translateY(0px);
  }
`;

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
  font-size: 13px;
  font-weight: 400;
  color: #7a9bcf;
  margin-left: 4px;
  opacity: 0.85;
`;

const UpdateTooltip = styled.div`
  font-size: 12px;
  line-height: 1.5;
  
  .version {
    color: #52c41a;
    font-weight: 500;
  }
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

const UpdateIcon = styled(UpCircleOutlined)`
  color: #52c41a;
  font-size: 14px;
  margin-left: 4px;
  cursor: pointer;
  animation: ${pulse} 2s infinite ease-in-out, ${float} 3s infinite ease-in-out;
  transition: all 0.3s;
  
  &:hover {
    color: #389e0d;
    transform: scale(1.3);
  }
`;

const MainLayout: React.FC = () => {
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const [isMobile, setIsMobile] = useState(responsive.isMobile());
  const { t } = useTranslation();
  const [version, setVersion] = useState<string>("");
  const [needUpdate, setNeedUpdate] = useState<boolean>(false);
  const [latestVersion, setLatestVersion] = useState<string>("");
  const hasRequestedRef = useRef(false);

  const fetchVersion = useCallback(async () => {
    try {
      const timestamp = new Date().getTime();
      const versionData = await server.getVersion(`?_t=${timestamp}`);
      if (versionData && versionData.version) {
        console.log('成功获取版本信息:', versionData);
        setVersion(versionData.version);
        setNeedUpdate(versionData.needUpdate);
        setLatestVersion(versionData.latestVersion || "");
      } else {
        console.warn('版本数据不完整:', versionData);
        if (!version) {
          setVersion("0.0.1");
        }
      }
    } catch (error) {
      console.error('获取服务器版本失败:', error);
      if (!version) {
        setVersion("0.0.1");
      }
    }
  }, [version]);

  useEffect(() => {
    const fetchWithDelay = setTimeout(() => {
      if (!hasRequestedRef.current) {
        hasRequestedRef.current = true;
        fetchVersion();
      }
    }, 50);

    return () => clearTimeout(fetchWithDelay);
  }, [fetchVersion]);

  useEffect(() => {
    setIsMobile(responsive.isMobile());
    
    const handleResize = () => {
      setIsMobile(responsive.isMobile());
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const handleUpdateClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    window.open('https://github.com/nannanStrawberry314/license/releases/latest', '_blank');
  };

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

  // 版本更新提示内容
  const updateTooltipContent = (
    <UpdateTooltip>
      发现新版本 <span className="version">v{latestVersion}</span> 可用
      <br />
      点击查看更新详情
    </UpdateTooltip>
  );

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
            {version && (
              <>
                <VersionText>v{version}</VersionText>
                {needUpdate && (
                  <Tooltip title={updateTooltipContent} placement="bottom">
                    <UpdateIcon onClick={handleUpdateClick} />
                  </Tooltip>
                )}
              </>
            )}
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
              {version && (
                <>
                  <span style={{ fontSize: '13px', marginLeft: '4px', color: '#7a9bcf', opacity: 0.85 }}>v{version}</span>
                  {needUpdate && (
                    <Tooltip title={updateTooltipContent} placement="bottom">
                      <UpdateIcon onClick={handleUpdateClick} />
                    </Tooltip>
                  )}
                </>
              )}
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