import React from 'react';
import { FormInstance } from 'antd/lib/form';
import {Drawer, Form, Input, Button, Modal} from 'antd';
import { post } from './api/request';
import {API} from "./api/apis";

interface TemplateProps {
  visible: boolean;
  closeDrawer: () => void;
}

interface TemplateRequestData {
  // 这里声明你的请求数据变量，变量名和Form.Item的key要一样
  tid: number;
  resource_type: string;
  bandwidth: string;
  number_of: number;
  last_time: number;
  image_id: string;
  instance_type: string;
  instance_name: string;
  description: string;
  tags_name: string;
}

interface TemplateResponseData {
  // 这里声明你的响应数据变量
  code: number;
  msg: string;
  data: {
    rid: number;
    resource_type: string;
    bandwidth: string;
    number_of: number;
    last_time: number;
    image_id: string;
    instance_type: string;
    instance_name: string;
    region_id: string;
  };
}

interface Resource {
  rid: number;
  resource_type: string;
  bandwidth: string;
  number_of: number;
  last_time: number;
  image_id: string;
  instance_type: string;
  instance_name: string;
  region_id: string;
}

export class TemplateDrawer extends React.Component<TemplateProps> {
  public formRef = React.createRef<FormInstance>();

  submit = () => {
    const { closeDrawer } = this.props;
    const values: TemplateRequestData = this.formRef.current?.getFieldsValue();
    let msg = "";
    console.log("values: ", values);
    let url = API.Host + "/ServiceManagement/addCloudInstance";
    if (values.tid != undefined) {
      url += "?tid=" + values.tid
    }
    post(url, values).then(data => { // post第一个参数是请求url，在api文件中维护；第二个参数是请求数据，可以用上面的values收集表单数据
      const res = data as TemplateResponseData;
      if (res.code === 1) { // 处理请求数据
        console.log("succ: ", res.msg)
        msg = res.msg + "可至官网查阅对应虚拟机资源，名称为：" + values.instance_name
      } else {
        // 错误处理
        console.log("fail: ", res)
        msg = res.msg
      }
      Modal.info({
        title: '部署得怎么样啦？',
        content: (
            <div>
              <p>{msg}</p>
            </div>
        ),
        onOk() {},
      });
      closeDrawer();
    });
  }

  render() {
    const { visible, closeDrawer } = this.props;
    // Form.Item 自定义自己需要的参数，注意key要和请求数据中的变量名一致
    return (
      <Drawer
        title={'使用模板创建实例'}
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
            label="tid"
            name="tid"
            key="tid"
            rules={[{ required: true, message: '请输入 tid!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="resource_type"
              name="resource_type"
              key="resource_type"
              rules={[{ required: false, message: '请输入 resource_type!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="bandwidth"
              name="bandwidth"
              key="bandwidth"
              rules={[{ required: false, message: '请输入 bandwidth!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="number_of"
              name="number_of"
              key="number_of"
              rules={[{ required: false, message: '请输入 number_of!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="last_time"
              name="last_time"
              key="last_time"
              rules={[{ required: false, message: '请输入 last_time!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="image_id"
              name="image_id"
              key="image_id"
              rules={[{ required: false, message: '请输入 image_id!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="instance_type"
              name="instance_type"
              key="instance_type"
              rules={[{ required: false, message: '请输入 instance_type!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="instance_name"
              name="instance_name"
              key="instance_name"
              rules={[{ required: true, message: '请输入 instance_name!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="description"
              name="description"
              key="description"
              rules={[{ required: false, message: '请输入 description!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
              label="tags_name"
              name="tags_name"
              key="tags_name"
              rules={[{ required: true, message: '请输入 tags_name!' }]}
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
