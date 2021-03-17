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
        component: '../layouts/SecurityLayout',
        routes: [
          {
            path: '/',
            component: '../layouts/BasicLayout',
            authority: ['admin', 'user', 'guest'],
            routes: [
              {
                path: '/',
                name: '极乐世界',
                icon: 'smile',
                component: './Player',
              },
              {
                name: '日志列表',
                icon: 'table',
                path: '/log',
                component: './Log',
              },
              {
                name: '上传日志',
                icon: 'UploadOutlined',
                path: '/parse',
                component: './Parse',
                authority: ['admin'],
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
