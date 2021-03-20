import {history} from 'umi';
import {setAuthority} from '@/utils/authority';
import {getPageQuery} from '@/utils/utils';
import {message} from 'antd';
import {login} from "@/services/user";

const Model = {
  namespace: 'login',
  state: {status: undefined},
  effects: {
    * login({payload}, {call, put}) {
      const response = yield call(login, payload);
      yield put({
        type: 'changeLoginStatus',
        payload: response,
      }); // Login successfully
      if (response.code === 200) {
        const urlParams = new URL(window.location.href);
        const params = getPageQuery();
        message.success('🎉 🎉 🎉  登录成功！');
        let {redirect} = params;
        if (redirect) {
          const redirectUrlParams = new URL(redirect);
          if (redirectUrlParams.origin === urlParams.origin) {
            redirect = redirect.substr(urlParams.origin.length);
            if (redirect.match(/^\/.*#/)) {
              redirect = redirect.substr(redirect.indexOf('#') + 1);
            }
          } else {
            window.location.href = '/';
            return;
          }
        }
        history.push(redirect || '/');
      }
    },
    * logout() {
      localStorage.clear()
      window.location.reload()
    },
  },
  reducers: {
    changeLoginStatus(state, {payload}) {
      let data = payload && payload.data
      setAuthority(data && data.currentAuthority);
      localStorage.setItem("token", data && data.token)
      return {...state, status: payload.msg};
    },
  },
};
export default Model;
