import React from 'react';
import { Typography, Row, Col, Card, Button } from 'antd';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import {
  CodeOutlined,
  BranchesOutlined,
  DesktopOutlined,
  CodeSandboxOutlined,
  AppstoreOutlined,
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

  const tools = [
    {
      title: 'JetBrains 激活工具',
      description: '生成JetBrains全系列产品的激活码，包括IntelliJ IDEA, PyCharm, WebStorm等',
      icon: <CodeOutlined />,
      path: '/jetbrains',
      color: '#f0f5ff',
    },
    {
      title: 'GitLab 许可证',
      description: '为GitLab创建企业版许可证，解锁所有高级功能',
      icon: <BranchesOutlined />,
      path: '/gitlab',
      color: '#f3f4f6',
    },
    {
      title: 'FinalShell 注册机',
      description: '生成FinalShell SSH工具的注册码',
      icon: <DesktopOutlined />,
      path: '/finalshell',
      color: '#f0fdf4',
    },
    {
      title: 'MobaXterm 激活工具',
      description: '解锁MobaXterm高级功能，生成专业版激活码',
      icon: <CodeSandboxOutlined />,
      path: '/mobaxterm',
      color: '#eff6ff',
    },
    {
      title: 'JRebel 激活服务',
      description: '提供JRebel热部署工具激活服务',
      icon: <AppstoreOutlined />,
      path: '/jrebel',
      color: '#f5f3ff',
    },
  ];

  const handleCardClick = (path: string) => {
    navigate(path);
  };

  return (
    <>
      <HeroSection>
        <HeroTitle level={1}>软件许可证生成服务</HeroTitle>
        <HeroDescription>
          一站式解决开发工具的许可证需求，支持各类常用开发工具的激活服务
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
                立即使用
              </ActionButton>
            </ToolCard>
          </Col>
        ))}
      </ToolsGrid>
    </>
  );
};

export default Home; 