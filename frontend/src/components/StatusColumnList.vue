<template>
    <div class="statuses" :style="{ '--statuses-min-width': minScrollWidth }">
        <!-- Верхний кастомный скроллбар -->
        <div class="statuses-scrollbar" ref="topScrollRef" @scroll="syncColumnsScroll">
            <div class="statuses-scrollbar-content"></div>
        </div>

        <!-- Основная область с колонками -->
        <div class="statuses-viewport" ref="columnsScrollRef" @scroll="syncTopScroll">
            <div class="status-columns">
                <div class="status-column" v-for="status_name in config.statuses" :key="status_name">
                    <div class="status-name"><span>{{ status_name }}</span></div>

                    <!-- Зона перетаскивания карточек -->
                    <Sortable
                        :list="tasksByStatus[status_name] || []"
                        itemKey="id"
                        tag="div"
                        class="status-column-body"
                        :options="{ group: 'tasks-kanban', animation: 200, ghostClass: 'ghost-card' }"
                        @change="(event) => handleTaskMove(event, status_name)"
                    >
                        <template #item="{ element: task }">
                            <!-- Карточка с динамическим классом приоритета -->
                            <div class="card" :key="task.id" :class="priorityClass(task.priority)">
                                <TaskCard :task="task" />
                            </div>
                        </template>
                    </Sortable>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { computed, ref } from 'vue';
    import { Sortable } from 'sortablejs-vue3';
    import TaskCard from './TaskCard.vue';

    const props = defineProps({
        tasks: {
            type: Array,
            required: true
        },
        command: {
            type: Object,
            required: true
        }
    });

    const emit = defineEmits(['update:tasks', 'task-status-changed']);

    const topScrollRef = ref(null);
    const columnsScrollRef = ref(null);
    const columnMinWidth = 200;
    const columnGap = 10;
    const config = computed(() => props.command.config_dashboard);

    const minScrollWidth = computed(() => {
        const statusesCount = config.value.statuses.length;
        const gapsWidth = Math.max(statusesCount - 1, 0) * columnGap;
        return `${statusesCount * columnMinWidth + gapsWidth}px`;
    });

    // Функция вычисления веса на основе данных вашей таблицы SQL priorities
    function getTaskWeight(task) {
        // Получаем значение приоритета (это может быть id, текстовое имя или объект)
        const pValue = task.priority_id || task.priority;
        if (pValue == null) return 999; // Если приоритет не указан, отправляем вниз

        const pStr = String(pValue).trim().toLowerCase();

        // Карта весов строго по вашей SQL-таблице (id и названия)
        const WEIGHT_MAP = {
            // По ID
            '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6,
            // По именам (priority_name)
            'critical': 1, 'high': 2, 'medium': 3, 'low': 4, 'backlog': 5, 'blocked': 6,
            // По ключевым словам (keywords)
            'crit': 1, 'urgent': 1,
            'important': 2, 'major': 2,
            'normal': 3, 'standard': 3,
            'minor': 4, 'trivial': 4,
            'later': 5, 'someday': 5, 'plan': 5,
            'block': 6
        };

        return WEIGHT_MAP[pStr] !== undefined ? WEIGHT_MAP[pStr] : 999;
    }

    // Группировка задач по текстовым статусам с автоматической сортировкой по важности (weight)
    const tasksByStatus = computed(() => {
        const map = {};
        if (config.value?.statuses) {
            config.value.statuses.forEach(status => {
                // Фильтруем исходные задачи по названию колонки
                const filtered = props.tasks.filter(task => {
                    // Поддержка фильтрации как по имени (task.status), так и по объекту (task.status.name)
                    const currentStatus = typeof task.status === 'object' ? task.status?.name : task.status;
                    return String(currentStatus).trim().toLowerCase() === String(status).trim().toLowerCase();
                });

                // Сортируем: карточки с меньшим weight (Critical = 1) поднимутся наверх (a - b)
                map[status] = filtered.sort((a, b) => getTaskWeight(a) - getTaskWeight(b));
            });
        }
        return map;
    });

    // Обработка переноса карточки в другую колонку
    function handleTaskMove(event, newStatus) {
        if (event.added) {
            const movedTask = event.added.element;

            const updatedTasks = props.tasks.map(task => {
                if (task.id === movedTask.id) {
                    return { ...task, status: newStatus };
                }
                return task;
            });

            emit('update:tasks', updatedTasks);
            emit('task-status-changed', { taskId: movedTask.id, newStatus });
        }
    }

    // Синхронизация скролла
    function syncColumnsScroll() {
        if (!topScrollRef.value || !columnsScrollRef.value) return;
        columnsScrollRef.value.scrollLeft = topScrollRef.value.scrollLeft;
    }

    function syncTopScroll() {
        if (!topScrollRef.value || !columnsScrollRef.value) return;
        topScrollRef.value.scrollLeft = columnsScrollRef.value.scrollLeft;
    }

    // Карта соответствия приоритетов css-классам (поддерживает числа и строки)
    const PRIORITY_MAP = {
      1: 'priority-critical', 'critical': 'priority-critical',
      2: 'priority-high', 'high': 'priority-high',
      3: 'priority-medium', 'medium': 'priority-medium',
      4: 'priority-low', 'low': 'priority-low',
      5: 'priority-backlog', 'backlog': 'priority-backlog',
      6: 'priority-blocked', 'blocked': 'priority-blocked'
    };

    function normalizePriorityValue(p) {
      if (p == null) return null;
      if (typeof p === 'string' && /^\d+$/.test(p)) return Number(p);
      if (typeof p === 'string') return p.trim().toLowerCase();
      return p;
    }

    function priorityClass(p) {
      const key = normalizePriorityValue(p);
      return PRIORITY_MAP[key] || '';
    }
