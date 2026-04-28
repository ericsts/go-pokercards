<template>
  <div class="seat" :class="{ 'seat--me': isMe }">
    <div class="seat__card-wrap">
      <PokerCard
        :value="cardValue"
        :back="player.has_voted && !revealed"
        :empty="!player.has_voted"
        :winner="isWinner"
        class="seat__card"
      />
    </div>

    <div class="seat__avatar" :style="{ background: avatarColor }">
      {{ initials }}
    </div>

    <div class="seat__name">
      <span class="seat__master-icon" v-if="isMaster" title="Scrum Master">👑</span>
      {{ player.name }}
      <span v-if="isMe" class="seat__you">(você)</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import PokerCard from './PokerCard.vue'

const props = defineProps({
  player:   { type: Object,  required: true },
  isMaster: { type: Boolean, default: false },
  isMe:     { type: Boolean, default: false },
  revealed: { type: Boolean, default: false },
  winners:  { type: Array,   default: () => [] }
})

const COLORS = [
  '#e63946','#457b9d','#2a9d8f','#e9c46a',
  '#f4a261','#6a4c93','#06a77d','#d62839'
]

// Deterministic color from player ID
const avatarColor = computed(() => {
  let hash = 0
  for (const ch of props.player.id) hash = (hash * 31 + ch.charCodeAt(0)) & 0xffffffff
  return COLORS[Math.abs(hash) % COLORS.length]
})

const initials = computed(() =>
  props.player.name
    .split(' ')
    .slice(0, 2)
    .map(w => w[0]?.toUpperCase() ?? '')
    .join('')
)

const cardValue = computed(() => props.revealed ? (props.player.vote || '') : '')
const isWinner  = computed(() => props.revealed && props.winners.includes(props.player.vote))
</script>

<style scoped>
.seat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  width: 90px;
}

.seat__card { transform: scale(0.6); transform-origin: bottom center; }
.seat__card-wrap { height: 62px; display: flex; align-items: flex-end; }

.seat__avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.9rem;
  color: white;
  text-shadow: 0 1px 3px rgba(0,0,0,0.4);
  box-shadow: 0 2px 8px rgba(0,0,0,0.3);
  flex-shrink: 0;
}

.seat--me .seat__avatar {
  box-shadow: 0 0 0 3px var(--color-gold), 0 2px 8px rgba(0,0,0,0.3);
}

.seat__name {
  font-size: 0.78rem;
  color: var(--color-text);
  text-align: center;
  max-width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.seat__master-icon { font-size: 0.85rem; }
.seat__you { color: var(--color-text-muted); font-size: 0.7rem; }
</style>
