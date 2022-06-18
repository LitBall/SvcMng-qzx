import React, { useState } from 'react';
import { Layout } from 'antd';
import { Toolbar } from './toolbar';
import { CreateInstanceDrawer } from './create_instance_drawer'; // 使用自定义配置创建实例
import { TemplateDrawer } from './create_by_template_drawer';// 使用模板创建实例
import { ProfileDrawer } from './change_profile_drawer';// 管理密钥
import { Template } from './template';// 模板列表
import './main.css';



interface MainContentProps { };

export function MainContent(props: MainContentProps) {

  const { Header, Content } = Layout;

  const [visible1, setVisible1] = useState(false); // 使用模板创建实例
  const [visible2, setVisible2] = useState(false); // 使用自定义配置创建实例
  const [visible3, setVisible3] = useState(false); // 管理密钥

  const createInstanceByTemplate = () => {// 使用模板创建实例
    setVisible1(true);
  }

  const createInstance = () => {// 使用自定义配置创建实例
    setVisible2(true);
  }

  const manageKey = () => {// 管理密钥
    setVisible3(true);
  }

  const closeDrawer1 = () => {// 使用模板创建实例
    setVisible1(false);
  }

  const closeDrawer2 = () => {// 使用自定义配置创建实例
    setVisible2(false);
  }

  const closeDrawer3 = () => {// 管理密钥
    setVisible3(false);
  }

  return (
      <Layout className="layout">
        <Header className="header">
          <div className='name'>多云管理平台</div>
        </Header>
        <Layout>
          <Layout style={{ position: 'relative'}}>
            <Content
              className="content-background"
            >
              <Toolbar
                createInstanceByTemplate={createInstanceByTemplate}
                createInstance={createInstance}
                manageKey={manageKey}
              />
              <Template></Template>
              <TemplateDrawer visible={visible1} closeDrawer={closeDrawer1}></TemplateDrawer>
              <CreateInstanceDrawer visible={visible2} closeDrawer={closeDrawer2}></CreateInstanceDrawer>
              <ProfileDrawer visible={visible3} closeDrawer={closeDrawer3}></ProfileDrawer>
            </Content>
          </Layout>
        </Layout>
      </Layout>

  );
}