</script>

<style scoped>
    .statuses {
        grid-column: 1 / -1;
        width: 100%;
        max-width: 100%;
        min-width: 0;
        display: grid;
        gap: 8px;
    }

    .statuses-scrollbar {
        width: 100%;
        max-width: 100%;
        overflow-x: auto;
        overflow-y: hidden;
    }

    .statuses-scrollbar-content {
        width: max(100%, var(--statuses-min-width));
        height: 1px;
    }

    .statuses-viewport {
        width: 100%;
        max-width: 100%;
        min-width: 0;
        overflow-x: auto;
        scrollbar-width: none;
    }

    .statuses-viewport::-webkit-scrollbar {
        display: none;
    }

    .status-columns {
        display: flex;
        flex-direction: row;
        gap: 10px;
        width: max(100%, var(--statuses-min-width));
    }

    .status-column {
        flex: 1 0 200px;
        min-width: 0;
        display: flex;
        flex-direction: column;
    }

    .status-name {
        display: grid;
        place-items: center;
        font-size: 20px;
        font-weight: 700;
        color: black;
        min-height: 45px;
        border: 2px solid black;
        border-radius: 14px 14px 0 0;
        background-color: #fff;
    }

    .status-column-body {
        min-height: 500px;
        height: 100%;
        padding: 10px 5px;
        border: 2px solid black;
        border-top: none;
        border-radius: 0 0 14px 14px;
        background-color: #fafafa;
    }

    .card {
        margin-bottom: 10px;
        padding: 8px;
        border-radius: 8px;
        cursor: grab;
        background-color: white;
        border: 1px solid #e0e0e0;
        box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        transition: transform 0.15s ease, box-shadow 0.15s ease;
    }

    .card:active {
        cursor: grabbing;
    }

    .ghost-card {
        opacity: 0.4;
        background-color: #e2e8f0 !important;
        border: 2px dashed #cbd5e1 !important;
    }

    .card.priority-critical { border-left: 6px solid #b00020; background-color: #fff5f6; }
    .card.priority-high { border-left: 6px solid #ff6f00; background-color: #fff8f0; }
    .card.priority-medium { border-left: 6px solid #f59e0b; background-color: #fffaf0; }
    .card.priority-low { border-left: 6px solid #06b6d4; background-color: #f0f9fb; }
    .card.priority-backlog { border-left: 6px solid #6b7280; background-color: #f7f7f8; }
    .card.priority-blocked { border-left: 6px solid #6f42c1; background-color: #f7f2fb; }
</style>
