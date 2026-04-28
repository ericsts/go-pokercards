<template>
  <div class="home">
    <div class="home__logo">🃏</div>
    <h1 class="home__title">Planning Poker</h1>
    <p class="home__subtitle">Estime histórias com seu time de forma simples e divertida</p>

    <div class="home__box">
      <label class="home__label">Seu nome</label>
      <input
        v-model="name"
        placeholder="Ex: Alice"
        maxlength="30"
        @keyup.enter="createRoom"
        autofocus
      />
      <button class="home__btn" :disabled="!name.trim() || loading" @click="createRoom">
        {{ loading ? 'Criando…' : 'Criar nova sala' }}
      </button>
      <p v-if="errorMsg" class="home__error">{{ errorMsg }}</p>
    </div>

    <p class="home__hint">Você receberá um link para compartilhar com seu time.</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router  = useRouter()
const name    = ref(sessionStorage.getItem('poker_name') || '')
const loading = ref(false)
const errorMsg = ref('')

async function createRoom() {
  if (!name.value.trim()) return
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await fetch('/api/rooms', { method: 'POST' })
    if (!res.ok) throw new Error('Falha ao criar sala')
    const { room_id, creator_id } = await res.json()
    sessionStorage.setItem('poker_name', name.value.trim())
    sessionStorage.setItem(`poker_creator_${room_id}`, creator_id)
    router.push(`/room/${room_id}`)
  } catch (e) {
    errorMsg.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  gap: 12px;
}

.home__logo  { font-size: 4rem; }
.home__title { font-size: 2rem; font-weight: 800; color: var(--color-text); }
.home__subtitle {
  color: var(--color-text-muted);
  text-align: center;
  max-width: 380px;
  margin-bottom: 12px;
}

.home__box {
  width: 100%;
  max-width: 360px;
  background: var(--color-surface);
  border-radius: 16px;
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.4);
}

.home__label {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.home__btn {
  background: var(--color-gold);
  color: #1a1a1a;
  font-weight: 700;
  padding: 12px;
  font-size: 1rem;
}

.home__btn:not(:disabled):hover { background: #e0af20; }

.home__error { color: var(--color-error); font-size: 0.875rem; }
.home__hint  { color: var(--color-text-muted); font-size: 0.8rem; }
</style>
