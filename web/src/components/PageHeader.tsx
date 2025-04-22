import React from 'react';
import styled from 'styled-components';
import { Typography, Breadcrumb } from 'antd';
import { Link } from 'react-router-dom';
import { HomeOutlined } from '@ant-design/icons';

const { Title } = Typography;

interface PageHeaderProps {
  title: string;
  subTitle?: string;
  breadcrumbs?: Array<{
    path: string;
    breadcrumbName: string;
  }>;
}

const HeaderContainer = styled.div`
  margin-bottom: 32px;
`;

const TitleContainer = styled.div`
  display: flex;
  align-items: baseline;
  margin-bottom: 12px;
`;

const StyledTitle = styled(Title)`
  font-weight: 700 !important;
  color: #111827 !important;
  margin-bottom: 0 !important;
`;

const SubTitle = styled.span`
  font-size: 16px;
  font-weight: 400;
  color: #6b7280;
  margin-left: 12px;
`;

const StyledBreadcrumb = styled(Breadcrumb)`
  margin-bottom: 16px;
  font-size: 13px;
  
  .ant-breadcrumb-link a {
    color: #4b5563;
    
    &:hover {
      color: #2563eb;
    }
  }
  
  .ant-breadcrumb-separator {
    color: #9ca3af;
  }
`;

const PageHeader: React.FC<PageHeaderProps> = ({ title, subTitle, breadcrumbs }) => {
  return (
    <HeaderContainer>
      {breadcrumbs && breadcrumbs.length > 0 && (
        <StyledBreadcrumb>
          <Breadcrumb.Item>
            <Link to="/">
              <HomeOutlined />
            </Link>
          </Breadcrumb.Item>
          {breadcrumbs.map((item, index) => (
            <Breadcrumb.Item key={index}>
              {item.path ? (
                <Link to={item.path}>{item.breadcrumbName}</Link>
              ) : (
                item.breadcrumbName
              )}
            </Breadcrumb.Item>
          ))}
        </StyledBreadcrumb>
      )}
      <TitleContainer>
        <StyledTitle level={1}>
          {title}
        </StyledTitle>
        {subTitle && <SubTitle>{subTitle}</SubTitle>}
      </TitleContainer>
    </HeaderContainer>
  );
};

export default PageHeader; 