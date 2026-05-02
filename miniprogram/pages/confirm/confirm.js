var api = require('../../utils/api.js');
var app = getApp();

Page({
  data: {
    song: null,
    singers: [],
    selectedSingerIndex: -1,
    message: '',
    price: 10,
    submitting: false
  },

  onLoad: function (options) {
    if (options.song) {
      try {
        var song = JSON.parse(decodeURIComponent(options.song));
        this.setData({ song: song });
      } catch (e) {
        wx.showToast({ title: '歌曲信息错误', icon: 'none' });
      }
    }
    this.loadSingers();
  },

  requireLogin: function (callback) {
    if (!app.globalData.userInfo) {
      wx.navigateTo({ url: '/pages/login/login' });
      return false;
    }
    if (callback) callback();
    return true;
  },

  loadSingers: function () {
    var that = this;
    api.getSingers(1, 50).then(function (res) {
      that.setData({ singers: (res.data && res.data.list) || [] });
    }).catch(function () {});
  },

  onSingerChange: function (e) {
    this.setData({ selectedSingerIndex: e.detail.value });
  },

  onMessageInput: function (e) {
    this.setData({ message: e.detail.value });
  },

  onSubmit: function () {
    if (this.data.submitting) return;

    var that = this;
    // 提交订单需要登录
    if (!that.requireLogin(function () {
      that.doSubmit();
    })) return;
  },

  doSubmit: function () {
    var song = this.data.song;
    var idx = this.data.selectedSingerIndex;
    if (!song) {
      wx.showToast({ title: '歌曲信息缺失', icon: 'none' });
      return;
    }
    if (idx < 0 || idx >= this.data.singers.length) {
      wx.showToast({ title: '请选择歌手', icon: 'none' });
      return;
    }

    var singer = this.data.singers[idx];
    this.setData({ submitting: true });

    var orderData = {
      song_id: song.id,
      song_name: song.name,
      artist: song.artist,
      cover: song.cover || song.pic || '',
      source: song.source || 'netease',
      message: this.data.message,
      singer_id: singer.id,
      price: this.data.price
    };

    var that = this;
    api.createOrder(orderData).then(function (res) {
      that.setData({ submitting: false });
      wx.showToast({ title: '下单成功！', icon: 'success', duration: 1500 });
      setTimeout(function () {
        wx.redirectTo({
          url: '/pages/success/success?order=' + encodeURIComponent(JSON.stringify(res.data))
        });
      }, 1500);
    }).catch(function () {
      that.setData({ submitting: false });
    });
  }
});
