import {LogoutOutlined, SettingOutlined, UserOutlined} from '@ant-design/icons';
import {Avatar, Menu, Spin} from 'antd';
import React from 'react';
import {connect, history} from 'umi';
import HeaderDropdown from '../HeaderDropdown';
import styles from './index.less';

class AvatarDropdown extends React.Component {
  onMenuClick = (event) => {
    const {key} = event;

    if (key === 'logout') {
      const {dispatch} = this.props;
      if (dispatch) {
        dispatch({
          type: 'login/logout',
        });
      }
    }
  };

  render() {
    const {
      currentUser = {
        avatar: '',
        name: '',
      },
      menu,
    } = this.props;
    const token = localStorage.getItem("token")
    const menuHeaderDropdown = (
      <Menu className={styles.menu} selectedKeys={[]} onClick={this.onMenuClick}>
        {menu && (
          <Menu.Item key="center">
            <UserOutlined/>
            个人中心
          </Menu.Item>
        )}
        {menu && (
          <Menu.Item key="settings">
            <SettingOutlined/>
            个人设置
          </Menu.Item>
        )}
        {menu && <Menu.Divider/>}
        {token && <Menu.Item key="logout">
          <LogoutOutlined/>
          退出登录
        </Menu.Item>}
      </Menu>
    );
    return currentUser && currentUser.name ? (
      <HeaderDropdown overlay={menuHeaderDropdown}>
        <span className={`${styles.action} ${styles.account}`}>
          <Avatar size="small" className={styles.avatar} src='./icons/icon.ico' alt="avatar" onClick={() => {
            history.push('/user/login')
          }}/>
          <span className={`${styles.name} anticon`} onClick={() => {
            history.push('/user/login')
          }}>{currentUser.name}</span>
        </span>
      </HeaderDropdown>
    ) : (
      <span className={`${styles.action} ${styles.account}`}>
        <Spin
          size="small"
          style={{
            marginLeft: 8,
            marginRight: 8,
          }}
        />
      </span>
    );
  }
}

export default connect(({user}) => ({
  currentUser: user.currentUser,
}))(AvatarDropdown);
