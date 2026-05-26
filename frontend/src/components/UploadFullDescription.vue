<template>
  <div class="dropzone dropzone--full-description">
    <p class="dropzone__section-title">Полное описание модели (файл):</p>
    <div
      class="dropzone-container"
      @dragover="dragover"
      @dragleave="dragleave"
      @drop="drop"
    >
      <input
        ref="fileInput"
        type="file"
        class="hidden-input"
        :accept="accept"
        @change="onFileSelected"
      />

      <label
        v-if="!displayName"
        class="file-label"
        @click.prevent="openFileDialog"
      >
        <span v-if="isDragging">Отпустите, чтобы загрузить файл полного описания.</span>
        <span v-else>Перетащите файл (.txt, .doc, .docx, .pdf) или нажмите для выбора.</span>
      </label>

      <div v-else class="preview-card preview-card--full-description">
        <button
          class="ml-2"
          type="button"
          title="Удалить файл"
          @click="removeFile"
        >
          <b>&times;</b>
        </button>
        <p :title="displayName">{{ displayName }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import {
  FULL_DESCRIPTION_ACCEPT,
  isFullDescriptionFile,
} from '@/utils/fullDescriptionFormats';

export default {
  props: {
    existingFile: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      isDragging: false,
      accept: FULL_DESCRIPTION_ACCEPT,
      localFile: null,
      uploadedFileId: null,
      keepExisting: true,
    };
  },
  computed: {
    displayName() {
      if (this.localFile) {
        return this.localFile.name;
      }
      if (this.keepExisting && this.existingFile?.title) {
        return this.existingFile.title;
      }
      return '';
    },
    currentFileId() {
      if (this.uploadedFileId) {
        return this.uploadedFileId;
      }
      if (this.keepExisting && this.existingFile?.id) {
        return this.existingFile.id;
      }
      return null;
    },
  },
  methods: {
    emitChange() {
      this.$emit('onChange', {
        full_description_file_id: this.currentFileId,
      });
    },
    openFileDialog() {
      this.$refs.fileInput.click();
    },
    dragover(event) {
      event.preventDefault();
      this.isDragging = true;
    },
    dragleave() {
      this.isDragging = false;
    },
    drop(event) {
      event.preventDefault();
      this.isDragging = false;
      const file = event.dataTransfer.files?.[0];
      if (file) {
        this.handleSelectedFile(file);
      }
    },
    onFileSelected() {
      const file = this.$refs.fileInput.files?.[0];
      if (file) {
        this.handleSelectedFile(file);
      }
      this.$refs.fileInput.value = '';
    },
    handleSelectedFile(file) {
      if (!isFullDescriptionFile(file.name)) {
        console.log('Неподдерживаемый формат файла полного описания');
        return;
      }

      this.localFile = file;
      this.keepExisting = false;
      this.uploadWithAuth(file);
    },
    async uploadWithAuth(file) {
      const accessToken = $cookies.get('access_token');
      const refreshToken = $cookies.get('refresh_token');

      const upload = async (access) => {
        const fileData = new FormData();
        fileData.append('file', file);
        const response = await axios.post('/api/file', fileData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            Authorization: `Bearer ${access}`,
          },
        });
        this.uploadedFileId = response.data.id;
        this.emitChange();
      };

      try {
        if (!accessToken && refreshToken) {
          const response = await axios.post('/api/users/refresh', {
            refresh_token: refreshToken,
          });
          const access = response.data.tokens.access_token;
          $cookies.set('access_token', access, '15min', '/');
          $cookies.set('refresh_token', response.data.tokens.refresh_token, '7d', '/');
          localStorage.setItem('isAuth', true);
          await upload(access);
          return;
        }

        await upload(accessToken);
      } catch (error) {
        console.log(error);
        this.localFile = null;
        this.uploadedFileId = null;
        this.emitChange();
      }
    },
    async removeFile() {
      const fileId = this.currentFileId;

      this.localFile = null;
      this.uploadedFileId = null;
      this.keepExisting = false;

      if (fileId) {
        await this.deleteWithAuth(fileId);
      }

      this.$emit('onChange', { full_description_file_id: null });
    },
    async deleteWithAuth(fileId) {
      const accessToken = $cookies.get('access_token');
      const refreshToken = $cookies.get('refresh_token');

      const remove = async (access) => {
        await axios.delete(`/api/file/${fileId}`, {
          headers: {
            Authorization: `Bearer ${access}`,
          },
        });
      };

      try {
        if (!accessToken && refreshToken) {
          const response = await axios.post('/api/users/refresh', {
            refresh_token: refreshToken,
          });
          const access = response.data.tokens.access_token;
          $cookies.set('access_token', access, '15min', '/');
          $cookies.set('refresh_token', response.data.tokens.refresh_token, '7d', '/');
          await remove(access);
          return;
        }

        await remove(accessToken);
      } catch (error) {
        console.log(error);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.dropzone--full-description {
  margin-bottom: 1.5rem;
}

.dropzone__section-title {
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.dropzone-container {
  padding: 2rem;
  background: #f7fafc;
  border: 1px solid #e2e8f0;
}

.hidden-input {
  opacity: 0;
  overflow: hidden;
  position: absolute;
  width: 1px;
  height: 1px;
}

.file-label {
  font-size: 16px;
  display: block;
  cursor: pointer;
  text-align: center;
}

.preview-card--full-description {
  display: flex;
  flex-direction: column;
  border: 1px solid #a2a2a2;
  padding: 5px;
  max-width: 100%;

  p {
    word-break: break-word;
    margin: 0;
  }

  .ml-2 {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    margin-left: auto;
    margin-bottom: 0.5rem;
  }
}
</style>
