<template>
  <div class="profile-container">
    <div class="user-info">
      <h2>{{ username }}</h2>
      <div class="action-buttons">
        <button @click="toggleFollow" :class="{ active: isFollowing }">
          {{ isFollowing ? 'Unfollow' : 'Follow' }}
        </button>
        <button @click="toggleBan" :class="{ active: isBanned }">
          {{ isBanned ? 'Unban' : 'Ban' }}
        </button>
      </div>
      <div class="userStats">
        <span>Following ({{ followingList.length }})</span>
        <span>Banned ({{ bannedList.length }})</span>
      </div>
    </div>

    <div class="image-grid">
      <div v-for="img in images" :key="img.id" class="image-item">
        <img :src="img.imageurl" alt="photo" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../services/axios.js'

const route = useRoute()
const router = useRouter()

const username = ref(route.params.username || '')
const myUsername = ref(localStorage.getItem('username') || '')

const followingList = ref([])
const bannedList = ref([])
const images = ref([])

const isFollowing = ref(false)
const isBanned = ref(false)

const fetchProfile = async () => {
  try {
    const res = await axios.get(`/users/${username.value}`)
    followingList.value = res.data.Following ? res.data.Following.split(',').filter(v => v) : []
    bannedList.value = res.data.Banned ? res.data.Banned.split(',').filter(v => v) : []
  } catch (err) {
    console.error('Failed to load profile', err)
  }
}

const fetchImages = async () => {
  try {
    const res = await axios.get(`/users/${username.value}/photos`)
    images.value = res.data || []
  } catch (err) {
    console.error('Failed to load images', err)
  }
}

const fetchRelation = async () => {
  try {
    const res = await axios.get(`/users/${myUsername.value}`)
    const fList = res.data.Following ? res.data.Following.split(',').filter(v => v) : []
    const bList = res.data.Banned ? res.data.Banned.split(',').filter(v => v) : []
    isFollowing.value = fList.includes(username.value)
    isBanned.value = bList.includes(username.value)
  } catch (err) {
    console.error('Failed to load relation info', err)
  }
}

const toggleFollow = async () => {
  try {
    if (isFollowing.value) {
      await axios.delete(`/users/${myUsername.value}/follow`, { data: { username: username.value } })
      isFollowing.value = false
    } else {
      await axios.put(`/users/${myUsername.value}/follow`, { username: username.value })
      isFollowing.value = true
    }
  } catch (err) {
    console.error('Failed to toggle follow', err)
  }
}

const toggleBan = async () => {
  try {
    if (isBanned.value) {
      await axios.delete(`/users/${myUsername.value}/ban`, { data: { username: username.value } })
      isBanned.value = false
    } else {
      await axios.put(`/users/${myUsername.value}/ban`, { username: username.value })
      isBanned.value = true
    }
  } catch (err) {
    console.error('Failed to toggle ban', err)
  }
}

onMounted(() => {
  if (!myUsername.value) {
    router.push('/login')
    return
  }
  fetchProfile()
  fetchImages()
  fetchRelation()
})

watch(() => route.params.username, (newVal) => {
  username.value = newVal
  fetchProfile()
  fetchImages()
  fetchRelation()
})
</script>

<style scoped>
.profile-container {
  padding: 1rem;
}

.user-info {
  border: 1px solid #ddd;
  padding: 1rem;
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.action-buttons button {
  cursor: pointer;
  margin: 5px;
  padding: 8px 12px;
  background-color: #000000;
  color: white;
  border: none;
  border-radius: 5px;
  transition: background-color 0.3s ease;
}

.action-buttons button.active {
  background-color: #4caf50;
  color: white;
}

.action-buttons button:hover {
  background-color: #707070;
}

.stats span {
  margin-right: 1rem;
  color: white;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 10px;
}

.image-item img {
  width: 100%;
  height: 300px;
  object-fit: cover;
}

.userStats span {
  margin: 5px;
  padding: 8px 12px;
  background-color: #000000;
  color: white;
  border: none;
  border-radius: 5px;
  transition: background-color 0.3s ease;
}

</style>
