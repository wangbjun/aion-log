export default [
  {
    path: '/',
    component: '../layouts/BlankLayout',
    routes: [
      {
        path: '/user',
        component: '../layouts/UserLayout',
        routes: [
          {
            name: 'login',
            path: '/user/login',
            component: './User/login',
          },
        ],
      },
      {
        path: '/',
        routes: [
          {
            path: '/',
            component: '../layouts/BasicLayout',
            routes: [
              {
                path: '/',
                name: '卡多尔',
                icon: 'smile',
                component: './Player',
              },
              {
                name: '普埃塔',
                icon: 'table',
                path: '/log',
                component: './Log',
              },
              {
                name: '封神榜',
                icon: 'LockOutlined',
                path: '/rank',
                component: './Rank',
              },
              {
                name: '导入日志',
                icon: 'UploadOutlined',
                path: '/import',
                component: './ImportLog',
              },
              {
                component: './404',
              },
            ],
          },
          {
            component: './404',
          },
        ],
      },
    ],
  },
  {
    component: './404',
  },
];
