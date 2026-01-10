import axios from "axios";
import VueCookies from 'vue-cookies';

export const uploadModule = {
    state: () => ({
        name: '',
        description: '',
        keywords: ''
    }),
    mutations: {
      setName(state, name) {
        state.name = name;
      },
      setDescription(state, description) {
        state.description = description;
      },
      setKeywords(state, keywords) {
        state.keywords = keywords;
      },
    },
    getters: { 
    },
    actions: {
      handleSubmitUpload({state, commit}) {
        axios.post('/api/model', {
          name: state.email,
          description: state.password
        })
        .then(response => {
        })
        .catch(error => {
          console.loog(error);
        });
      },
    },
    namespaced: true
}