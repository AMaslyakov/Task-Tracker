<template>
    <div class="statuses" :style="{ '--statuses-min-width': minScrollWidth }">
        <div class="statuses-scrollbar" ref="topScrollRef" @scroll="syncColumnsScroll">
            <div class="statuses-scrollbar-content"></div>
        </div>

        <div class="statuses-viewport" ref="columnsScrollRef" @scroll="syncTopScroll">
            <div class="status-columns">
                <div class="status-column" v-for="status_name in config.statuses" :key="status_name">
                    <div class="status-name"><span>{{ status_name }}</span></div>
                    <template v-for="task in tasks" :key="task.id">
                        <div class="card" v-if="status_name === task.status">
                            <TaskCard :task="task" :key="task.id" />
                        </div>
                    </template>
                </div>
            </div>

        </div>
    </div>
</template> 

<script setup>
    import { computed, ref } from 'vue';
    import { command1 } from '../data/dashboard';
    import TaskCard from './TaskCard.vue';

    const config = command1.config_dashboard;

    defineProps({
        tasks: {
            type: Array,
            required: true
        }
    });

    const topScrollRef = ref(null);
    const columnsScrollRef = ref(null);
    const columnMinWidth = 200;
    const columnGap = 10;

    const minScrollWidth = computed(() => {
        const statusesCount = config.statuses.length;
        const gapsWidth = Math.max(statusesCount - 1, 0) * columnGap;

        return `${statusesCount * columnMinWidth + gapsWidth}px`;
    });

    function syncColumnsScroll() {
        if (!topScrollRef.value || !columnsScrollRef.value) {
            return;
        }

        columnsScrollRef.value.scrollLeft = topScrollRef.value.scrollLeft;
    }

    function syncTopScroll() {
        if (!topScrollRef.value || !columnsScrollRef.value) {
            return;
        }

        topScrollRef.value.scrollLeft = columnsScrollRef.value.scrollLeft;
    }

</script>

<style scoped>

    .statuses{
        grid-column: 1 / -1;
        width: 100%;
        max-width: 100%;
        min-width: 0;
        display: grid;
        gap: 8px;
    }

    .statuses-scrollbar{
        width: 100%;
        max-width: 100%;
        overflow-x: auto;
        overflow-y: hidden;
    }

    .statuses-scrollbar-content{
        width: max(100%, var(--statuses-min-width));
        height: 1px;
    }

    .statuses-viewport{
        width: 100%;
        max-width: 100%;
        min-width: 0;
        overflow-x: auto;
        scrollbar-width: none;
    }

    .statuses-viewport::-webkit-scrollbar{
        display: none;
    }

    .status-columns{
        display: flex;
        flex-direction: row;
        gap: 10px;
        width: max(100%, var(--statuses-min-width));
    }

    .status-column{
        /* background-color: rgb(240, 223, 223); */
        flex: 1 0 200px;
        min-width: 0;
        /* border: 2px solid black; */
        /* border-radius: 14px 14px 0 0; */

    }

    .status-name{
        display: grid;
        place-items: center;
        font-size: 20px;
        font-weight: 700;
        color: black;

        min-height: 45px;
        border: 2px solid black;
        border-radius: 14px 14px 0 0;
    }

    .card{
        margin-bottom: 10px;
    }
</style>
