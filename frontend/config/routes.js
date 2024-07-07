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
                name: '极乐世界',
                icon: 'smile',
                component: './Player',
              },
              {
                name: '诸神黄昏',
                icon: 'LockOutlined',
                path: '/rank',
                component: './Rank',
              },
              {
                name: '日志列表',
                icon: 'table',
                path: '/log',
                component: './Log',
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
