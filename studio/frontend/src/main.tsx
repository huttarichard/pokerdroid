import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";

window.addEventListener("keypress", (event) => {
  if (event.metaKey && event.key === "c") {
    document.execCommand("copy");
    event.preventDefault();
  }

  if (event.metaKey && event.key === "v") {
    document.execCommand("paste");
    event.preventDefault();
  }

  if (event.metaKey && event.key === "x") {
    document.execCommand("cut");
    event.preventDefault();
  }

  if (event.metaKey && event.key === "z") {
    document.execCommand("undo");
    event.preventDefault();
  }

  if (event.metaKey && event.key === "y") {
    document.execCommand("redo");
    event.preventDefault();
  }

  if (event.metaKey && event.key === "a") {
    document.execCommand("selectAll");
    event.preventDefault();
  }

  if (event.metaKey && event.key === "r") {
    window.location.reload();
    event.preventDefault();
  }
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
