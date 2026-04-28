<template>
  <div class="deck">
    <div class="deck__header">
      <p class="deck__label">Selecione sua carta:</p>
      <Transition name="badge">
        <span v-if="selected" class="deck__badge">✓ Sua carta: {{ selected }}</span>
      </Transition>
    </div>
    <div class="deck__cards">
      <PokerCard
        v-for="val in cards"
        :key="val"
        :value="val"
        :selected="selected === val"
        :selectable="!disabled"
        @select="emit('vote', val)"
      />
    </div>
  </div>
</template>

<script setup>
import PokerCard from './PokerCard.vue'

defineProps({
  cards:    { type: Array,  required: true },
  selected: { type: String, default: '' },
  disabled: { type: Boolean, default: false }
})
const emit = defineEmits(['vote'])
</script>

<style scoped>
.deck {
  padding: 16px 20px 24px;
  background: var(--color-surface);
  border-top: 1px solid rgba(255,255,255,0.08);
}

.deck__header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.deck__label {
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--color-text-muted);
}

.deck__badge {
  font-size: 0.8rem;
  font-weight: 700;
  background: #c0392b;
  color: #fff;
  border-radius: 999px;
  padding: 3px 12px;
}

.badge-enter-active, .badge-leave-active { transition: opacity 0.25s, transform 0.25s; }
.badge-enter-from  { opacity: 0; transform: scale(0.8); }
.badge-leave-to    { opacity: 0; transform: scale(0.8); }

.deck__cards {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
}
</style>
