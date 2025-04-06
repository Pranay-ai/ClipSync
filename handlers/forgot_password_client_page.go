package handlers

import (
	"fmt"
	"net/http"
)

func ForgotPasswordClientPage(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/reset-password-client"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Forgot Password - ClipSync</title>
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
			background: #2196F3;
			color: white;
			border: none;
			border-radius: 8px;
			cursor: pointer;
			font-size: 16px;
		}
		button:hover {
			background: #1976D2;
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
		<h2>Reset Your Password</h2>
		<form id="forgotForm">
			<input type="email" id="email" placeholder="Email address" required>
			<button type="submit">Send OTP</button>
		</form>
		<div id="status"></div>
	</div>
	<script>
		const form = document.getElementById("forgotForm");
		form.onsubmit = async (e) => {
			e.preventDefault();
			const email = document.getElementById("email").value;

			const res = await fetch("/forgot-password", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email })
			});

			const data = await res.json();
			if (res.ok) {
				window.location.href = "%s"; // Redirect after success (optional)
			} else {
				document.getElementById("status").textContent = data.message || "Failed to send OTP.";
			}
		};
	</script>
</body>
</html>`, redirectURI)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
