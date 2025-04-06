package handlers

import (
	"fmt"
	"net/http"
)

func LoginClientPage(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:8000/callback"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>ClipSync Login</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background: #f0f2f5;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
		}
		.container {
			background: white;
			padding: 2rem;
			border-radius: 12px;
			box-shadow: 0 0 10px rgba(0,0,0,0.1);
			width: 300px;
			text-align: center;
		}
		input {
			width: 100%%;
			padding: 10px;
			margin: 10px 0;
			border: 1px solid #ccc;
			border-radius: 8px;
		}
		button {
			width: 100%%;
			padding: 10px;
			background: #4CAF50;
			color: white;
			border: none;
			border-radius: 8px;
			cursor: pointer;
			font-size: 16px;
		}
		button:hover {
			background: #45a049;
		}
		.links {
			margin-top: 15px;
			font-size: 14px;
		}
		.links a {
			text-decoration: none;
			color: #4CAF50;
			margin: 0 5px;
		}
		#status {
			margin-top: 10px;
			color: #d00;
			font-size: 14px;
		}
	</style>
</head>
<body>
	<div class="container">
		<h2>Login to ClipSync</h2>
		<form id="loginForm">
			<input type="email" id="email" placeholder="Email" required>
			<input type="password" id="password" placeholder="Password" required>
			<button type="submit">Login</button>
		</form>
		<div class="links">
			<a href="/forgot-password-client">Forgot Password?</a> |
			<a href="/register-client">Sign up</a>
		</div>
		<div id="status"></div>
	</div>
	<script>
		const form = document.getElementById("loginForm");
		form.onsubmit = async (e) => {
			e.preventDefault();
			const email = document.getElementById("email").value;
			const password = document.getElementById("password").value;

			const res = await fetch("/login", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email, password })
			});

			const data = await res.json();
			if (res.ok && data.token) {
				window.location.href = "%s?token=" + data.token;
			} else {
				document.getElementById("status").textContent = data.message || "Login failed";
			}
		};
	</script>
</body>
</html>`, redirectURI)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
