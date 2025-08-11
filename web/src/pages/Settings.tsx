import React from 'react';
import {Button, Card, Form, Input, Select, Switch, Typography} from 'antd';
import styled from 'styled-components';
import PageHeader from '../components/PageHeader';

const { Title, Paragraph } = Typography;
const { Option } = Select;

const SettingCard = styled(Card)`
  margin-bottom: 24px;
`;

const Settings: React.FC = () => {
  const [form] = Form.useForm();

  const handleSaveSettings = (values: any) => {
    console.log('保存设置:', values);
    // 实际保存将在API集成后实现
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: '首页',
    },
    {
      path: '',
      breadcrumbName: '设置',
    },
  ];

  return (
    <div>
      <PageHeader
        title="系统设置"
        subTitle="配置许可证生成服务的相关参数"
        breadcrumbs={breadcrumbs}
      />

      <Form
        form={form}
        layout="vertical"
        onFinish={handleSaveSettings}
        initialValues={{
          apiBaseUrl: '/api',
          theme: 'light',
          language: 'zh-CN',
          enableNotifications: true,
          saveHistory: true,
        }}
      >
        <SettingCard title="基本设置">
          <Form.Item
            name="apiBaseUrl"
            label="API基础URL"
            rules={[{ required: true, message: '请输入API基础URL' }]}
          >
            <Input placeholder="请输入API基础URL" />
          </Form.Item>

          <Form.Item name="theme" label="主题">
            <Select>
              <Option value="light">浅色</Option>
              <Option value="dark">深色</Option>
              <Option value="system">跟随系统</Option>
            </Select>
          </Form.Item>

          <Form.Item name="language" label="语言">
            <Select>
              <Option value="zh-CN">简体中文</Option>
              <Option value="en-US">English</Option>
            </Select>
          </Form.Item>
        </SettingCard>

        <SettingCard title="通知设置">
          <Form.Item
            name="enableNotifications"
            label="启用通知"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
          
          <Paragraph type="secondary">
            启用后，系统将在生成许可证成功或失败时显示通知。
          </Paragraph>
        </SettingCard>

        <SettingCard title="历史记录">
          <Form.Item name="saveHistory" label="保存历史记录" valuePropName="checked">
            <Switch />
          </Form.Item>
          
          <Paragraph type="secondary">
            启用后，系统将保存生成的许可证历史记录。
          </Paragraph>

          <Form.Item>
            <Button danger>清除历史记录</Button>
          </Form.Item>
        </SettingCard>

        <Form.Item>
          <Button type="primary" htmlType="submit">
            保存设置
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default Settings; 