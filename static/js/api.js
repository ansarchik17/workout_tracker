// API Client for Workout Tracker
class API {
    constructor() {
        this.baseURL = window.location.origin;
        this.token = localStorage.getItem('authToken');
    }

    setToken(token) {
        this.token = token;
        localStorage.setItem('authToken', token);
    }

    removeToken() {
        this.token = null;
        localStorage.removeItem('authToken');
    }

    getHeaders(includeAuth = false) {
        const headers = {
            'Content-Type': 'application/json',
        };

        if (includeAuth && this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        return headers;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const config = {
            ...options,
            headers: this.getHeaders(options.auth),
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (response.ok) {
                throw new Error(data.error || 'Request failed');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // Auth endpoints
    async signUp(name, email, password) {
        return this.request('/user/signUp', {
            method: 'POST',
            body: JSON.stringify({ name, email, password }),
        });
    }

    async signIn(email, password) {
        const data = await this.request('/user/signIn', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        });
        
        if (data.token) {
            this.setToken(data.token);
        }
        
        return data;
    }

    // Exercise endpoints
    async createExercise(name, muscle_group) {
        return this.request('/exercise', {
            method: 'POST',
            auth: true,
            body: JSON.stringify({ name, muscle_group }),
        });
    }

    async getUserExercises() {
        return this.request('/user/exercises', {
            method: 'GET',
            auth: true,
        });
    }

    // Workout endpoints
    async createWorkout(workoutData) {
        return this.request('/workout', {
            method: 'POST',
            auth: true,
            body: JSON.stringify(workoutData),
        });
    }

    async getUserWorkouts() {
        return this.request('/user/workout', {
            method: 'GET',
            auth: true,
        });
    }
    

    // Helper to check if user is authenticated
    isAuthenticated() {
        return !!this.token;
    }

    // Logout
    logout() {
        this.removeToken();
    }
}

// Create global API instance
const api = new API();
