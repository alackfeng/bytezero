<template>
  <div>
    <form id="uploadForm" enctype=" multipart/form-data" method="post" action="/api/v1/bridge/files/upload">
      <div v-for="(item, index) in uploadFiles" :key="index" class="row">
        <div style="display: flex">
          <input id="uploadFile" type="file" :name="index" @change="fileSelected" />
          <div>{{ item.size }}</div>
          <div> - {{ item.type }}</div>
          <div> - {{ item.lastModified }}</div>
          <div> - {{ item.status }}</div>
        </div>
      </div>
      <div>
        <van-button type='primary' size="small" @click="fileStartUpload">开始上传</van-button>
        <van-button type='second' size="small" @click="fileEndUpload">暂停上传</van-button>
      </div>
    </form>
  </div>
</template>

<script>
import Vue from 'vue';
import { formatTs } from '~/plugins/utils';

export default Vue.extend({
  name: "NetCardUpload",
  data() {
    return {
      uploadFiles: [
        { size: "0", name: "", type: "unknown", lastModified: "0", status: '' },
        { size: "0", name: "", type: "unknown", lastModified: "0", status: '' },
        { size: "0", name: "", type: "unknown", lastModified: "0", status: '' },
        { size: "0", name: "", type: "unknown", lastModified: "0", status: '' },
        { size: "0", name: "", type: "unknown", lastModified: "0", status: '' }
      ]
    }
  },
  methods: {
    adjustSize(size) {
      let fileSize = '0 KB'
      if (size > 1024 * 1024)
        fileSize = (Math.round(size * 100 / (1024 * 1024)) / 100).toString() + 'MB';
      else
        fileSize = (Math.round(size * 100 / 1024) / 100).toString() + 'KB';
      return fileSize
    },
    fileSelected(event) {
      const file = event.target.files[0]
      const key = parseInt(event.target.name)

      this.uploadFiles[key].key = key
      this.uploadFiles[key].size = this.adjustSize(file.size)
      this.uploadFiles[key].name = file.name
      this.uploadFiles[key].type = file.type
      this.uploadFiles[key].lastModified = formatTs(file.lastModified)
      this.uploadFiles[key].file = file
    },
    fileStartUpload(event) {
      event.preventDefault()

      this.uploadFiles.forEach(item => {
        if (!item.file) {
          return
        }
        // const key = item.key
        // this.$api.uploadFile(item).then(res => {
        //   console.log('-------------- fileStartUpload ', res)
        //   this.uploadFiles[key].status = res.code
        // })

        const formData = new FormData()
        formData.append('total', item.size)
        formData.append('files', item.file)

        const config = {
          headers: {
            'Content-Type': 'multipart/form-data'
          },
          onUploadProgress: (progressEvent) => {
            const progressPrecent = (progressEvent.loaded / progressEvent.total * 100);
            this.uploadFiles[key].status = progressPrecent
          }
        }
        const key = item.key
        this.$axios.post('http://192.168.90.23:7790/api/v1/bridge/files/upload', formData, config).then(res => {
          console.log('>>>>>fileStartUpload ', res.data)
          // this.uploadFiles[key].status = res.data.info
        })
      })
    },
    fileEndUpload() {

    }
  }
})
</script>
