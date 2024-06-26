const passwordInput = document.getElementById("password");
const showPasswordCheckbox = document.getElementById("show-password");

showPasswordCheckbox.addEventListener("change", function() {
  if (showPasswordCheckbox.checked) {
    passwordInput.type = "text";
  } else {
    passwordInput.type = "password";
  }
});

const authForm = document.getElementById("auth-form");
const authButton = document.getElementById("auth-button");
const toggleButton = document.getElementById("toggle-button");  
let isLoginMode = true; 

toggleButton.addEventListener("click", function () {
    isLoginMode = !isLoginMode; 
    if (isLoginMode) {
        authForm.action = "/login";
        authButton.textContent = "Log in";
        toggleButton.textContent = "Switch to Sign up";
    } else {
        authForm.action = "/signup";
        authButton.textContent = "Sign up";
        toggleButton.textContent = "Switch to Log in";
    }
});