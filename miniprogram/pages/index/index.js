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
    source: 'netease'
  },

  onLoad: function () {
    this.loadCategories();
    this.loadSingers();
    this.loadRecommendSongs();
  },

  onPullDownRefresh: function () {
    this.setData({ page: 1, songs: [], hasMore: true });
    if (this.data.searchMode) {
      this.doSearch();
    } else if (this.data.activeCategory) {
      this.loadCategorySongs();
    } else {
      this.loadRecommendSongs();
    }
    wx.stopPullDownRefresh();
  },

  onReachBottom: function () {
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
    wx.navigateTo({
      url: '/pages/confirm/confirm?song=' + encodeURIComponent(JSON.stringify(song))
    });
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
  }
});
