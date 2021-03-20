import {DefaultFooter, getMenuData, getPageTitle} from '@ant-design/pro-layout';
import {Helmet, HelmetProvider} from 'react-helmet-async';
import {connect, Link, useIntl} from 'umi';
import React from 'react';
import logo from '../assets/logo.ico';
import styles from './UserLayout.less';

const UserLayout = (props) => {
  const {
    route = {
      routes: [],
    },
  } = props;
  const {routes = []} = route;
  const {
    children,
    location = {
      pathname: '',
    },
  } = props;
  const {formatMessage} = useIntl();
  const {breadcrumb} = getMenuData(routes);
  const title = getPageTitle({
    pathname: location.pathname,
    formatMessage,
    breadcrumb,
    ...props,
  });
  return (
    <HelmetProvider>
      <Helmet>
        <title>{title}</title>
        <meta name="description" content={title}/>
      </Helmet>

      <div className={styles.container}>
        <div className={styles.content}>
          <div className={styles.top}>
            <div className={styles.header}>
              <Link to="/">
                <img alt="logo" className={styles.logo} src={logo}/>
                <span className={styles.title}>AION</span>
              </Link>
            </div>
            <div className={styles.desc}>
              永恒之塔
            </div>
          </div>
          {children}
        </div>
        <DefaultFooter copyright={`${new Date().getFullYear()} 永恒之塔`} links={[]}/>
      </div>
    </HelmetProvider>
  );
};

export default connect(({settings}) => ({...settings}))(UserLayout);
