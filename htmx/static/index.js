const button = document.getElementById("button");
const message = document.getElementById("message");
button.onclick = function() {
  button.innerText = "Loading...";
  const body = JSON.stringify({ message: "Hello from HTMX!" });
  fetch("/api", { method: "POST", headers: { "Content-Type": "application/json" }, body })
    .then(response => response.json())
    .then(data => {
      message.innerText = data.message;
    });
};

