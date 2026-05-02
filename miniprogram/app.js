var api = require('./utils/api.js');

App({
  globalData: {
    apiBase: 'http://localhost:8080/api',
    userInfo: null
  },

  onLaunch: function () {
    var token = api.getToken();
    if (token) {
      this.checkLogin();
    }
    // 未登录也允许进入首页浏览，不强制跳转登录页
  },

  checkLogin: function () {
    var that = this;
    api.getProfile().then(function (res) {
      that.globalData.userInfo = res.data;
    }).catch(function () {
      api.removeToken();
      that.globalData.userInfo = null;
    });
  },

  redirectToLogin: function () {
    wx.navigateTo({
      url: '/pages/login/login'
    });
  }
});
