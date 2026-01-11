<template>
  <div class="filter-container">
    <button 
      @click="toggleFilter" 
      class="filter-button"
      :class="{ 'filter-button--active': isFilterOpen }"
    >
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M3 18C3 18.55 3.45 19 4 19H9V17H4C3.45 17 3 17.45 3 18ZM3 6C3 6.55 3.45 7 4 7H13V5H4C3.45 5 3 5.45 3 6ZM13 20V19H20C20.55 19 21 18.55 21 18C21 17.45 20.55 17 20 17H13V16C13 15.45 12.55 15 12 15C11.45 15 11 15.45 11 16V20C11 20.55 11.45 21 12 21C12.55 21 13 20.55 13 20ZM7 10V11H4C3.45 11 3 11.45 3 12C3 12.55 3.45 13 4 13H7V14C7 14.55 7.45 15 8 15C8.55 15 9 14.55 9 14V10C9 9.45 8.55 9 8 9C7.45 9 7 9.45 7 10ZM21 12C21 11.45 20.55 11 20 11H11V13H20C20.55 13 21 12.55 21 12ZM16 9C16.55 9 17 8.55 17 8V7H20C20.55 7 21 6.55 21 6C21 5.45 20.55 5 20 5H17V4C17 3.45 16.55 3 16 3C15.45 3 15 3.45 15 4V8C15 8.55 15.45 9 16 9Z"/>
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
  transition: background-color 0.3s ease, border-color 0.3s ease;
  color: #333;

  & svg {
    fill: #000;
    transition: fill 0.3s ease;
  }

  &:hover {
    background-color: #f5f5f5;
    border-color: #999;
  }

  &--active {
    background-color: #007bff;
    border-color: #007bff;
    color: white;

    & svg {
      fill: #fff;
    }

    &:hover {
      background-color: #0869d1;
      border-color: #007bff;
    }
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
  cursor: pointer;

  input[type="checkbox"] {
    margin-right: 8px;
    cursor: pointer;
    width: 20px;
    height: 20px;
  }

  span {
    font-size: 14px;
    line-height: 1;
  }

  &:not(:last-child) {
    margin-bottom: 8px;
  }
}
</style>
