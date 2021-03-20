import {LockOutlined, UserOutlined,} from '@ant-design/icons';
import {Alert} from 'antd';
import React from 'react';
import ProForm, {ProFormCheckbox, ProFormText} from '@ant-design/pro-form';
import {connect} from 'umi';
import styles from './index.less';

const LoginMessage = ({content}) => (
  <Alert
    style={{
      marginBottom: 24,
    }}
    message={content}
    type="error"
    showIcon
  />
);

const Login = (props) => {
  const {userLogin = {}, submitting} = props;
  const {status} = userLogin;
  const handleSubmit = (values) => {
    const {dispatch} = props;
    dispatch({
      type: 'login/login',
      payload: {...values},
    });
  };

  return (
    <div className={styles.main}>
      <ProForm
        initialValues={{
          autoLogin: true,
        }}
        submitter={{
          render: (_, dom) => dom.pop(),
          submitButtonProps: {
            loading: submitting,
            size: 'large',
            style: {
              width: '100%',
            },
          },
        }}
        onFinish={(values) => {
          handleSubmit(values);
          return Promise.resolve();
        }}
      >
        {status !== undefined && !submitting && (
          <LoginMessage content={status}/>
        )}
        <>
          <ProFormText
            name="email"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined className={styles.prefixIcon}/>,
            }}
            placeholder="用户名"
            rules={[
              {
                required: true,
                message: "请输入用户名!",
              },
            ]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className={styles.prefixIcon}/>,
            }}
            placeholder="密码"
            rules={[
              {
                required: true,
                message: "请输入密码！",
              },
            ]}
          />
        </>
        <div
          style={{
            marginBottom: 24,
          }}
        >
          <ProFormCheckbox noStyle name="autoLogin">
            自动登录
          </ProFormCheckbox>
          <a
            style={{
              float: 'right',
            }}
          >
          </a>
        </div>
      </ProForm>
    </div>
  );
};

export default connect(({login, loading}) => ({
  userLogin: login,
  submitting: loading.effects['login/login'],
}))(Login);
