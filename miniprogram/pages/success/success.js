Page({
  data: {
    order: null
  },

  onLoad: function (options) {
    if (options.order) {
      try {
        var order = JSON.parse(decodeURIComponent(options.order));
        this.setData({ order: order });
      } catch (e) {
        wx.showToast({ title: '订单信息错误', icon: 'none' });
      }
    }
  },

  onBackHome: function () {
    wx.reLaunch({
      url: '/pages/index/index'
    });
  },

  onViewOrders: function () {
    wx.reLaunch({
      url: '/pages/index/index'
    });
  }
});
