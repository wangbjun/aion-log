/**
 * Ant Design Pro v4 use `@ant-design/pro-layout` to handle Layout.
 *
 * @see You can view component api by: https://github.com/ant-design/ant-design-pro-layout
 */
import ProLayout, {DefaultFooter} from '@ant-design/pro-layout';
import React, {useEffect, useMemo, useRef} from 'react';
import {connect, history, Link} from 'umi';
import {Button, Modal, Result} from 'antd';
import Authorized from '@/utils/Authorized';
import RightContent from '@/components/GlobalHeader/RightContent';
import {getMatchMenu} from '@umijs/route-utils';
import logo from '../assets/logo.ico';

const noMatch = (
  <Result
    status={403}
    title="403"
    subTitle="Sorry, you are not authorized to access this page."
    extra={
      <Button type="primary">
        <Link to="/user/login">Go Login</Link>
      </Button>
    }
  />
);

/** Use Authorized check all menu item */
const menuDataRender = (menuList) =>
  menuList.map((item) => {
    const localItem = {
      ...item,
      children: item.children ? menuDataRender(item.children) : undefined,
    };
    return Authorized.check(item.authority, localItem, null);
  });

const defaultFooterDom = (
  <DefaultFooter
    copyright={`${new Date().getFullYear()} 永恒之塔`}
    links={[]}
  />
);

const BasicLayout = (props) => {
  let {
    dispatch,
    children,
    settings,
    location = {
      pathname: '/',
    },
    visible
  } = props;
  visible = sessionStorage.getItem("modalClose") !== "true"
  const menuDataRef = useRef([]);
  // useEffect(() => {
  //   if (dispatch) {
  //     dispatch({
  //       type: 'user/fetchCurrent',
  //     });
  //   }
  // }, []);
  /** Init variables */
  const handleMenuCollapse = (payload) => {
    if (dispatch) {
      dispatch({
        type: 'global/changeLayoutCollapsed',
        payload,
      });
    }
  }; // get children authority

  const authorized = useMemo(
    () =>
      getMatchMenu(location.pathname || '/', menuDataRef.current).pop() || {
        authority: undefined,
      },
    [location.pathname],
  );
  const cancel = () => {
    if (dispatch) {
      dispatch({
        type: 'global/closeModal',
      });
    }
  }
  return (
    <ProLayout
      logo={logo}
      {...props}
      {...settings}
      onCollapse={handleMenuCollapse}
      onMenuHeaderClick={() => history.push('/')}
      menuItemRender={(menuItemProps, defaultDom) => {
        if (
          menuItemProps.isUrl ||
          !menuItemProps.path ||
          location.pathname === menuItemProps.path
        ) {
          return defaultDom;
        }

        return <Link to={menuItemProps.path}>{defaultDom}</Link>;
      }}
      itemRender={(route, params, routes, paths) => {
        const first = routes.indexOf(route) === 0;
        return first ? (
          <Link to={paths.join('/')}>{route.breadcrumbName}</Link>
        ) : (
          <span>{route.breadcrumbName}</span>
        );
      }}
      footerRender={() => {
        if (settings.footerRender || settings.footerRender === undefined) {
          return defaultFooterDom;
        }
        return null;
      }}
      menuDataRender={menuDataRender}
      rightContentRender={() => <RightContent/>}
      postMenuData={(menuData) => {
        menuDataRef.current = menuData || [];
        return menuData || [];
      }}
    >
      <Authorized authority={authorized.authority} noMatch={noMatch}>
        {children}
      </Authorized>
      <Modal visible={visible} onCancel={cancel} footer={null} title="友情提示" width="60%">
        <div style={{fontSize: 17, color: "red"}}>
          <p>一、本站纯属个人兴趣爱好而建立！</p>
          <p>二、本站数据来源于本人游戏客户端，其原理类似咱们用的DPS水表工具，从Chatlog日志分析所得，完全绿色，无侵入游戏！</p>
          <p>三、Chalog日志只能记录本人角色100米范围内的战斗行为，所以本日志仅供参考，无法保证数据完整性！</p>
          <p>四、大神榜单模块，分为3个段位，星耀段位是1秒3技能，最强王者是1秒4技能，荣耀王者是1秒5技能，都包含平砍技能！</p>
          <p>五、理性分析，上榜玩家不代表就一定非绿色，在延迟和FPS极致的情况，绿色玩家也可以做到1秒2技能+平砍，或者纯3技能！</p>
          <p>六、根据老G反馈，在要塞战期间，由于延迟原因，客户端记录日志可能存在偏差，不一定代表实际行为，此日志也不可以作为封禁证据！</p>
          <p>七、条件有限，没有其他区的号，若有人愿意提供其它时间段数据，可以联系本人！</p>
        </div>
      </Modal>
    </ProLayout>
  );
};

export default connect(({global, settings}) => ({
  collapsed: global.collapsed,
  visible: global.visible,
  settings,
}))(BasicLayout);
