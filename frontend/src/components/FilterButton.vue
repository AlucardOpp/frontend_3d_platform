<template>
  <div class="filter-container">
    <button 
      @click="toggleFilter" 
      class="filter-button"
      :class="{ 'filter-button--active': isFilterOpen }"
    >
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M3.33334 5H16.6667M5.00001 10H15M6.66667 15H13.3333" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </button>
    <div v-if="isFilterOpen" class="filter-panel">
      <p class="filter-panel__title">Ключевые слова</p>
      <label class="filter-checkbox">
        <input 
          type="checkbox" 
          v-model="selectedKeywords" 
          value="история"
          @change="applyFilter"
        />
        <span>история</span>
      </label>
      <label class="filter-checkbox">
        <input 
          type="checkbox" 
          v-model="selectedKeywords" 
          value="архитектура"
          @change="applyFilter"
        />
        <span>архитектура</span>
      </label>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FilterButton',
  data() {
    return {
      isFilterOpen: false,
      selectedKeywords: []
    }
  },
  methods: {
    toggleFilter() {
      this.isFilterOpen = !this.isFilterOpen;
    },
    applyFilter() {
      this.$emit('filter-changed', this.selectedKeywords);
    }
  }
}
</script>

<style lang="scss" scoped>
.filter-container {
  position: relative;
  display: inline-block;
  margin-bottom: 20px;
}

.filter-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #333;

  &:hover {
    background: #f5f5f5;
    border-color: #999;
  }

  &--active {
    background: #007bff;
    border-color: #007bff;
    color: white;
  }
}

.filter-panel {
  position: absolute;
  top: 45px;
  left: 0;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 15px;
  min-width: 200px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.filter-panel__title {
  margin: 0 0 10px 0;
  font-weight: 600;
  font-size: 14px;
}

.filter-checkbox {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  cursor: pointer;

  input[type="checkbox"] {
    margin-right: 8px;
    cursor: pointer;
  }

  span {
    font-size: 14px;
  }
}
</style>
