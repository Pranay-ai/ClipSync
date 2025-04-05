package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"clipsync.com/m/db"
	"clipsync.com/m/models"
	"clipsync.com/m/utils"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordRequest

	if err := json.NewDecoder((r.Body)).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(newHashedPassword)
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})

}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	otp := fmt.Sprintf("%06d", rand.Intn(100000)) // 6-digit OTP

	// Store OTP in Redis for 5 minutes

	ctx := context.Background()
	redisKey := fmt.Sprintf("reset_otp:%s", req.Email)
	if err := db.RedisClient.Set(ctx, redisKey, otp, 5*time.Minute).Err(); err != nil {
		http.Error(w, "Failed to store OTP", http.StatusInternalServerError)
		return
	}

	// TODO: Replace with actual email logic
	fmt.Printf("Mock OTP for %s is: %s\n", req.Email, otp)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTP sent to your email",
	})
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	NewPassword string `json:"new_password"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("reset_otp:%s", req.Email)
	storedOTP, err := db.RedisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil || storedOTP != req.OTP {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(hashedPassword)
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	// Cleanup OTP from Redis
	db.RedisClient.Del(ctx, redisKey)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}

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

func ForgotPasswordClientPage(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:8000/callback"
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
