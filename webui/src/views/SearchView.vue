<template>
  <div class="search-container">
    <input
      v-model="searchName"
      type="text"
      placeholder="Enter username"
      class="search-input"
    />
    <button @click="search" class="search-button">Search</button>
    <ErrorMsg v-if="error" :msg="error" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'

const searchName = ref('')
const error = ref('')
const router = useRouter()

const search = async () => {
  error.value = ''
  if (!searchName.value) return
  try {
    const res = await axios.get(`/users/${searchName.value}`)
    console.log(res.data)
    router.push(`/user/${searchName.value}`)
  } catch (err) {
    error.value = "username doesn't exist"
  }
}
</script>

<style scoped>
.search-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
}
.search-input {
  width: 300px; 
  padding: 0.5rem;
  margin-bottom: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.search-button {
  padding: 0.5rem 1rem;
  background-color: #000000;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}
.search-button:hover {
  background-color: #cecece;
}
</style>
