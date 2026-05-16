<template>

  <div class="toolbar">


    <div>
      <p class="eyebrow">Victory Group</p>
      <h1>Task Tracker</h1>
    </div>
    <!-- <button type="button">Новая задача</button> -->

    <div class="user">
      <img class="user-icon" src="../assets/user.png" alt={{ user.name }}>
      <div class="user-data">
          <span class="user-id">Id: {{ user2.id}}</span>
          <span class="user-name">{{ user2.name}}</span>
          <span class="user-email">{{ user2.email}}</span>
      </div>
    </div>

  </div>
  
  
</template>

<script setup>
  import {user1, user2} from "../data/dashboard"
</script>

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

.user-icon{
  max-width: 40px;
}

.user-data{
  display: flex;
  flex-direction: column;
  align-items:self-start;
  font-size: 22px;
  row-gap: 8px;
}

</style>
