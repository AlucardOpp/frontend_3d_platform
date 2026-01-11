<template>
  <section class="models">
      <div class="models__controls">
        <filter-button @filter-changed="handleFilterChange" />
        <div class="sort-select-wrapper">
          <select 
            v-model="sortType" 
            @change="handleSortChange"
            class="sort-select"
          >
            <option value="newest">Сначала новые</option>
            <option value="oldest">Сначала старые</option>
            <option value="alphabetical">С начала алфавита</option>
            <option value="alphabetical-reverse">С конца алфавита</option>
          </select>
        </div>
      </div>
      <models-list
        :models="sortedModels"
        @remove="removeModel"
        v-if="!isModelsLoading"
      />
      <div v-else>Идет загрузка...</div>
      <div v-intersection="loadMoreModels" class="observer"></div>
  </section>
</template>

<script>
import {mapState, mapGetters, mapActions, mapMutations} from 'vuex';
import axios from "axios";
import FilterButton from "@/components/FilterButton.vue";

export default {
  components: {
    FilterButton
  },
  data() {
    return {
      sortType: 'newest'
    }
  },
  methods: {
    ...mapMutations({
      setPage: 'models/setPage',
    }),
    ...mapActions({
      loadMoreModels: 'models/loadMoreModels',
      fetchModels: 'models/fetchModels',
      setPagesToOne: 'models/setPagesToOne',
    }),
    handleFilterChange(keywords) {
      this.$store.commit('models/setSelectedKeywords', keywords);
      this.setPagesToOne();
      this.fetchModels();
    },
    handleSortChange() {
    },
    removeModel(model){
      const accessToken = $cookies.get("access_token");
      const modelId = model.model.id;
      axios.delete(`/api/model/${modelId}`,
      { 
        data: { 
          id: modelId
        }, 
        headers: { 
          "Authorization": `Bearer ${accessToken}`
        } 
      })
      .then(response => {
        const model = document.querySelector(`[data-model-id="${modelId}"`);
        model.remove();
        const models = document.querySelectorAll('.model');
        if(models.length <= 1) {
          this.loadMoreModels();
        }
      })
      .catch(error => {
        console.log(error);
      });
    }
  },
  mounted() {
    this.setPagesToOne();
    this.fetchModels();
    let modelViewerScript = document.createElement('script');
    modelViewerScript.setAttribute('src', 'https://unpkg.com/@google/model-viewer/dist/model-viewer.min.js');
    modelViewerScript.setAttribute('type', 'module');
    document.head.appendChild(modelViewerScript);
  },
  computed: {
    ...mapState({
      models: state => state.models.models,
      isModelsLoading: state => state.models.isModelsLoading,
      page: state => state.models.page,
      limit: state => state.models.limit,
      totalPages: state => state.models.totalPages,
    }),
    filteredModels() {
      return this.models;
    },
    sortedModels() {
      const models = [...this.filteredModels];
      
      switch (this.sortType) {
        case 'newest':
          return models.sort((a, b) => {
            const dateA = new Date(a.model.created_at);
            const dateB = new Date(b.model.created_at);
            return dateB - dateA;
          });
        case 'oldest':
          return models.sort((a, b) => {
            const dateA = new Date(a.model.created_at);
            const dateB = new Date(b.model.created_at);
            return dateA - dateB;
          });
        case 'alphabetical':
          return models.sort((a, b) => {
            const titleA = (a.model.title || '').toLowerCase();
            const titleB = (b.model.title || '').toLowerCase();
            return titleA.localeCompare(titleB, 'ru');
          });
        case 'alphabetical-reverse':
          return models.sort((a, b) => {
            const titleA = (a.model.title || '').toLowerCase();
            const titleB = (b.model.title || '').toLowerCase();
            return titleB.localeCompare(titleA, 'ru');
          });
        default:
          return models;
      }
    }
  },
}
</script>

<style lang="scss" scoped>
.models__controls {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.sort-select-wrapper {
  position: relative;
  display: inline-block;
  
  &::after {
    content: '';
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%) rotate(0deg);
    width: 0;
    height: 0;
    border-left: 5px solid transparent;
    border-right: 5px solid transparent;
    border-top: 6px solid #333;
    pointer-events: none;
    transition: transform 0.3s ease;
    z-index: 1;
  }
  
  &:focus-within::after {
    transform: translateY(-50%) rotate(180deg);
  }
}

.sort-select {
  padding: 8px 35px 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  font-size: 14px;
  cursor: pointer;
  transition: border-color 0.3s ease, background-color 0.3s ease;
  min-width: 180px;
  width: 100%;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  min-height: 40px;
  
  &::-ms-expand {
    display: none;
  }
  
  &:hover {
    border-color: #999;
    background-color: #f5f5f5;
  }
  
  &:focus {
    outline: none;
    border-color: #007bff;
    background-color: white;
  }
}

.observer {
  height: 30px;
}
</style>
