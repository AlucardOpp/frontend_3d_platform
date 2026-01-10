<template>
  <section class="models">
      <filter-button @filter-changed="handleFilterChange" />
      <models-list
        :models="filteredModels"
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
    }
  },
}
</script>

<style lang="scss" scoped>
.observer {
  height: 30px;
}
</style>
