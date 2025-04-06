package handlers

import "net/http"

func RegisterClientPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Register - ClipSync</title>
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
		#status {
			margin-top: 10px;
			color: #d00;
			font-size: 14px;
		}
	</style>
</head>
<body>
	<div class="container">
		<h2>Register for ClipSync</h2>
		<form id="registerForm">
			<input type="text" id="name" placeholder="Full Name" required>
			<input type="email" id="email" placeholder="Email" required>
			<input type="password" id="password" placeholder="Password" required>
			<input type="password" id="confirmPassword" placeholder="Confirm Password" required>
			<button type="submit">Register</button>
		</form>
		<div id="status"></div>
	</div>
	<script>
		const form = document.getElementById("registerForm");
		form.onsubmit = async (e) => {
			e.preventDefault();

			const name = document.getElementById("name").value;
			const email = document.getElementById("email").value;
			const password = document.getElementById("password").value;
			const confirmPassword = document.getElementById("confirmPassword").value;

			if (password !== confirmPassword) {
				document.getElementById("status").textContent = "Passwords do not match.";
				return;
			}

			const res = await fetch("/register", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ name, email, password })
			});

			const data = await res.json();
			if (res.ok && data.token) {
				const params = new URLSearchParams(window.location.search);
				const redirect = params.get("redirect_uri") || "";
				window.location.href = "/login-client" + (redirect ? "?redirect_uri=" + encodeURIComponent(redirect) : "");
			} else {
				document.getElementById("status").textContent = data.message || "Registration failed.";
			}
		};
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
