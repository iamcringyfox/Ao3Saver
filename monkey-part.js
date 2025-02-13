// ==UserScript==
// @name         AO3 Link Sender
// @namespace    http://tampermonkey.net/
// @version      1.0
// @description  Отправляет ссылку страницы AO3 на сервер, если это works или series
// @author       Ваше имя
// @match        https://archiveofourown.org/works/*
// @match        https://archiveofourown.org/series/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    // URL сервера для отправки данных
    const SERVER_URL = 'http://localhost:8080'; // Замените на ваш сервер

    // Проверяем, что мы находимся на странице works или series
    if (window.location.href.includes('/works/') || window.location.href.includes('/series/')) {
        const currentUrl = window.location.href;

        // Отправляем POST-запрос на сервер
        fetch(SERVER_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: currentUrl }),
        })
        .then(response => {
            if (!response.ok) {
                console.error('Ошибка при отправке данных:', response.status);
            } else {
                console.log('Ссылка успешно отправлена на сервер');
            }
        })
        .catch(error => {
            console.error('Ошибка при отправке POST-запроса:', error);
        });
    }
})();
