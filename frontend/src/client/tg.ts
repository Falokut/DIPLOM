export function checkInitData() {
    const webViewStatus = document.querySelector('#webview_data_status');
    if (Telegram.WebApp.initDataUnsafe.query_id &&
        Telegram.WebApp.initData &&
        webViewStatus.classList.contains('status_need')
    ) {
        webViewStatus.classList.remove('status_need');
        Telegram.WebApp.apiRequest('checkInitData', {}, function (result) {
            if (result.ok) {
                webViewStatus.textContent = 'Hash is correct (async)';
                webViewStatus.className = 'ok';
            } else {
                webViewStatus.textContent = result.error + ' (async)';
                webViewStatus.className = 'err';
            }
        });
    }
},