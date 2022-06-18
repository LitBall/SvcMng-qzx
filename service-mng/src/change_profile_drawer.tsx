import React from 'react';
import { FormInstance } from 'antd/lib/form';
import { Drawer, Form, Input, Button } from 'antd';
import { get, post } from './api/request';
import { API } from "./api/apis";
interface ProfileProps {
  visible: boolean;
  closeDrawer: () => void;
}

interface ProfileRequestData {
  // 这里声明你的请求数据变量，变量名和Form.Item的key要一样
  aws_a_key: string;
  aws_s_key: string;
  ali_a_key: string;
  ali_s_key: string;
  hw_a_key: string;
  hw_s_key: string;
  aws_region: string;
  ali_region: string;
  hw_region: string;
}

interface ProfileResponseData {
  // 这里声明你的响应数据变量
  code: number;
  msg: string;
  data: {
    accesskey: AccessKey;
  };
}

interface AccessKey {
  Aws_a_key: string;
  Aws_s_key: string;
  Aws_region: string;
  Ali_a_key: string;
  Ali_s_key: string;
  Ali_region: string;
  Hw_a_key: string;
  Hw_s_key: string;
  Hw_region: string;
}

export class ProfileDrawer extends React.Component<ProfileProps> {
  public formRef = React.createRef<FormInstance>();

  componentDidMount(){
    this.fill();
  }

  componentWillReceiveProps(nextProps: ProfileProps) {
    this.fill();
  }

  fill = () => {
    get(API.AccessKeyInfo, {  params: { uid: 1} }).then(data => { // 第一个参数为url，第二个参数是get的params，如果没有就传一个空对象
      const res = data as ProfileResponseData;
      if (res.code === 1) {
        const result = res.data.accesskey as AccessKey;
        if (this.formRef.current) {
          this.formRef.current.setFieldsValue({
            // 在表单中填入已经有的数据
            aws_a_key: result.Aws_a_key,
            aws_s_key: result.Aws_s_key,
            ali_a_key: result.Ali_a_key,
            ali_s_key: result.Ali_s_key,
            hw_a_key: result.Hw_a_key,
            hw_s_key: result.Hw_s_key,
            aws_region: result.Aws_region,
            ali_region: result.Ali_region,
            hw_region: result.Hw_region
          });
        }
      } else {
        //错误处理
      }
    })
  }

  submit = () => {
    const { closeDrawer } = this.props;
    const values: ProfileRequestData = this.formRef.current?.getFieldsValue();
    post(API.AccessKeyUpdate+"?uid=1", values).then(data => { // post第一个参数是请求url，在api文件中维护；第二个参数是请求数据，可以用上面的values收集表单数据
      const res = data as ProfileResponseData;
      if (res.code === 1) { // 处理请求数据
        console.log("succ: ", res.data)
      } else {
        console.log("fail: ", res.msg)
      }
      closeDrawer();
    });
  }

  render() {
    const { visible, closeDrawer } = this.props;
    return (
      <Drawer
        title={'管理密钥'}
        width={1000}
        onClose={closeDrawer}
        visible={visible}
        bodyStyle={{ paddingBottom: 80 }}
      >
        <Form
          layout="horizontal"
          requiredMark={true}
          labelCol={{ span: 5 }}
          wrapperCol={{ span: 16 }}
          initialValues={{ remember: true }}
          ref={this.formRef}
        >
          <Form.Item
            label="AWS access_key"
            name="aws_a_key"
            key="aws_a_key"
            rules={[{ required: true, message: '请输入 access_key!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
              label="AWS secret_key"
              name="aws_s_key"
              key="aws_s_key"
              rules={[{ required: true, message: '请输入 secret_key!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
              label="AWS region_id"
              name="aws_region"
              key="aws_region"
              rules={[{ required: true, message: '请输入 region_id!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="Ali access_key"
              name="ali_a_key"
              key="ali_a_key"
              rules={[{ required: true, message: '请输入 access_key!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
              label="Ali secret_key"
              name="ali_s_key"
              key="ali_s_key"
              rules={[{ required: true, message: '请输入 secret_key!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
              label="Ali region_id"
              name="ali_region"
              key="ali_region"
              rules={[{ required: true, message: '请输入 region_id!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="HW access_key"
              name="hw_a_key"
              key="hw_a_key"
              rules={[{ required: true, message: '请输入 access_key!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
              label="HW secret_key"
              name="hw_s_key"
              key="hw_s_key"
              rules={[{ required: true, message: '请输入 secret_key!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
              label="HW region_id"
              name="hw_region"
              key="hw_region"
              rules={[{ required: true, message: '请输入 region_id!' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item wrapperCol={{ offset: 5, span: 16 }} key="button">
            <Button onClick={closeDrawer} style={{ marginRight: 8 }}>
              取消
            </Button>
            <Button onClick={this.submit} type="primary">
              提交
            </Button>
          </Form.Item>
        </Form>
      </Drawer>
    )
  }
}
