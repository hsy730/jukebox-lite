var app = getApp();

function request(options) {
  return new Promise(function (resolve, reject) {
    wx.request({
      url: app.globalData.apiBase + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: Object.assign({
        'Content-Type': 'application/json'
      }, options.header || {}),
      success: function (res) {
        if (res.data.code === 0) {
          resolve(res.data);
        } else {
          wx.showToast({
            title: res.data.message || '请求失败',
            icon: 'none'
          });
          reject(res.data);
        }
      },
      fail: function (err) {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        });
        reject(err);
      }
    });
  });
}

function searchSongs(keyword, source, page, limit) {
  return request({
    url: '/songs/search',
    data: {
      keyword: keyword,
      source: source || '',
      page: page || 1,
      limit: limit || 20
    }
  });
}

function getCategories(source) {
  return request({
    url: '/songs/categories',
    data: {
      source: source || ''
    }
  });
}

function getCategorySongs(toplistId, source, page, limit) {
  return request({
    url: '/songs/category',
    data: {
      toplist_id: toplistId,
      source: source || '',
      page: page || 1,
      limit: limit || 20
    }
  });
}

function getSongInfo(source, id) {
  return request({
    url: '/songs/info',
    data: {
      source: source,
      id: id
    }
  });
}

function createOrder(data) {
  return request({
    url: '/orders',
    method: 'POST',
    data: data
  });
}

function getOrder(id) {
  return request({
    url: '/orders/' + id
  });
}

function getSingerOrders(singerId, page, limit) {
  return request({
    url: '/orders/singer',
    data: {
      singer_id: singerId,
      page: page || 1,
      limit: limit || 20
    }
  });
}

function getUserOrders(page, limit) {
  return request({
    url: '/orders/user',
    data: {
      page: page || 1,
      limit: limit || 20
    }
  });
}

function getSingers(page, limit) {
  return request({
    url: '/singers',
    data: {
      page: page || 1,
      limit: limit || 20
    }
  });
}

module.exports = {
  request: request,
  searchSongs: searchSongs,
  getCategories: getCategories,
  getCategorySongs: getCategorySongs,
  getSongInfo: getSongInfo,
  createOrder: createOrder,
  getOrder: getOrder,
  getSingerOrders: getSingerOrders,
  getUserOrders: getUserOrders,
  getSingers: getSingers
};
