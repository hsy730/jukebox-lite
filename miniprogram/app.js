App({
  globalData: {
    apiBase: 'http://localhost:8080/api',
    userInfo: null
  },
  onLaunch: function () {
    var that = this;
    wx.login({
      success: function (res) {
        console.log('login code:', res.code);
      }
    });
  }
});
