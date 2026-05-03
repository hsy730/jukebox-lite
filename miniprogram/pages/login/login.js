var api = require('../../utils/api.js');
var app = getApp();

Page({
  data: {
    loading: false,
    avatarUrl: '',
    nickName: '',
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
    }).catch(function (err) {
      console.error('[Login] checkLogin failed:', err);
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

  onChooseAvatar: function (e) {
    var avatarUrl = e.detail.avatarUrl;
    console.log('[Login] avatar chosen:', avatarUrl);
    this.setData({ avatarUrl: avatarUrl });
  },

  onNicknameInput: function (e) {
    this.setData({ nickName: e.detail.value });
  },

  onNicknameBlur: function (e) {
    if (e.detail.value) {
      this.setData({ nickName: e.detail.value });
    }
  },

  onConfirmLogin: function () {
    var nickName = this.data.nickName || '微信用户';
    var avatar = this.data.avatarUrl || '';
    this.doLogin(nickName, avatar);
  },

  onQuickLogin: function () {
    if (this.data.loading) return;
    this.doLogin('微信用户', '');
  },

  doLogin: function (nickName, avatar) {
    var that = this;
    if (that.data.loading) return;
    that.setData({ loading: true });

    wx.login({
      success: function (loginRes) {
        if (loginRes.code) {
          console.log('[Login] wx.login success, code:', loginRes.code);
          api.login(loginRes.code, nickName, avatar).then(function (res) {
            console.log('[Login] login API success');
            api.setToken(res.data.token);
            app.globalData.userInfo = res.data.user;
            that.navigateAfterLogin();
          }).catch(function (err) {
            console.error('[Login] login API failed:', err);
            wx.showToast({ title: '登录失败，请重试', icon: 'none' });
            that.setData({ loading: false });
          });
        } else {
          console.error('[Login] wx.login failed, no code:', loginRes);
          wx.showToast({ title: '微信登录失败', icon: 'none' });
          that.setData({ loading: false });
        }
      },
      fail: function (err) {
        console.error('[Login] wx.login error:', err);
        wx.showToast({ title: '微信登录失败', icon: 'none' });
        that.setData({ loading: false });
      }
    });
  }
});
