import React from "react";
import { Card, Descriptions, Empty, Button, Modal } from 'antd';
import { get } from './api/request';
import {API} from "./api/apis";

interface TemplateProps {

}

interface TemplateInfo {
  // 这里声明模板需要展示的东西
  Name: string;
  Tid: number;
  ResourceType: string;
  Bandwidth: string;
  NumberOf: number;
  LastTime: number;
  ImageId: string;
  InstanceType: string;
}

interface TemplateRequestData {
  // 这里声明你的请求数据变量
  // aws_a_key: string;
}

interface TemplateResponseData {
  // 这里声明你的响应数据变量
  code: number;
  msg: string;
  data: {
    template: TemplateInfo[]
  };
}

interface TemplateState {
  // 定义 state，把你从后端拿到的数据声明在这里
  xxx: TemplateInfo[] | null; // 一个 TemplateInfo 类型的数组
}


export class Template extends React.Component<TemplateProps, TemplateState> {

  constructor(props: TemplateProps) {
    super(props);
    this.state = {
      xxx: null,
    }

    this.getTemplate();
  }


  getTemplate = () => {
    get(API.TemplateList, { uid:1}).then(data => {
      const res = data as TemplateResponseData;
      if (res.code === 1) {
        // 把响应数据设置到 state 里面，这样才能驱动页面更新
        const values = res.data.template
        this.setState({
          xxx: values,
        })
      }
    });
  }

  getTemplateInfo = () => {
    // 填充第三方接口信息
    Modal.info({
      title: '模板详情',
      content: (
          <div>
            <p>some messages...some messages...</p>
            <p>some messages...some messages...</p>
          </div>
      ),
      onOk() {},
    });
  }

  renderTemplate = (info: TemplateInfo) => {
    // 需要展示多少信息都写在 Descriptions.Item 中
    return (
      <div className='template-item'>
          <Card title={"模板"+info.Tid}>
            <Descriptions
                bordered
                size="small">
              <Descriptions.Item label='tid'>{info.Tid}</Descriptions.Item>
              <Descriptions.Item label='instance_type'>{info.InstanceType}</Descriptions.Item>
              <Descriptions.Item label='image_id'>{info.ImageId}</Descriptions.Item>
              <Descriptions.Item label='resource_type'>{info.ResourceType}</Descriptions.Item>
              <Descriptions.Item label='bandwidth'>{info.Bandwidth}</Descriptions.Item>
              <Descriptions.Item label='last_time'>{info.LastTime}</Descriptions.Item>
            </Descriptions>
            <Button size="large" style={{float: "right", margin:"15px"}} onClick={this.getTemplateInfo}>查看更多关于此模板的信息</Button>
          </Card>

        </div>
    )
  }


  render() {
    const { xxx } = this.state;
    return (
      <>
        {xxx ?
        <div>
          {xxx.map(item => this.renderTemplate(item))}
        </div>
        :
         <Empty></Empty>}
      </>
    );
  }
}
