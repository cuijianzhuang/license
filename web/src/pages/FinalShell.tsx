import React, { useState } from 'react';
import { Typography, Form, Button, Input, Alert } from 'antd';
import styled from 'styled-components';
import PageHeader from '../components/PageHeader';
import ResultCard from '../components/ResultCard';
import { finalshell } from '../api';
import { FinalShellLicense } from '../types';

const { Paragraph } = Typography;

const FormWrapper = styled.div`
  max-width: 600px;
  margin-bottom: 32px;
`;

const FinalShell: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [license, setLicense] = useState<FinalShellLicense | null>(null);
  const [form] = Form.useForm();

  const handleGenerateLicense = async (values: { username: string }) => {
    setLoading(true);
    try {
      const data = await finalshell.generateLicense(values.username);
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
      breadcrumbName: 'FinalShell 注册机',
    },
  ];

  return (
    <div>
      <PageHeader
        title="FinalShell 注册机"
        subTitle="生成FinalShell SSH工具的注册码"
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        FinalShell是一款优秀的SSH客户端工具，填写以下表单生成FinalShell的注册码，解锁所有专业功能。
      </Paragraph>

      <Alert
        message="注意事项"
        description="生成的注册码仅供学习和测试使用，请支持正版软件。"
        type="warning"
        showIcon
        style={{ marginBottom: 24 }}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
            initialValue="FinalShell_User"
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              生成注册码
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {license && (
        <ResultCard
          title="FinalShell注册码生成成功"
          data={{
            '用户名': license.username || '未指定',
            '注册码': license.license,
          }}
          fileName="finalshell-license.txt"
        />
      )}

      <Alert
        message="使用说明"
        description={
          <div>
            <p>1. 打开FinalShell软件</p>
            <p>2. 点击"帮助" &gt; "注册"</p>
            <p>3. 输入上面生成的用户名和注册码</p>
            <p>4. 点击"确定"完成注册</p>
          </div>
        }
        type="info"
        showIcon
        style={{ marginTop: 24 }}
      />
    </div>
  );
};

export default FinalShell; 