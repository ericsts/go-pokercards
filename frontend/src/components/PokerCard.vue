<template>
  <div
    class="card"
    :class="{
      'card--back':       back,
      'card--empty':      empty,
      'card--selected':   selected,
      'card--selectable': selectable,
      'card--winner':     winner
    }"
    @click="selectable && emit('select')"
  >
    <template v-if="!back && !empty && value">
      <span class="card__corner card__corner--tl">{{ value }}</span>
      <span class="card__center">{{ value }}</span>
      <span class="card__corner card__corner--br">{{ value }}</span>
    </template>

    <template v-else-if="back">
      <div class="card__back-pattern">
        <div class="card__back-inner" />
      </div>
    </template>

    <template v-else>
      <span class="card__placeholder">?</span>
    </template>
  </div>
</template>

<script setup>
defineProps({
  value:     { type: String,  default: '' },
  back:      { type: Boolean, default: false },
  empty:     { type: Boolean, default: false },
  selected:  { type: Boolean, default: false },
  selectable:{ type: Boolean, default: false },
  winner:    { type: Boolean, default: false }
})
const emit = defineEmits(['select'])
</script>

<style scoped>
.card {
  position: relative;
  width: 70px;
  height: 100px;
  border-radius: var(--radius-card);
  border: 2px solid #ddd;
  background: var(--color-card-bg);
  box-shadow: var(--shadow-card);
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease, background 0.18s ease;
  flex-shrink: 0;
}

.card--selectable { cursor: pointer; }
.card--selectable:hover {
  transform: translateY(-10px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.3);
}

.card--selected {
  transform: translateY(-18px) scale(1.08);
  background: #c0392b;
  border-color: #e74c3c;
  box-shadow: 0 0 0 3px rgba(231,76,60,0.45), 0 18px 32px rgba(0,0,0,0.45);
}

.card--selected .card__center,
.card--selected .card__corner {
  color: #fff;
}

.card--winner {
  border-color: var(--color-gold);
  box-shadow: 0 0 0 4px var(--color-gold-glow), 0 8px 20px rgba(212,160,23,0.4);
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 4px var(--color-gold-glow); }
  50%       { box-shadow: 0 0 0 8px rgba(212,160,23,0.2); }
}

.card--back {
  background: linear-gradient(145deg, #1a3d6b 0%, #2d5fa0 50%, #1a3d6b 100%);
  border-color: #3a5fa8;
}

.card--empty {
  background: rgba(255,255,255,0.05);
  border: 2px dashed rgba(255,255,255,0.2);
  box-shadow: none;
}

.card__center {
  font-size: 1.8rem;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1;
}

.card__corner {
  position: absolute;
  font-size: 0.65rem;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1;
}

.card__corner--tl { top: 5px;  left: 6px; }
.card__corner--br { bottom: 5px; right: 6px; transform: rotate(180deg); }

.card__back-pattern {
  width: 56px;
  height: 86px;
  border-radius: 6px;
  border: 2px solid rgba(255,255,255,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.card__back-inner {
  width: 36px;
  height: 66px;
  border-radius: 4px;
  background: repeating-linear-gradient(
    45deg,
    rgba(255,255,255,0.15) 0px,
    rgba(255,255,255,0.15) 2px,
    transparent 2px,
    transparent 8px
  );
  border: 1px solid rgba(255,255,255,0.2);
}

.card__placeholder {
  font-size: 1.4rem;
  color: rgba(255,255,255,0.15);
}
</style>
