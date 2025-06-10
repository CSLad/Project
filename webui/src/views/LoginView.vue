<template>
    <div class="login-container">
      <h1>Welcome to WasaPhoto 📸</h1>
      <form @submit.prevent="login">
        <input
          v-model="username"
          type="text"
          placeholder="Enter your username"
          required
          class="login-input"
        />
        <button type="submit" class="login-button">Login</button>
      </form>
      <ErrorMsg v-if="error" :message="error" />
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import ErrorMsg from '../components/ErrorMsg.vue'
  import axios from '../services/axios.js'
  
  const username = ref('')
  const error = ref('')
  const router = useRouter()
  
  const login = async () => {
  error.value = ''
  try {
    const res = await axios.post('/session', { Username: username.value })
    console.log('✅ API success:', res.data)
    localStorage.setItem('username', res.data.Username)
    router.push('/home')
  } catch (err) {
    error.value = 'Login failed. Try again.'
    console.error('❌ API error:', err)
  }
}

  </script>
  
  <style scoped>
  .login-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4rem 1rem;
    max-width: 400px;
    margin: 0 auto;
    text-align: center;
  }
  
  .login-input {
    width: 100%;
    padding: 0.8rem;
    margin: 1rem 0;
    border: 1px solid #ccc;
    border-radius: 8px;
  }
  
  .login-button {
    padding: 0.8rem 2rem;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
    transition: background-color 0.3s;
  }
  
  .login-button:hover {
    background-color: #369f70;
  }
  </style>
  