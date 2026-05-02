var api = require('../../utils/api.js');
var app = getApp();

Page({
  data: {
    loading: false,
    canIUseGetUserProfile: wx.getUserProfile ? true : false,
    redirect: ''
  },

  onLoad: function (options) {
    if (options.redirect) {
      this.setData({ redirect: decodeURIComponent(options.redirect) });
    }
    var token = api.getToken();
    if (token) {
      this.checkLogin();
    }
  },

  checkLogin: function () {
    var that = this;
    api.getProfile().then(function (res) {
      app.globalData.userInfo = res.data;
      that.navigateAfterLogin();
    }).catch(function () {
      api.removeToken();
      app.globalData.userInfo = null;
    });
  },

  navigateAfterLogin: function () {
    var redirect = this.data.redirect;
    if (redirect) {
      wx.redirectTo({ url: redirect, fail: function () { wx.reLaunch({ url: '/pages/index/index' }); } });
    } else {
      var pages = getCurrentPages();
      if (pages.length > 1) {
        wx.navigateBack({ delta: 1 });
      } else {
        wx.reLaunch({ url: '/pages/index/index' });
      }
    }
  },

  onGetUserProfile: function () {
    var that = this;
    if (that.data.loading) return;
    that.setData({ loading: true });

    wx.getUserProfile({
      desc: '用于完善用户资料',
      success: function (userInfoRes) {
        var nickName = userInfoRes.userInfo.nickName;
        var avatar = userInfoRes.userInfo.avatarUrl;
        that.doLogin(nickName, avatar);
      },
      fail: function () {
        that.doLogin('微信用户', '');
      }
    });
  },

  onQuickLogin: function () {
    if (this.data.loading) return;
    this.doLogin('微信用户', '');
  },

  doLogin: function (nickName, avatar) {
    var that = this;
    wx.login({
      success: function (loginRes) {
        if (loginRes.code) {
          api.login(loginRes.code, nickName, avatar).then(function (res) {
            api.setToken(res.data.token);
            app.globalData.userInfo = res.data.user;
            that.navigateAfterLogin();
          }).catch(function () {
            that.setData({ loading: false });
          });
        } else {
          wx.showToast({ title: '微信登录失败', icon: 'none' });
          that.setData({ loading: false });
        }
      },
      fail: function () {
        wx.showToast({ title: '微信登录失败', icon: 'none' });
        that.setData({ loading: false });
      }
    });
  }
});
