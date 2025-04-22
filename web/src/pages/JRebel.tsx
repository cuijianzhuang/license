import React, { useState } from 'react';
import { Typography, Form, Button, Input, Alert, Divider } from 'antd';
import styled from 'styled-components';
import PageHeader from '../components/PageHeader';
import ResultCard from '../components/ResultCard';
import { jrebel } from '../api';
import { JRebelLicense } from '../types';

const { Title, Paragraph } = Typography;

const FormWrapper = styled.div`
  max-width: 600px;
  margin-bottom: 32px;
`;

const StepItem = styled.div`
  margin-bottom: 16px;
`;

const StepNumber = styled.span`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background-color: #1890ff;
  color: #fff;
  border-radius: 50%;
  margin-right: 8px;
  font-size: 14px;
`;

const JRebel: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [license, setLicense] = useState<JRebelLicense | null>(null);
  const [form] = Form.useForm();

  const handleGenerateLicense = async (values: {
    username: string;
    email: string;
    teamName: string;
  }) => {
    setLoading(true);
    try {
      const data = await jrebel.generateLicense(
        values.username,
        values.email,
        values.teamName
      );
      setLicense(data);
    } catch (error) {
      console.error('生成许可证失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: '首页',
    },
    {
      path: '',
      breadcrumbName: 'JRebel 激活服务',
    },
  ];

  const serverUrl = 'http://jrebel.license.server/';

  return (
    <div>
      <PageHeader
        title="JRebel 激活服务"
        subTitle="为JRebel热部署工具提供激活服务"
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        JRebel是一款强大的Java热部署工具，可以让你在不重启应用服务器的情况下，实时看到代码改动的效果。
      </Paragraph>

      <Alert
        message="激活说明"
        description="本服务提供JRebel的激活。请填写以下表单来获取激活信息，或使用本服务器地址进行离线激活。"
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <Title level={3}>方法一：离线激活</Title>
      <Paragraph>
        您可以直接使用以下服务器地址在JRebel中进行离线激活：
      </Paragraph>

      <ResultCard
        title="JRebel激活服务器"
        data={{
          '服务器地址': serverUrl,
        }}
      />

      <Divider />

      <Title level={3}>方法二：生成授权信息</Title>
      <Paragraph>
        填写以下表单生成JRebel的授权信息：
      </Paragraph>

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="email"
            label="邮箱地址"
            rules={[
              { required: true, message: '请输入邮箱地址' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input placeholder="请输入邮箱地址" />
          </Form.Item>

          <Form.Item
            name="teamName"
            label="团队名称"
            rules={[{ required: true, message: '请输入团队名称' }]}
            initialValue="JRebel Team"
          >
            <Input placeholder="请输入团队名称" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              生成授权信息
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {license && (
        <ResultCard
          title="JRebel授权信息生成成功"
          data={{
            '用户名': license.username || '未指定',
            '邮箱': license.email || '未指定',
            '团队名称': license.teamName || '未指定',
            '授权密钥': license.validKey || '',
          }}
          fileName="jrebel-license.txt"
        />
      )}

      <Divider />

      <Title level={3}>JRebel激活步骤</Title>

      <div>
        <StepItem>
          <StepNumber>1</StepNumber>
          <span>打开IDE (如IntelliJ IDEA)</span>
        </StepItem>
        <StepItem>
          <StepNumber>2</StepNumber>
          <span>找到JRebel插件的设置界面</span>
        </StepItem>
        <StepItem>
          <StepNumber>3</StepNumber>
          <span>选择"Team URL"激活方式</span>
        </StepItem>
        <StepItem>
          <StepNumber>4</StepNumber>
          <span>在URL中填入上方的服务器地址</span>
        </StepItem>
        <StepItem>
          <StepNumber>5</StepNumber>
          <span>如果生成了授权信息，可填入相关的用户名和授权密钥</span>
        </StepItem>
        <StepItem>
          <StepNumber>6</StepNumber>
          <span>点击"Activate"完成激活</span>
        </StepItem>
      </div>
    </div>
  );
};

export default JRebel; 