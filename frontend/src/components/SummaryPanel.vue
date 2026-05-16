<template>
  <aside class="summary-panel">
    <h2>Сводка</h2>
    <dl>
      <div v-for="item in items" :key="item.label">
        <dt>{{ item.label }}</dt>
        <dd>{{ item.value }}</dd>
      </div>
    </dl>
  </aside>
</template>

<script setup>
defineProps({
  items: {
    type: Array,
    required: true
  }
})
</script>

<style scoped>
.summary-panel {
  border: 1px solid #dce3ee;
  border-radius: 8px;
  padding: 20px;
  background: #ffffff;
}

.summary-panel h2 {
  margin-bottom: 0;
  font-size: 22px;
}

dl {
  display: grid;
  gap: 12px;
  margin: 18px 0 0;
}

dl div {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 12px;
}

dt {
  color: #64748b;
}

dd {
  margin: 0;
  font-size: 24px;
  font-weight: 800;
}
</style>
