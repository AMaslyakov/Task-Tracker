<template>
  <div class="toolbar">
    <div>
      <p class="eyebrow">Victory Group</p>
      <h1>Task Tracker</h1>
    </div>
    <button type="button">Новая задача</button>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 28px;
}

@media (max-width: 760px) {
  .toolbar {
    display: grid;
  }

  button {
    width: 100%;
  }
}
</style>
