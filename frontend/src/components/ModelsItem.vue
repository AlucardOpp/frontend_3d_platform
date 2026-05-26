<template>
  <div class="model" :data-model-id="`${model.model.id}`">
    <div class="model__content">
      <h2 class="model__title">{{ model.model.title }}</h2>
      <div class="model__swipers">
        <swiper
          :slides-per-view="1" 
          :space-between="20"
          :allowTouchMove="false"
          :modules="[Thumbs]"
          :thumbs="{ swiper: thumbsSwiper }"
          class="model__default-swiper"
        >
          <swiper-slide v-for="model in model.files.glb">
            <model-viewer auto-rotate class="model__model" :src="getFileUrl(model.url)" camera-controls=""></model-viewer>
          </swiper-slide>
          <swiper-slide v-for="img in model.files.img">
            <img :src="getFileUrl(img.url)">
          </swiper-slide>
          <swiper-slide v-for="video in model.files.video">
            <video controls :src="getFileUrl(video.url)"></video>
          </swiper-slide>
        </swiper>
        <swiper v-if="model.files.glb || model.files.img || model.files.other || model.files.video" class="model__thumbs-swiper"
          @swiper="setThumbsSwiper"
          :slides-per-view="2"
          :space-between="10"
          :freeMode="true"
          :watchSlidesProgress="true"
          :grabCursor="true"
          :modules="[Thumbs, FreeMode, Navigation]"
          :mousewheel="true"
          :navigation="true"
          :breakpoints="{
            '768': {
              slidesPerView: 4,
            },
          }"
        >
          <swiper-slide v-for="model in model.files.glb">
            <model-viewer auto-rotate class="model__model" :src="getFileUrl(model.url)"></model-viewer>
          </swiper-slide>
          <swiper-slide v-for="img in model.files.img">
            <img :src="getFileUrl(img.url)">
          </swiper-slide>
          <swiper-slide v-for="video in model.files.video">
            <video class="model__video" :src="getFileUrl(video.url)"></video>
          </swiper-slide>
        </swiper>
      </div>
      <span class="model__author">Создано пользователем {{ model.author.name }}</span>
      <span v-if="model.model.keywords && model.model.keywords.trim()" class="model__keywords">Ключевые слова: {{ formatKeywords(model.model.keywords) }}</span>
      <span class="model__data">Загружено {{ toDateString(model.model.created_at) }}</span>
    </div>
    <div class="model__panel">
      <div class="model__info">
        <h3 class="model__sub-title">Описание</h3>
        <div class="model__info-container">
          <p class="model__description">{{ model.model.description }}</p>
          <a
            v-if="model.full_description"
            class="model__full-description-link"
            :href="getFileUrl(model.full_description.url)"
            :download="model.full_description.title"
          >
            {{ model.full_description.title }}
          </a>
          <ul v-if="isAuth && (this.isAdmin || this.userId == model.author.id)" class="model__actions"> 
            <li class="model__action">
              <button @click="$router.push(`/${model.model.id}`)" class="model__action-btn btn btn--white">Редактировать</button>
            </li>
            <li class="model__action">
              <button @click="$emit('remove', model)" class="model__action-btn btn">Удалить</button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { Swiper, SwiperSlide } from 'swiper/vue';
import 'swiper/css';
import "swiper/css/free-mode"
import "swiper/css/thumbs"
import "swiper/css/navigation"
import {FreeMode, Thumbs, Navigation} from 'swiper';

export default {
  components: {
    Swiper,
    SwiperSlide,
  },
  methods: {
    toDateString(date) {
      return new Date(date).toLocaleDateString('ru-RU', { year: 'numeric', month: 'numeric', day: 'numeric' })
    },
    formatKeywords(keywords) {
      if (!keywords) return '';
      // Преобразуем ключевые слова в формат #ключевое_слово
      return keywords.split(' ').filter(k => k.trim()).map(k => {
        const trimmed = k.trim();
        return trimmed.startsWith('#') ? trimmed : '#' + trimmed;
      }).join(' ');
    },
    getOrigin() {
      return window.location.origin;
    },
    getFileUrl(fileUrl) {
      if (!fileUrl) {
        return '';
      }
      
      if (fileUrl.startsWith('/upload/')) {
        return this.getOrigin() + fileUrl;
      }
      
      if (fileUrl.startsWith('http://') || fileUrl.startsWith('https://')) {
        return fileUrl;
      }
      
      if (fileUrl.includes('upload/')) {
        const fileName = fileUrl.split('upload/').pop();
        return this.getOrigin() + '/upload/' + fileName;
      }
      
      if (fileUrl.includes('/')) {
        const fileName = fileUrl.split('/').pop();
        return this.getOrigin() + '/upload/' + fileName;
      }
      
      return this.getOrigin() + '/upload/' + fileUrl;
    }
  },
  setup() {
    const thumbsSwiper = ref(null);
    const setThumbsSwiper = (swiper) => {
        thumbsSwiper.value = swiper;
      };

    return {
      FreeMode,
      Thumbs,
      Navigation,
      thumbsSwiper,
      setThumbsSwiper,
    };
  },
  props: {
    model: {
      type: Object,
      required: true,
    },
    isAdmin: {
      type: Boolean,
      required: true,
    },
    userId: {
      type: Number,
      required: true
    }
  },
  computed: {
    isAuth() {
      return localStorage.isAuth === 'true';
    }
  },
  mounted() {
    const currentModel = document.querySelector(`[data-model-id="${this.model.model.id}"`);
    const videos = currentModel.querySelectorAll('.model__video');
    videos.forEach((video) => {
      video.currentTime = video.currentTime + 1;
    });
  },
  name: 'models-item',
}
</script>

<style lang="scss" scoped>
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
</style>