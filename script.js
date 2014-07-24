function post() {
  var http = new XMLHttpRequest(),
    url = "/send",
    params = document.getElementById('input').value;
  document.getElementById('input').value = '';
  document.getElementById('input').style.backgroundColor = '#cccccc';
  http.open("POST", url, true);
  http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  http.onreadystatechange = function() {
      if (http.readyState == 4 && http.status == 200) {
          document.getElementById('input').style.backgroundColor = '#ffffff';
      }
  }
  http.send(params);
}
