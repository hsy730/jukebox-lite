var api = require('../../utils/api.js');
var app = getApp();

Page({
  data: {
    keyword: '',
    categories: [],
    activeCategory: '',
    activeCategoryName: '推荐',
    songs: [],
    singers: [],
    page: 1,
    limit: 20,
    hasMore: true,
    loading: false,
    searchMode: false,
    source: 'netease',
    currentRole: 'user',
    singerOrders: [],
    singerOrdersPage: 1,
    singerOrdersTotal: 0,
    singerOrdersLoading: false,
    userInfo: null
  },

  onLoad: function () {
    var userInfo = app.globalData.userInfo;
    this.setData({
      userInfo: userInfo,
      currentRole: (userInfo && userInfo.role) || 'user'
    });
    this.loadByRole();
  },

  onShow: function () {
    var userInfo = app.globalData.userInfo;
    this.setData({
      userInfo: userInfo,
      currentRole: (userInfo && userInfo.role) || 'user'
    });
  },

  requireLogin: function (callback) {
    if (!app.globalData.userInfo) {
      wx.navigateTo({ url: '/pages/login/login' });
      return false;
    }
    if (callback) callback();
    return true;
  },

  onTapLogin: function () {
    if (!app.globalData.userInfo) {
      wx.navigateTo({ url: '/pages/login/login' });
    }
  },

  onBecomeSinger: function () {
    var that = this;
    if (that.data.currentRole === 'singer') return;
    if (app.globalData.userInfo) {
      api.switchRole('singer').then(function (res) {
        app.globalData.userInfo = res.data;
        that.setData({ userInfo: res.data });
      }).catch(function () {});
    }
    that.setData({
      currentRole: 'singer',
      page: 1,
      songs: [],
      hasMore: true,
      searchMode: false,
      activeCategory: '',
      activeCategoryName: '推荐',
      singerOrders: [],
      singerOrdersPage: 1
    });
    that.loadByRole();
  },

  switchToSinger: function () {
    this.onBecomeSinger();
  },

  loadByRole: function () {
    if (this.data.currentRole === 'singer') {
      this.loadSingerOrders();
    } else {
      this.loadCategories();
      this.loadSingers();
      this.loadRecommendSongs();
    }
  },

  onSwitchRole: function () {
    var that = this;
    var newRole = that.data.currentRole === 'user' ? 'singer' : 'user';

    if (app.globalData.userInfo) {
      api.switchRole(newRole).then(function (res) {
        app.globalData.userInfo = res.data;
        that.setData({ userInfo: res.data });
      }).catch(function () {});
    }

    that.setData({
      currentRole: newRole,
      page: 1,
      songs: [],
      hasMore: true,
      searchMode: false,
      activeCategory: '',
      activeCategoryName: '推荐',
      singerOrders: [],
      singerOrdersPage: 1
    });
    that.loadByRole();
  },

  onLogout: function () {
    var that = this;
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: function (res) {
        if (res.confirm) {
          api.removeToken();
          app.globalData.userInfo = null;
          that.setData({
            userInfo: null,
            currentRole: 'user',
            page: 1,
            songs: [],
            hasMore: true,
            searchMode: false,
            activeCategory: '',
            activeCategoryName: '推荐',
            singerOrders: [],
            singerOrdersPage: 1
          });
          that.loadByRole();
        }
      }
    });
  },

  loadSingerOrders: function () {
    var that = this;
    if (!app.globalData.userInfo) {
      that.setData({ singerOrders: [], singerOrdersTotal: 0, singerOrdersLoading: false });
      return;
    }
    that.setData({ singerOrdersLoading: true });
    api.getUserOrders(that.data.singerOrdersPage, 20).then(function (res) {
      var newList = (res.data && res.data.list) || [];
      that.setData({
        singerOrders: that.data.singerOrdersPage === 1 ? newList : that.data.singerOrders.concat(newList),
        singerOrdersTotal: (res.data && res.data.total) || 0,
        singerOrdersLoading: false
      });
    }).catch(function () {
      that.setData({ singerOrdersLoading: false });
    });
  },

  onPullDownRefresh: function () {
    if (this.data.currentRole === 'singer') {
      this.setData({ singerOrdersPage: 1, singerOrders: [] });
      this.loadSingerOrders();
    } else {
      this.setData({ page: 1, songs: [], hasMore: true });
      if (this.data.searchMode) {
        this.doSearch();
      } else if (this.data.activeCategory) {
        this.loadCategorySongs();
      } else {
        this.loadRecommendSongs();
      }
    }
    wx.stopPullDownRefresh();
  },

  onReachBottom: function () {
    if (this.data.currentRole === 'singer') {
      var loaded = this.data.singerOrders.length;
      if (loaded < this.data.singerOrdersTotal && !this.data.singerOrdersLoading) {
        this.setData({ singerOrdersPage: this.data.singerOrdersPage + 1 });
        this.loadSingerOrders();
      }
      return;
    }
    if (this.data.hasMore && !this.data.loading) {
      this.setData({ page: this.data.page + 1 });
      if (this.data.searchMode) {
        this.doSearch();
      } else if (this.data.activeCategory) {
        this.loadCategorySongs();
      } else {
        this.loadRecommendSongs();
      }
    }
  },

  loadCategories: function () {
    var that = this;
    api.getCategories(that.data.source).then(function (res) {
      that.setData({ categories: res.data || [] });
    }).catch(function () {});
  },

  loadSingers: function () {
    var that = this;
    api.getSingers(1, 50).then(function (res) {
      that.setData({ singers: (res.data && res.data.list) || [] });
    }).catch(function () {});
  },

  loadRecommendSongs: function () {
    var that = this;
    that.setData({ loading: true });
    api.getCategorySongs('', that.data.source, that.data.page, that.data.limit).then(function (res) {
      var newSongs = res.data || [];
      that.setData({
        songs: that.data.page === 1 ? newSongs : that.data.songs.concat(newSongs),
        hasMore: newSongs.length >= that.data.limit,
        loading: false
      });
    }).catch(function () {
      that.setData({ loading: false });
    });
  },

  onSearchInput: function (e) {
    this.setData({ keyword: e.detail.value });
  },

  onSearch: function () {
    if (!this.data.keyword.trim()) {
      wx.showToast({ title: '请输入搜索关键词', icon: 'none' });
      return;
    }
    this.setData({ page: 1, songs: [], hasMore: true, searchMode: true, activeCategory: '', activeCategoryName: '搜索结果' });
    this.doSearch();
  },

  doSearch: function () {
    var that = this;
    that.setData({ loading: true });
    api.searchSongs(that.data.keyword, that.data.source, that.data.page, that.data.limit).then(function (res) {
      var newSongs = res.data || [];
      that.setData({
        songs: that.data.page === 1 ? newSongs : that.data.songs.concat(newSongs),
        hasMore: newSongs.length >= that.data.limit,
        loading: false
      });
    }).catch(function () {
      that.setData({ loading: false });
    });
  },

  onClearSearch: function () {
    this.setData({ keyword: '', searchMode: false, page: 1, songs: [], hasMore: true, activeCategory: '', activeCategoryName: '推荐' });
    this.loadRecommendSongs();
  },

  onCategoryTap: function (e) {
    var id = e.currentTarget.dataset.id;
    var name = e.currentTarget.dataset.name;
    this.setData({
      activeCategory: id,
      activeCategoryName: name,
      searchMode: false,
      page: 1,
      songs: [],
      hasMore: true
    });
    this.loadCategorySongs();
  },

  loadCategorySongs: function () {
    var that = this;
    that.setData({ loading: true });
    api.getCategorySongs(that.data.activeCategory, that.data.source, that.data.page, that.data.limit).then(function (res) {
      var newSongs = res.data || [];
      that.setData({
        songs: that.data.page === 1 ? newSongs : that.data.songs.concat(newSongs),
        hasMore: newSongs.length >= that.data.limit,
        loading: false
      });
    }).catch(function () {
      that.setData({ loading: false });
    });
  },

  onSongTap: function (e) {
    var song = e.currentTarget.dataset.song;
    var that = this;
    // 点歌需要登录
    if (!that.requireLogin(function () {
      wx.navigateTo({
        url: '/pages/confirm/confirm?song=' + encodeURIComponent(JSON.stringify(song))
      });
    })) return;
  },

  onSourceChange: function (e) {
    var sources = ['netease', 'kuwo', 'qq'];
    var idx = e.detail.value;
    this.setData({
      source: sources[idx],
      page: 1,
      songs: [],
      hasMore: true,
      categories: [],
      activeCategory: '',
      activeCategoryName: '推荐'
    });
    this.loadCategories();
    this.loadRecommendSongs();
  },

  onOrderStatusTap: function (e) {
    var orderId = e.currentTarget.dataset.id;
    var status = e.currentTarget.dataset.status;
    var newStatus = '';
    if (status === 'pending') newStatus = 'singing';
    else if (status === 'singing') newStatus = 'done';
    if (!newStatus) return;

    var that = this;
    wx.showModal({
      title: '确认操作',
      content: '确定更改订单状态吗？',
      success: function (res) {
        if (res.confirm) {
          api.request({
            url: '/orders/' + orderId + '/status',
            method: 'PUT',
            data: { status: newStatus }
          }).then(function () {
            wx.showToast({ title: '状态已更新', icon: 'success' });
            that.setData({ singerOrdersPage: 1, singerOrders: [] });
            that.loadSingerOrders();
          }).catch(function () {});
        }
      }
    });
  }
});
