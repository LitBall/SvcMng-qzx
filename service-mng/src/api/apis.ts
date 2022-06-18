interface UrlRoute {
  [key: string]: string;
}

export const API: UrlRoute = {
  Host: 'http://127.0.0.1:8888',
  TemplateList: 'http://127.0.0.1:8888/ServiceManagement/getTemplateList',
  AccessKeyInfo: 'http://127.0.0.1:8888/ServiceManagement/getKey',
  AccessKeyUpdate: 'http://127.0.0.1:8888/ServiceManagement/updateKey',
  InstanceCreate: 'http://127.0.0.1:8888/ServiceManagement/addCloudInstance'
};
