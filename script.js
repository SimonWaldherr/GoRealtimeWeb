/*jslint browser: true, plusplus: true, indent: 2 */

var post = function () {
  "use strict";
  var http = new XMLHttpRequest(),
    url = "/send",
    params = document.getElementById('input').value;
  document.getElementById('input').value = '';
  document.getElementById('input').style.backgroundColor = '#cccccc';
  http.open("POST", url, true);
  http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  http.onreadystatechange = function () {
    if (http.readyState === 4 && http.status === 200) {
      document.getElementById('input').style.backgroundColor = '#ffffff';
    }
  };
  http.send(params);
};

var htmlEscape = function (str) {
  "use strict";
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
};
