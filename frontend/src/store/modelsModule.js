import axios from "axios";
import VueCookies from 'vue-cookies';

export const modelsModule = {
    state: () => ({
        models: [],
        userId: null,
        isModelsLoading: false,
        page: 1,
        limit: 1,
        selectedKeywords: [],
    }),
    mutations: {
        setModels(state, models) {
            state.models = models;
        },
        setLoading(state, bool) {
            state.isModelsLoading = bool;
        },
        setUserId(state, userId) {
            state.userId = userId;
        },
        setPage(state, page) { 
            state.page = page;
        },
        setSelectedKeywords(state, keywords) {
            state.selectedKeywords = keywords;
        },
    },
    actions: {
        async fetchModels({state, commit}) {
            try {
                commit('setLoading', true);
                const params = {
                    _page: state.page,
                    _limit: state.limit
                };
                if (state.userId !== null) {
                    params.user_id = state.userId;
                }
                if (state.selectedKeywords && state.selectedKeywords.length > 0) {
                    params.keywords = state.selectedKeywords.join(',');
                }
                const response = await axios.get('/api/model', { params });
                const models = Array.isArray(response.data) ? response.data : [];
                commit('setModels', models)
            } catch (e) {
                console.log(e);
            } finally {
                commit('setLoading', false);
            }
        },
        async loadMoreModels({state, commit}) {
            try {
                commit('setPage', state.page + 1)
                const params = {
                    _page: state.page,
                    _limit: state.limit
                };
                if (state.userId !== null) {
                    params.user_id = state.userId;
                }
                if (state.selectedKeywords && state.selectedKeywords.length > 0) {
                    params.keywords = state.selectedKeywords.join(',');
                }
                const response = await axios.get('/api/model', { params });
                const newModels = Array.isArray(response.data) ? response.data : [];
                commit('setModels', [...state.models, ...newModels]);
            } catch (e) {
                console.log(e)
            }
        },
        setPagesToOne({state, commit}) {
            commit('setPage', 1);
        }
    },
    namespaced: true
}