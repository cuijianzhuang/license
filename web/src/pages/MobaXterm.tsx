import React, { useState } from 'react';
import { Typography, Form, Button, Input, Select, Alert } from 'antd';
import styled from 'styled-components';
import PageHeader from '../components/PageHeader';
import ResultCard from '../components/ResultCard';
import { mobaxterm } from '../api';
import { MobaXtermLicense } from '../types';

const { Paragraph } = Typography;
const { Option } = Select;

const FormWrapper = styled.div`
  max-width: 600px;
  margin-bottom: 32px;
`;

const versions = [
  '23.0',
  '22.3',
  '22.2',
  '22.1',
  '22.0',
  '21.5',
  '21.4',
  '21.3',
  '21.2',
  '21.1',
  '21.0',
  '20.6',
  '20.5',
  '20.4',
  '20.3',
  '20.2',
  '20.1',
  '20.0',
];

const MobaXterm: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [license, setLicense] = useState<MobaXtermLicense | null>(null);
  const [form] = Form.useForm();

  const handleGenerateLicense = async (values: { username: string; version: string }) => {
    setLoading(true);
    try {
      const data = await mobaxterm.generateLicense(values.username, values.version);
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
      breadcrumbName: 'MobaXterm 激活工具',
    },
  ];

  return (
    <div>
      <PageHeader
        title="MobaXterm 激活工具"
        subTitle="生成MobaXterm专业版激活码"
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        MobaXterm是一款功能强大的终端工具，填写以下表单生成MobaXterm的专业版激活码。
      </Paragraph>

      <Alert
        message="注意事项"
        description="生成的激活码仅供学习和测试使用，请支持正版软件。"
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
            initialValue="MobaXterm User"
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="version"
            label="软件版本"
            rules={[{ required: true, message: '请选择MobaXterm版本' }]}
            initialValue={versions[0]}
          >
            <Select placeholder="请选择版本">
              {versions.map((version) => (
                <Option key={version} value={version}>
                  {version}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              生成激活码
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {license && (
        <ResultCard
          title="MobaXterm激活码生成成功"
          data={{
            '用户名': license.username || '未指定',
            '版本': license.version || '未指定',
            '激活码': license.license,
          }}
          fileName="mobaxterm-license.txt"
        />
      )}

      <Alert
        message="使用说明"
        description={
          <div>
            <p>1. 打开MobaXterm软件</p>
            <p>2. 点击右上角的"?"按钮，然后选择"Register"</p>
            <p>3. 输入上面生成的用户名和激活码</p>
            <p>4. 点击"OK"完成激活</p>
          </div>
        }
        type="info"
        showIcon
        style={{ marginTop: 24 }}
      />
    </div>
  );
};

export default MobaXterm; 