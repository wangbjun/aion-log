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
                name: '总览',
                icon: 'smile',
                component: './Player',
              },
              {
                name: '玩家明细',
                icon: 'TeamOutlined',
                path: '/players',
                component: './PlayerDetail',
              },
              {
                name: '战斗日志',
                icon: 'table',
                path: '/log',
                component: './Log',
              },
              {
                name: '异常分析',
                icon: 'LockOutlined',
                path: '/rank',
                component: './Rank',
              },
              {
                name: '数据导入',
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
