// Dashboard Logic
let userExercises = [];
let workoutExercises = [];

document.addEventListener('DOMContentLoaded', () => {
    // Check authentication
    if (!api.isAuthenticated()) {
        window.location.href = '/';
        return;
    }

    // Initialize dashboard
    initializeDashboard();
    loadUserExercises();
});

function initializeDashboard() {
    // Logout
    document.getElementById('logoutBtn').addEventListener('click', () => {
        api.logout();
        window.location.href = '/';
    });

    // View navigation
    const sidebarButtons = document.querySelectorAll('.sidebar-btn');
    sidebarButtons.forEach(btn => {
        btn.addEventListener('click', () => {
            const view = btn.dataset.view;
            switchView(view);
            
            // Update active state
            sidebarButtons.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
        });
    });

    // Create workout button in header
    document.getElementById('createWorkoutBtn').addEventListener('click', () => {
        switchView('create-workout');
        document.querySelectorAll('.sidebar-btn').forEach(b => b.classList.remove('active'));
        document.querySelector('[data-view="create-workout"]').classList.add('active');
    });

    // Exercise modal
    document.getElementById('createExerciseBtn').addEventListener('click', openExerciseModal);
    document.getElementById('closeExerciseModal').addEventListener('click', closeExerciseModal);
    document.getElementById('cancelExerciseBtn').addEventListener('click', closeExerciseModal);
    document.getElementById('modalOverlay').addEventListener('click', closeExerciseModal);

    // Exercise form submission
    document.getElementById('exerciseForm').addEventListener('submit', handleCreateExercise);

    // Workout form
    document.getElementById('addExerciseBtn').addEventListener('click', addExerciseToWorkout);
    document.getElementById('cancelWorkoutBtn').addEventListener('click', () => {
        switchView('workouts');
        document.querySelectorAll('.sidebar-btn').forEach(b => b.classList.remove('active'));
        document.querySelector('[data-view="workouts"]').classList.add('active');
    });
    document.getElementById('workoutForm').addEventListener('submit', handleCreateWorkout);
}

function switchView(viewName) {
    // Hide all views
    document.querySelectorAll('.view-container').forEach(view => {
        view.classList.add('hidden');
    });

    // Show selected view
    const targetView = document.getElementById(`${viewName}View`);
    if (targetView) {
        targetView.classList.remove('hidden');
    }

    // Load data for specific views
    if (viewName === 'exercises') {
        displayExercises();
    } else if (viewName === 'create-workout') {
        resetWorkoutForm();
    }
}

async function loadUserExercises() {
    try {
        userExercises = await api.getUserExercises();
        console.log('Loaded exercises:', userExercises);
    } catch (error) {
        console.error('Failed to load exercises:', error);
        userExercises = [];
    }
}

