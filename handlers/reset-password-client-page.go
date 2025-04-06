package handlers

import "net/http"

func ResetPasswordClientPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Reset Password - ClipSync</title>
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
			width: 320px;
			text-align: center;
		}
		input {
			width: 100%;
			padding: 10px;
			margin: 10px 0;
			border: 1px solid #ccc;
			border-radius: 8px;
		}
		button {
			width: 100%;
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
		<form id="resetForm">
			<input type="email" id="email" placeholder="Email" required>
			<input type="text" id="otp" placeholder="OTP Code" required>
			<input type="password" id="newPassword" placeholder="New Password" required>
			<input type="password" id="confirmNewPassword" placeholder="Confirm New Password" required>
			<button type="submit">Reset Password</button>
		</form>
		<div id="status"></div>
	</div>
	<script>
		const form = document.getElementById("resetForm");
		form.onsubmit = async (e) => {
			e.preventDefault();

			const email = document.getElementById("email").value;
			const otp = document.getElementById("otp").value;
			const newPassword = document.getElementById("newPassword").value;
			const confirmNewPassword = document.getElementById("confirmNewPassword").value;

			if (newPassword !== confirmNewPassword) {
				document.getElementById("status").textContent = "Passwords do not match.";
				return;
			}

			const res = await fetch("/reset-password", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email, otp, new_password: newPassword })
			});

			const data = await res.json();
			if (res.ok) {
				window.location.href = "/login-client";
			} else {
				document.getElementById("status").textContent = data.message || "Reset failed.";
			}
		};
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
