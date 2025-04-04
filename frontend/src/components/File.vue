<template>
  <ul class="file-list"> <!-- 添加 ul 标签包裹列表项 -->
    <li v-for="(file, index) in files" :key="index" class="file-item">
      <span class="file-icon">📄</span>
      <span class="file-name">{{ file.name }}</span>
      <span class="file-size">{{ file.size }}</span>
      <!-- 添加下载按钮 -->
      <button @click="emitDownload(file)" class="download-button">下载</button>
    </li>
  </ul>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'; // 引入 defineEmits

const props = defineProps({
  files: {
    type: Array,
    required: true,
    default: () => [] // 提供一个默认空数组
  }
});

// 定义可以触发的事件
const emit = defineEmits(['download']);

// 触发下载事件的方法
const emitDownload = (file) => {
  emit('download', file);
};
</script>

<style scoped>
.file-list-container {
  font-family: 'Arial', sans-serif;
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  max-width: 500px; /* 限制最大宽度 */
  margin: 20px auto; /* 居中显示 */
}

.file-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.file-item {
  display: flex;
  align-items: center;
  background-color: #fff;
  padding: 10px 15px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: 1px solid #eee;
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
}

.file-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.file-icon {
  margin-right: 10px;
  font-size: 1.2em;
}

.file-name {
  flex-grow: 1; /* 文件名占据剩余空间 */
  color: #555;
  margin-right: 10px; /* 与大小保持距离 */
}

.file-size {
  color: #888;
  font-size: 0.9em;
  white-space: nowrap; /* 防止大小换行 */
  margin-right: 10px; /* 与下载按钮保持距离 */
}

.download-button {
  background-color: #007bff; /* 蓝色背景 */
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8em;
  transition: background-color 0.2s;
  white-space: nowrap; /* 防止按钮文字换行 */
}

.download-button:hover {
  background-color: #0056b3; /* 悬停时深一点的蓝色 */
}
</style>