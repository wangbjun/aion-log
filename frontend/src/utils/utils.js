import {parse} from 'querystring';
/* eslint no-useless-escape:0 import/prefer-default-export:0 */

const reg = /(((^https?:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(?::\d+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)$/;
export const isUrl = (path) => reg.test(path);
export const isAntDesignPro = () => {
  if (ANT_DESIGN_PRO_ONLY_DO_NOT_USE_IN_YOUR_PRODUCTION === 'site') {
    return true;
  }

  return window.location.hostname === 'preview.pro.ant.design';
}; // 给官方演示站点用，用于关闭真实开发环境不需要使用的特性

export const isAntDesignProOrDev = () => {
  const {NODE_ENV} = process.env;

  if (NODE_ENV === 'development') {
    return true;
  }

  return isAntDesignPro();
};
export const getPageQuery = () => parse(window.location.href.split('?')[1]);

export const playerPros = [
  {name: "其它", logo: "unknown.png", class: 0},
  {name: "剑星", logo: "jx.png", class: 1},
  {name: "守护", logo: "sh.png", class: 2},
  {name: "杀星", logo: "sx.png", class: 3},
  {name: "弓星", logo: "gx.png", class: 4},
  {name: "治愈", logo: "zy.png", class: 5},
  {name: "护法", logo: "hf.png", class: 6},
  {name: "精灵", logo: "jl.png", class: 7},
  {name: "魔道", logo: "md.png", class: 8},
  {name: "执行", logo: "zxz.png", class: 9},
]

export const getTypeColor = (type) => {
  if (type === 1) {
    return ["green","天族"]
  } else if (type === 2) {
    return ["blue","魔族"]
  } else {
    return ["orange","NPC"]
  }
}
