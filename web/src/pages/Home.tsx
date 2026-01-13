import React from 'react';
import {Button, Card, Col, Row, Typography} from 'antd';
import {useNavigate} from 'react-router-dom';
import styled from 'styled-components';
import {useTranslation} from 'react-i18next';
import {
    AppstoreOutlined,
    BranchesOutlined,
    CodeOutlined,
    CodeSandboxOutlined,
    DesktopOutlined,
} from '@ant-design/icons';

const { Title, Paragraph } = Typography;

const HeroSection = styled.div`
  text-align: center;
  margin-bottom: 48px;
`;

const HeroTitle = styled(Title)`
  margin-bottom: 16px !important;
`;

const HeroDescription = styled(Paragraph)`
  font-size: 18px;
  max-width: 800px;
  margin: 0 auto 32px;
`;

const ToolsGrid = styled(Row)`
  margin-top: 32px;
`;

const ToolCard = styled(Card)`
  height: 100%;
  border-radius: 12px;
  transition: all 0.3s ease;
  overflow: hidden;
  cursor: pointer;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 25px rgba(37, 99, 235, 0.1);
  }
`;

const IconWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 16px;
  margin-bottom: 16px;
  font-size: 28px;
  background-color: #f0f5ff;
  color: #2563eb;
`;

const ToolTitle = styled(Title)`
  margin-bottom: 8px !important;
`;

const ToolDescription = styled(Paragraph)`
  margin-bottom: 16px;
  color: #6b7280;
`;

const ActionButton = styled(Button)`
  border-radius: 6px;
`;

const Home: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useTranslation();

  const tools = [
    {
      title: t('home.tools.jetbrains.title'),
      description: t('home.tools.jetbrains.description'),
      icon: <CodeOutlined />,
      path: '/page/jetbrains',
      color: '#f0f5ff',
    },
    {
      title: t('home.tools.jrebel.title'),
      description: t('home.tools.jrebel.description'),
      icon: <AppstoreOutlined />,
      path: '/page/jrebel',
      color: '#f5f3ff',
    },
    {
      title: t('home.tools.gitlab.title'),
      description: t('home.tools.gitlab.description'),
      icon: <BranchesOutlined />,
      path: '/page/gitlab',
      color: '#f3f4f6',
    },
    {
      title: t('home.tools.finalshell.title'),
      description: t('home.tools.finalshell.description'),
      icon: <DesktopOutlined />,
      path: '/page/finalshell',
      color: '#f0fdf4',
    },
    {
      title: t('home.tools.mobaxterm.title'),
      description: t('home.tools.mobaxterm.description'),
      icon: <CodeSandboxOutlined />,
      path: '/page/mobaxterm',
      color: '#eff6ff',
    },
  ];

  const handleCardClick = (path: string) => {
    navigate(path);
  };

  return (
    <>
      <HeroSection>
        <HeroTitle level={1}>{t('home.welcome')}</HeroTitle>
        <HeroDescription>
          {t('home.description')}
        </HeroDescription>
      </HeroSection>
      
      <ToolsGrid gutter={[24, 24]}>
        {tools.map((tool, index) => (
          <Col xs={24} sm={12} md={8} key={index}>
            <ToolCard onClick={() => handleCardClick(tool.path)}>
              <IconWrapper style={{ backgroundColor: tool.color }}>
                {tool.icon}
              </IconWrapper>
              <ToolTitle level={4}>{tool.title}</ToolTitle>
              <ToolDescription>{tool.description}</ToolDescription>
              <ActionButton 
                type="primary" 
                onClick={(e) => {
                  e.stopPropagation();
                  navigate(tool.path);
                }}
              >
                {t('common.useNow')}
              </ActionButton>
            </ToolCard>
          </Col>
        ))}
      </ToolsGrid>
    </>
  );
};

export default Home; 