function displayExercises() {
    const exercisesGrid = document.getElementById('exercisesGrid');
    
    if (!userExercises || userExercises.length === 0) {
        exercisesGrid.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                </svg>
                <h3>No exercises yet</h3>
                <p>Create your first exercise to get started</p>
            </div>
        `;
        return;
    }

    exercisesGrid.innerHTML = userExercises.map(exercise => `
        <div class="exercise-card">
            <div class="card-header">
                <div>
                    <h3 class="card-title">${exercise.name}</h3>
                </div>
                <span class="card-badge">${exercise.muscle_group}</span>
            </div>
        </div>
    `).join('');
}

// Exercise Modal Functions
function openExerciseModal() {
    document.getElementById('createExerciseModal').classList.remove('hidden');
    document.getElementById('exerciseName').focus();
}

function closeExerciseModal() {
    document.getElementById('createExerciseModal').classList.add('hidden');
    document.getElementById('exerciseForm').reset();
    document.getElementById('exerciseFormError').classList.remove('show');
}

async function handleCreateExercise(e) {
    e.preventDefault();
    
    const name = document.getElementById('exerciseName').value;
    const muscle_group = document.getElementById('exerciseMuscleGroup').value;
    const errorElement = document.getElementById('exerciseFormError');

    try {
        const result = await api.createExercise(name, muscle_group);
        
        // Add to local exercises array
        userExercises.push({
            exercise_id: result.exercise_id,
            name,
            muscle_group
        });

        // Close modal and refresh exercises view
        closeExerciseModal();
        
        if (!document.getElementById('exercisesView').classList.contains('hidden')) {
            displayExercises();
        }

        showNotification('Exercise created successfully!');
    } catch (error) {
        errorElement.textContent = error.message || 'Failed to create exercise';
        errorElement.classList.add('show');
    }
}

// Workout Form Functions
function addExerciseToWorkout() {
    if (userExercises.length === 0) {
        showNotification('Please create some exercises first!', 'error');
        return;
    }

    const exercisesList = document.getElementById('exercisesList');
    const exerciseIndex = workoutExercises.length;
    
    const exerciseItem = document.createElement('div');
    exerciseItem.className = 'exercise-item';
    exerciseItem.dataset.index = exerciseIndex;
    
    exerciseItem.innerHTML = `
        <div class="exercise-item-header">
            <h4>Exercise ${exerciseIndex + 1}</h4>
            <button type="button" class="btn-remove" onclick="removeExercise(${exerciseIndex})">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"/>
                    <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
            </button>
        </div>
        <div class="exercise-item-fields">
            <div class="form-group">
                <label>Exercise</label>
                <select required data-field="exercise_id">
                    <option value="">Select exercise</option>
                    ${userExercises.map(ex => `
                        <option value="${ex.exercise_id}">${ex.name}</option>
                    `).join('')}
                </select>
            </div>
            <div class="form-group">
                <label>Sets</label>
                <input type="number" min="1" required placeholder="3" data-field="sets">
            </div>
            <div class="form-group">
                <label>Reps</label>
                <input type="number" min="1" required placeholder="10" data-field="reps">
            </div>
            <div class="form-group">
                <label>Weight (kg)</label>
                <input type="number" min="0" step="0.5" required placeholder="50" data-field="weight">
            </div>
        </div>
    `;
    
    exercisesList.appendChild(exerciseItem);
    workoutExercises.push({ index: exerciseIndex });
}

function removeExercise(index) {
    const exerciseItem = document.querySelector(`.exercise-item[data-index="${index}"]`);
    if (exerciseItem) {
        exerciseItem.remove();
        workoutExercises = workoutExercises.filter(ex => ex.index !== index);
    }
}

function resetWorkoutForm() {
    document.getElementById('workoutForm').reset();
    document.getElementById('exercisesList').innerHTML = '';
    workoutExercises = [];
    document.getElementById('workoutFormError').classList.remove('show');
    
    // Set default datetime to current time
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('workoutScheduledAt').value = now.toISOString().slice(0, 16);
}

async function handleCreateWorkout(e) {
    e.preventDefault();
    
    const errorElement = document.getElementById('workoutFormError');
    errorElement.classList.remove('show');

    // Get form values
    const title = document.getElementById('workoutTitle').value;
    const description = document.getElementById('workoutDescription').value;
    const comment = document.getElementById('workoutComment').value;
    const scheduled_at = new Date(document.getElementById('workoutScheduledAt').value).toISOString();

    // Get exercises from form
    const exerciseItems = document.querySelectorAll('.exercise-item');
    
    if (exerciseItems.length === 0) {
        errorElement.textContent = 'Please add at least one exercise';
        errorElement.classList.add('show');
        return;
    }

    const exercises = [];
    
    exerciseItems.forEach((item, index) => {
        const exercise_id = parseInt(item.querySelector('[data-field="exercise_id"]').value);
        const sets = parseInt(item.querySelector('[data-field="sets"]').value);
        const reps = parseInt(item.querySelector('[data-field="reps"]').value);
        const weight = parseFloat(item.querySelector('[data-field="weight"]').value);

        if (exercise_id && sets && reps && weight >= 0) {
            exercises.push({
                exercise_id,
                sets,
                reps,
                weight,
                order: index + 1
            });
        }
    });

    if (exercises.length === 0) {
        errorElement.textContent = 'Please fill in all exercise details';
        errorElement.classList.add('show');
        return;
    }

    // Prepare workout data
    const workoutData = {
        title,
        description,
        comment,
        scheduled_at,
        exercises
    };

    // Disable submit button
    const submitBtn = e.target.querySelector('button[type="submit"]');
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span>Creating workout...</span>';

    try {
        await api.createWorkout(workoutData);
        
        showNotification('Workout created successfully!');
        
        // Switch to workouts view
        switchView('workouts');
        document.querySelectorAll('.sidebar-btn').forEach(b => b.classList.remove('active'));
        document.querySelector('[data-view="workouts"]').classList.add('active');
        
        resetWorkoutForm();
    } catch (error) {
        errorElement.textContent = error.message || 'Failed to create workout';
        errorElement.classList.add('show');
        submitBtn.disabled = false;
        submitBtn.innerHTML = `
            <span>Create Workout</span>
            <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"/>
            </svg>
        `;
    }
}

// Notification System
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.style.cssText = `
        position: fixed;
        top: 2rem;
        right: 2rem;
        padding: 1rem 2rem;
        background: ${type === 'success' ? 'var(--success-color)' : 'var(--error-color)'};
        color: white;
        border-radius: 12px;
        box-shadow: var(--shadow-lg);
        z-index: 10000;
        animation: slideInRight 0.3s ease;
        font-weight: 600;
    `;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.animation = 'slideOutRight 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Add CSS for notification animations
const style = document.createElement('style');
style.textContent = `
    @keyframes slideInRight {
        from {
            opacity: 0;
            transform: translateX(100px);
        }
        to {
            opacity: 1;
            transform: translateX(0);
        }
    }
    
    @keyframes slideOutRight {
        from {
            opacity: 1;
            transform: translateX(0);
        }
        to {
            opacity: 0;
            transform: translateX(100px);
        }
    }
`;
document.head.appendChild(style);
