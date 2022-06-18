import axios from 'axios';

export function get(url: string, params: any) {
  return new Promise((resolve, reject) => {
    axios.get(url, params).then(response => {
      if (response.status === 200) {
        resolve(response.data);
      } else {
        reject(response.data.msg);
      }
    }).catch(e => {
      reject(e);
    })
  })
}

export function post(url: string, postData: any) {
  return new Promise((resolve, reject) => {
    axios.post(url, postData).then(response => {
      if (response.status === 200) {
        resolve(response.data);
      } else {
        reject(response.data.msg);
      }
    }).catch(e => {
      reject(e);
    })
  })
}
