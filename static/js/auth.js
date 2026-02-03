// Authentication Logic
document.addEventListener('DOMContentLoaded', () => {
    // Check if already authenticated
    if (api.isAuthenticated()) {
        window.location.href = '/dashboard';
        return;
    }

    // Get DOM elements
    const signInForm = document.getElementById('signInForm');
    const signUpForm = document.getElementById('signUpForm');
    const showSignUpBtn = document.getElementById('showSignUpBtn');
    const showSignInBtn = document.getElementById('showSignInBtn');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const loginError = document.getElementById('loginError');
    const registerError = document.getElementById('registerError');

    // Toggle between sign in and sign up
    showSignUpBtn.addEventListener('click', (e) => {
        e.preventDefault();
        signInForm.classList.add('hidden');
        signUpForm.classList.remove('hidden');
        clearErrors();
    });

    showSignInBtn.addEventListener('click', (e) => {
        e.preventDefault();
        signUpForm.classList.add('hidden');
        signInForm.classList.remove('hidden');
        clearErrors();
    });

    // Handle login
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearErrors();

        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;

        // Disable submit button
        const submitBtn = loginForm.querySelector('button[type="submit"]');
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<span>Signing in...</span>';

        try {
            await api.signIn(email, password);
            // Redirect to dashboard
            window.location.href = '/dashboard';
        } catch (error) {
            showError(loginError, error.message || 'Invalid credentials. Please try again.');
            submitBtn.disabled = false;
            submitBtn.innerHTML = `
                <span>Sign In</span>
                <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M5 12h14M12 5l7 7-7 7"/>
                </svg>
            `;
        }
    });

    // Handle registration
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearErrors();

        const name = document.getElementById('registerName').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;

        // Validate inputs
        if (name.length < 3) {
            showError(registerError, 'Name must be at least 3 characters long');
            return;
        }

        if (password.length < 6) {
            showError(registerError, 'Password must be at least 6 characters long');
            return;
        }

        // Disable submit button
        const submitBtn = registerForm.querySelector('button[type="submit"]');
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<span>Creating account...</span>';

        try {
            await api.signUp(name, email, password);
            
            // Auto sign in after registration
            await api.signIn(email, password);
            
            // Redirect to dashboard
            window.location.href = '/dashboard';
        } catch (error) {
            showError(registerError, error.message || 'Failed to create account. Please try again.');
            submitBtn.disabled = false;
            submitBtn.innerHTML = `
                <span>Create Account</span>
                <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M5 12h14M12 5l7 7-7 7"/>
                </svg>
            `;
        }
    });

    // Helper functions
    function showError(element, message) {
        element.textContent = message;
        element.classList.add('show');
    }

    function clearErrors() {
        loginError.textContent = '';
        loginError.classList.remove('show');
        registerError.textContent = '';
        registerError.classList.remove('show');
    }

    // Add input animations
    const inputs = document.querySelectorAll('input');
    inputs.forEach(input => {
        input.addEventListener('focus', () => {
            input.parentElement.style.transform = 'translateY(-2px)';
        });
        
        input.addEventListener('blur', () => {
            input.parentElement.style.transform = 'translateY(0)';
        });
    });
});
