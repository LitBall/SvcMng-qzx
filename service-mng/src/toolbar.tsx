import React from 'react';
import { Button } from 'antd';

interface ToolbarProps {
  createInstanceByTemplate: () => void;
  createInstance: () => void;
  manageKey: () => void;
}

export class Toolbar extends React.Component<ToolbarProps> {
  render() {
    const { createInstanceByTemplate, createInstance, manageKey } = this.props;
    return (
      <div className='toolbar-button-groups'>
        <Button type="text" size="large" onClick={createInstanceByTemplate}>使用模板创建实例</Button>
        <Button type="text" size="large" onClick={createInstance}>自定义配置创建实例</Button>
        <Button type="text" size="large" onClick={manageKey}>管理密钥</Button>
      </div>
    );
  }
}
