document.getElementById("calcForm").addEventListener("submit", function(e) {
  e.preventDefault();
  var expr = document.getElementById("expression").value;
  document.getElementById("result").innerHTML = "";
  document.getElementById("status").innerHTML = "<div class='d-flex align-items-center'><strong>Отправка запроса...</strong><div class='spinner-border ms-2 text-primary' role='status'><span class='visually-hidden'>Загрузка...</span></div></div>";
  fetch("/api/v1/calculate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ expression: expr })
  })
  .then(res => res.json())
  .then(data => {
    document.getElementById("result").innerHTML = "<div class='alert alert-success'>Задача принята. ID: " + data.id + "</div>";
    pollExpressionStatus(data.id);
  })
  .catch(err => {
    console.error("Ошибка отправки запроса", err);
    document.getElementById("result").innerHTML = "<div class='alert alert-danger'>Ошибка отправки запроса</div>";
    document.getElementById("status").innerHTML = "";
  });
});

function pollExpressionStatus(id) {
  var pollInterval = setInterval(function() {
    fetch("/api/v1/expressions/" + id)
      .then(res => res.json())
      .then(data => {
        if(data.status === "Выполнено.") {
          document.getElementById("status").innerHTML = "<div class='alert alert-success'>Результат: " + data.result + "</div>";
          clearInterval(pollInterval);
        } else {
          document.getElementById("status").innerHTML = "<div class='alert alert-info'>Статус: " + data.status + "</div>";
          if (data.status === "Ошибка.") {
            clearInterval(pollInterval);
          }
        }
      })
      .catch(err => {
        console.error("Ошибка получения статуса", err);
        document.getElementById("status").innerHTML = "<div class='alert alert-danger'>Ошибка получения статуса</div>";
        clearInterval(pollInterval);
      });
  }, 2000);
}