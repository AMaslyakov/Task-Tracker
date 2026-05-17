<template>
    <!-- Информационный блок о выбранной команде (Растянут и масштабирован) -->
    <div v-if="props.command?.name" class="team-info-widget middle-panel scaled-panel">

        <!-- Верхняя секция с текстом и увеличенным котиком -->
        <div class="team-header-grid">
            <div class="team-text-details">
                <!-- Название команды с глубоким и строгим градиентом -->
                <h2 class="team-display-name-middle">{{ props.command.name }}</h2>
                <!-- Сбалансированный мягкий фон описания -->
                <p v-if="props.command.description" class="team-display-desc-middle-balanced">
                    {{ props.command.description }}
                </p>
            </div>

            <!-- Котик увеличен в 1.3 раза и зафиксирован на правой позиции -->
            <div class="cat-container">
                <svg class="coder-cat" viewBox="0 0 120 100" xmlns="http://w3.org">
                    <!-- Ушки котика -->
                    <polygon points="25,40 15,10 40,25" fill="#475569" />
                    <polygon points="28,38 20,18 38,26" fill="#f43f5e" />
                    <polygon points="75,40 85,10 60,25" fill="#475569" />
                    <polygon points="72,38 80,18 62,26" fill="#f43f5e" />
                    <!-- Голова -->
                    <ellipse cx="50" cy="45" rx="35" ry="25" fill="#64748b" />
                    <!-- Глазки -->
                    <path d="M32,42 Q40,38 42,45" stroke="#0f172a" stroke-width="2.5" fill="none" stroke-linecap="round" />
                    <path d="M68,42 Q60,38 58,45" stroke="#0f172a" stroke-width="2.5" fill="none" stroke-linecap="round" />
                    <!-- Розовый носик и ротик -->
                    <polygon points="48,48 52,48 50,51" fill="#f43f5e" />
                    <path d="M46,54 Q50,57 54,54" stroke="#0f172a" stroke-width="1.5" fill="none" />
                    <!-- Усики -->
                    <line x1="10" y1="45" x2="25" y2="47" stroke="#cbd5e1" stroke-width="1.5" />
                    <line x1="10" y1="52" x2="25" y2="51" stroke="#cbd5e1" stroke-width="1.5" />
                    <line x1="90" y1="45" x2="75" y2="47" stroke="#cbd5e1" stroke-width="1.5" />
                    <line x1="90" y1="52" x2="75" y2="51" stroke="#cbd5e1" stroke-width="1.5" />

                    <!-- Тело и лапки -->
                    <path d="M25,65 Q50,55 75,65 L80,95 L20,95 Z" fill="#475569" />
                    <ellipse cx="38" cy="78" rx="7" ry="5" fill="#64748b" />
                    <ellipse cx="62" cy="78" rx="7" ry="5" fill="#64748b" />

                    <!-- Ноутбук -->
                    <rect x="25" y="65" width="50" height="12" rx="2" fill="#0f172a" stroke="#a855f7" stroke-width="1" />
                    <!-- Код на экране -->
                    <line x1="30" y1="69" x2="45" y2="69" stroke="#10b981" stroke-width="1.5" />
                    <line x1="30" y1="73" x2="60" y2="73" stroke="#6366f1" stroke-width="1.5" />
                    <line x1="50" y1="69" x2="70" y2="69" stroke="#f43f5e" stroke-width="1.5" />
                    <polygon points="20,78 100,78 90,92 30,92" fill="#334155" stroke="#a855f7" stroke-width="1" />
                    <circle cx="50" cy="85" r="2" fill="#10b981" class="blink-led" />
                </svg>
            </div>
        </div>

        <!-- Список участников -->
        <div v-if="props.command.members && props.command.members.length" class="team-roster-sub">
            <span class="roster-meta-label">Команда:</span>
            <div class="roster-chips-container">
                <!-- Светло-серые бейджи участников -->
                <div
                    v-for="memberName in props.command.members"
                    :key="memberName"
                    class="member-chip-pill-middle-light"
                >
                    <span class="activity-dot-indicator-middle-light"></span>
                    <span class="member-chip-name">{{ memberName }}</span>
                </div>
            </div>
        </div>
    </div>

    <!-- Исходный неизмененный контейнер доски с колонками -->
    <div class="statuses" :style="{ '--statuses-min-width': minScrollWidth }">
        <div class="statuses-scrollbar" ref="topScrollRef" @scroll="syncColumnsScroll">
            <div class="statuses-scrollbar-content"></div>
        </div>

        <div class="statuses-viewport" ref="columnsScrollRef" @scroll="syncTopScroll">
            <div class="status-columns">
                <div class="status-column" v-for="status_name in config.statuses" :key="status_name">
                    <div class="status-name"><span>{{ status_name }}</span></div>

                    <Sortable
                        :list="tasksByStatus[status_name] || []"
                        itemKey="id"
                        tag="div"
                        class="status-column-body"
                        :options="{ group: 'tasks-kanban', animation: 200, ghostClass: 'ghost-card' }"
                        @change="(event) => handleTaskMove(event, status_name)"
                    >
                        <template #item="{ element: task }">
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
        tasks: { type: Array, required: true },
        command: { type: Object, required: true }
    });

    const emit = defineEmits(['update:tasks', 'task-status-changed']);

    const topScrollRef = ref(null);
    const columnsScrollRef = ref(null);
    const columnMinWidth = 200;
    const columnGap = 10;
    const config = computed(() => props.command.config_dashboard);

    const minScrollWidth = computed(() => {
        const statusesCount = config.value?.statuses?.length || 0;
        const gapsWidth = Math.max(statusesCount - 1, 0) * columnGap;
        return `${statusesCount * columnMinWidth + gapsWidth}px`;
    });

    function getTaskWeight(task) {
        const pValue = task.priority_id || task.priority;
        if (pValue == null) return 999;
        const pStr = String(pValue).trim().toLowerCase();
        const WEIGHT_MAP = {
            '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6,
            'critical': 1, 'high': 2, 'medium': 3, 'low': 4, 'backlog': 5, 'blocked': 6,
            'crit': 1, 'urgent': 1, 'important': 2, 'major': 2, 'normal': 3, 'standard': 3,
            'minor': 4, 'trivial': 4, 'later': 5, 'someday': 5, 'plan': 5, 'block': 6
        };
        return WEIGHT_MAP[pStr] !== undefined ? WEIGHT_MAP[pStr] : 999;
    }

    const tasksByStatus = computed(() => {
        const map = {};
        if (config.value?.statuses) {
            config.value.statuses.forEach(status => {
                const filtered = props.tasks.filter(task => {
                    const currentStatus = typeof task.status === 'object' ? task.status?.name : task.status;
                    return String(currentStatus).trim().toLowerCase() === String(status).trim().toLowerCase();
                });
                map[status] = filtered.sort((a, b) => getTaskWeight(a) - getTaskWeight(b));
            });
        }
        return map;
    });

    function handleTaskMove(event, newStatus) {
        if (event.added) {
            const movedTask = event.added.element;
            const updatedTasks = props.tasks.map(task => {
                if (task.id === movedTask.id) return { ...task, status: newStatus };
                return task;
            });
            emit('update:tasks', updatedTasks);
            emit('task-status-changed', { taskId: movedTask.id, newStatus });
        }
    }

    function syncColumnsScroll() {
        if (!topScrollRef.value || !columnsScrollRef.value) return;
        columnsScrollRef.value.scrollLeft = topScrollRef.value.scrollLeft;
    }

    function syncTopScroll() {
        if (!topScrollRef.value || !columnsScrollRef.value) return;
        topScrollRef.value.scrollLeft = columnsScrollRef.value.scrollLeft;
    }

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
    /* Средне-серая подложка */
    .team-info-widget.middle-panel {
        font-family: system-ui, -apple-system, sans-serif;
        background-color: #e2e8f0;
        border: 1px solid #cbd5e1;
        border-radius: 14px;
        padding: 20px 24px; /* Оставляем оригинальные отступы для совпадения с Toolbar */
        box-shadow: inset 0 1px 2px rgba(255, 255, 255, 0.4), 0 4px 12px -2px rgba(0, 0, 0, 0.03);
        margin-bottom: 28px;
        width: 100%;
        max-width: 100%;
        margin-left: 0;
        margin-right: 0;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        gap: 18px; /* Немного увеличен зазор между блоками внутри */
    }

    /* Сеточная верстка для текста и котика */
    .team-header-grid {
        display: grid;
        grid-template-columns: 1fr auto;
        align-items: center;
        gap: 32px; /* Увеличен зазор до котика */
        width: 100%;
    }

    .team-text-details {
        display: flex;
        flex-direction: column;
        max-width: 85%; /* ИСПРАВЛЕНО: Текстовый блок стал шире и просторнее внутри панели */
    }

    /* Название команды */
    .team-display-name-middle {
        font-size: 26px; /* ИСПРАВЛЕНО: Шрифт увеличен пропорционально масштабу */
        font-weight: 800;
        margin: 0 0 6px 0;
        letter-spacing: -0.02em;
        background: linear-gradient(135deg, #0f172a 0%, #1e3a8a 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
    }

    /* Сбалансированный мягкий фон описания */
    .team-display-desc-middle-balanced {
        font-size: 14.5px; /* ИСПРАВЛЕНО: Текст стал чуть крупнее для удобства чтения */
        color: #334155;
        margin: 0;
        line-height: 1.6;
        background-color: #f8fafc;
        padding: 14px 18px; /* Увеличены внутренние поля плашки описания */
        border-radius: 10px;
        border: 1px solid #e2e8f0;
        border-left: 4px solid #6366f1;
    }

    /* Контейнер котика */
    .cat-container {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }

    /* ИСПРАВЛЕНО: Размер котика увеличен ровно в 1.3 раза (со 100px до 130px) */
    .coder-cat {
        width: 130px;
        height: auto;
    }

    .blink-led {
        animation: led-pulse 1.5s infinite;
    }

    @keyframes led-pulse {
        0%, 100% { opacity: 0.3; }
        50% { opacity: 1; }
    }

    /* Список участников */
    .team-roster-sub {
        border-top: 1px solid #cbd5e1;
        padding-top: 14px;
        display: flex;
        align-items: center;
        gap: 14px;
        flex-wrap: wrap;
    }

    .roster-meta-label {
        font-size: 11px;
        font-weight: 700;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.08em;
    }

    .roster-chips-container {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
    }

    /* Светло-серые бейджи участников */
    .member-chip-pill-middle-light {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        background-color: #f1f5f9;
        border: 1px solid #cbd5e1;
        border-radius: 20px;
        padding: 5px 14px;
        font-size: 13px;
        font-weight: 600;
        color: #1e293b;
        transition: border-color 0.15s ease, background-color 0.15s ease;
    }

    .member-chip-pill-middle-light:hover {
        background-color: #ffffff;
        border-color: #6366f1;
    }

    /* Синий маркер участника */
    .activity-dot-indicator-middle-light {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background-color: #2563eb;
        box-shadow: 0 0 6px rgba(37, 99, 235, 0.4);
        display: inline-block;
    }

    /* Стили доски */
    .statuses { grid-column: 1 / -1; width: 100%; max-width: 100%; min-width: 0; display: grid; gap: 8px; }
    .statuses-scrollbar { width: 100%; max-width: 100%; overflow-x: auto; overflow-y: hidden; }
    .statuses-scrollbar-content { width: max(100%, var(--statuses-min-width)); height: 1px; }
    .statuses-viewport { width: 100%; max-width: 100%; min-width: 0; overflow-x: auto; scrollbar-width: none; }
    .statuses-viewport::-webkit-scrollbar { display: none; }
    .status-columns { display: flex; flex-direction: row; gap: 10px; width: max(100%, var(--statuses-min-width)); }
    .status-column { flex: 1 0 200px; min-width: 0; display: flex; flex-direction: column; }
    .status-name { display: grid; place-items: center; font-size: 20px; font-weight: 700; color: #334155; min-height: 45px; border: 2px solid black; border-radius: 14px 14px 0 0; background-color: #cbd5e1; }
    .status-column-body { min-height: 500px; height: 100%; padding: 10px 5px; border: 2px solid black; border-top: none; border-radius: 0 0 14px 14px; background-color: #fafafa; }
    .card { margin-bottom: 10px; padding: 8px; border-radius: 8px; cursor: grab; background-color: white; border: 1px solid #e0e0e0; box-shadow: 0 2px 4px rgba(0,0,0,0.05); transition: transform 0.15s ease, box-shadow 0.15s ease; }
    .card:active { cursor: grabbing; }
    .ghost-card { opacity: 0.4; background-color: #e2e8f0 !important; border: 2px dashed #cbd5e1 !important; }
    .card.priority-critical { border-left: 6px solid #b00020; background-color: #fff5f6; }
    .card.priority-high { border-left: 6px solid #ff6f00; background-color: #fff8f0; }
    .card.priority-medium { border-left: 6px solid #f59e0b; background-color: #fffaf0; }
    .card.priority-low { border-left: 6px solid #06b6d4; background-color: #f0f9fb; }
    .card.priority-backlog { border-left: 6px solid #6b7280; background-color: #f7f7f8; }
    .card.priority-blocked { border-left: 6px solid #6f42c1; background-color: #f7f2fb; }
</style>
