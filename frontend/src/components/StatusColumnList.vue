<template>
    <div class="statuses">
        <div class="status-column" v-for="status_name in config.statuses">
            <div class="status-name"><span>{{ status_name }}</span></div>
            <template v-for="task in tasks" :key="task.id">
                <div class="card" v-if="status_name === task.status">
                    <TaskCard :task="task" :key="task.id" />
                </div>
            </template>
            
        </div>
    </div>
</template> 

<script setup>
    import { command1} from '../data/dashboard';
    import TaskCard from './TaskCard.vue';
    const config = command1.config_dashboard;
    defineProps({
    tasks: {
        type: Array,
        required: true
    }
});

</script>

<style scoped>
    .status-column{
        background-color: rgb(180, 149, 149);
        min-width: 200px;
        border: 2px;
        border-color: black;
    }
    .statuses{
        display: flex;
        flex-direction: row;
        gap: 10px;
    }

    .status-name{
        display: grid;
        place-items: center;
        font-size: 20px;
        font-weight: 700;
        color: black;
        /* text-align: center; */
        min-height: 45px;
    }
</style>