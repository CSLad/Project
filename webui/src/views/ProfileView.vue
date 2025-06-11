<template>
  <div class="profile-container">
    <div class="user-info">
      <h2>{{ username }}</h2>
      <button class="change-btn" @click="changeUsername">Change Username</button>
      <div class="stats">
        <span @click="showFollowingModal = true">Following ({{ followingList.length }})</span>
        <span @click="showBannedModal = true">Banned ({{ bannedList.length }})</span>
      </div>
    </div>

    <div class="image-grid">
      <div class="upload-tile" @click="showUploadModal = true">⬆️</div>
      <div v-for="img in images" :key="img.id" class="image-item">
        <img :src="img.imageurl" alt="photo" />
      </div>
    </div>

    <div v-if="showFollowingModal" class="modal-overlay">
      <div class="modal">
        <h3>Following</h3>
        <ul>
          <li v-for="user in followingList" :key="user">
            <span>{{ user }}</span>
            <button @click="toggleFollow(user)" :class="{ active: !isFollowing(user) }">
              {{ isFollowing(user) ? 'Unfollow' : 'Follow' }}
            </button>
          </li>
        </ul>
        <button @click="showFollowingModal = false">Close</button>
      </div>
    </div>

    <div v-if="showBannedModal" class="modal-overlay">
      <div class="modal">
        <h3>Banned</h3>
        <ul>
          <li v-for="user in bannedList" :key="user">
            <span>{{ user }}</span>
            <button @click="toggleBan(user)" :class="{ active: !isBanned(user) }">
              {{ isBanned(user) ? 'Unban' : 'Ban' }}
            </button>
          </li>
        </ul>
        <button @click="showBannedModal = false">Close</button>
      </div>
    </div>

    <div v-if="showUploadModal" class="modal-overlay">
      <div class="modal">
        <h3>Upload Image</h3>
        <input v-model="newImageUrl" placeholder="Image URL" />
        <button @click="uploadImage">Upload</button>
        <button @click="showUploadModal = false">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const followingList = ref([])
const bannedList = ref([])
const images = ref([])
const showFollowingModal = ref(false)
const showBannedModal = ref(false)
const showUploadModal = ref(false)
const newImageUrl = ref('')

const fetchProfile = async () => {
  try {
    const res = await axios.get(`/users/${username.value}`)
    followingList.value = res.data.following ? res.data.following.split(',').filter(v => v) : []
    bannedList.value = res.data.banned ? res.data.banned.split(',').filter(v => v) : []
  } catch (err) {
    console.error('Failed to load profile', err)
  }
}

const fetchImages = async () => {
  try {
    const res = await axios.get(`/images/${username.value}`)
    images.value = res.data || []
  } catch (err) {
    console.error('Failed to load images', err)
  }
}

const changeUsername = async () => {
  const newName = prompt('New username:')
  if (!newName) return
  try {
    await axios.put(`/users/${username.value}`, { username: newName })
    localStorage.setItem('username', newName)
    username.value = newName
    await fetchProfile()
    await fetchImages()
  } catch (err) {
    console.error('Failed to change username', err)
  }
}

const isFollowing = (user) => followingList.value.includes(user)
const isBanned = (user) => bannedList.value.includes(user)

const toggleFollow = async (user) => {
  try {
    if (isFollowing(user)) {
      await axios.delete(`/users/${username.value}/follow`, { data: { username: user } })
      followingList.value = followingList.value.filter(u => u !== user)
    } else {
      await axios.put(`/users/${username.value}/follow`, { username: user })
      followingList.value.push(user)
    }
  } catch (err) {
    console.error('Failed to toggle follow', err)
  }
}

const toggleBan = async (user) => {
  try {
    if (isBanned(user)) {
      await axios.delete(`/users/${username.value}/ban`, { data: { username: user } })
      bannedList.value = bannedList.value.filter(u => u !== user)
    } else {
      await axios.put(`/users/${username.value}/ban`, { username: user })
      bannedList.value.push(user)
    }
  } catch (err) {
    console.error('Failed to toggle ban', err)
  }
}

const uploadImage = async () => {
  if (!newImageUrl.value) return
  try {
    await axios.post('/images', { username: username.value, imageurl: newImageUrl.value })
    newImageUrl.value = ''
    showUploadModal.value = false
    await fetchImages()
  } catch (err) {
    console.error('Failed to upload image', err)
  }
}

onMounted(() => {
  if (!username.value) {
    router.push('/login')
    return
  }
  fetchProfile()
  fetchImages()
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

.change-btn {
  width: fit-content;
}

.stats span {
  margin-right: 1rem;
  cursor: pointer;
  color: blue;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;
}

.upload-tile, .image-item img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.upload-tile {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  cursor: pointer;
  font-size: 2rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: white;
  padding: 1rem;
  border-radius: 4px;
  max-width: 400px;
  width: 90%;
}

.modal ul {
  list-style: none;
  padding: 0;
}

.modal li {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.modal button.active {
  background: #4caf50;
  color: white;
}

.modal button {
  padding: 0.3rem 0.6rem;
}
</style>
