<template>
  <div class="dashboard-container">
    <!-- 操作按钮区域 -->
    <div class="action-buttons">
      <button @click="handleUpload" class="upload-button">上传文件</button>
      <!-- 这里可以添加更多按钮 -->
    </div>

    <!-- 文件列表 -->
    <File :files="sampleFiles" @download="handleDownload" />
  </div>
</template>

<script setup>
import { onBeforeMount, ref } from "vue";
import File from "../components/File.vue"; // 引入 File 组件
import { UploadFile } from "../../wailsjs/go/main/FileManager"; // 导入 Wails Runtime 的文件对话框功能
import { GetAllFilesInfo } from "../../wailsjs/go/main/Request"; // 导入 Wails Runtime 的文件对话框功能

// 示例文件数据
const sampleFiles = ref([]);
onBeforeMount(() => {
  GetAllFilesInfo()
    .then((res) => {
      console.log("res ===>" + res);
      
      // 将 res 解析为 JavaScript 对象
      const response = JSON.parse(res);
      
      // 获取 data 的值
      const data = response.data;
      
      // 将 data 赋值给 sampleFiles.value
      sampleFiles.value = data;
    })
    .catch((err) => {
      console.error(err);
    });
});

// 处理文件下载事件（示例）
const handleDownload = (file) => {
  console.log("Dashboard 收到下载请求:", file.name);
  // 在这里可以调用实际的后端接口或 Wails 函数来处理下载
  alert(`请求下载: ${file.name}`);
};

// 处理文件上传事件
const handleUpload = async () => { // 使用 async 关键字
  console.log("Dashboard 请求上传文件");
  try {
    // 调用 Wails Runtime 打开文件选择对话框
    // 可以配置选项，例如允许多选、设置标题、过滤器等
    const selectedFiles = await UploadFile();

    if (selectedFiles) {
      // 如果用户选择了文件 (selectedFiles 可能是一个包含路径的数组或单个路径字符串)
      console.log("选择的文件:", selectedFiles);
      // 这里可以接续处理上传逻辑，例如调用后端 API 或 Wails Go 函数
      // 如果允许多选，selectedFiles 是一个数组；否则是单个字符串
      // if (Array.isArray(selectedFiles)) {
      //    alert(`选择了 ${selectedFiles.length} 个文件:\n${selectedFiles.join('\n')}`);
      // } else {
      //    alert(`选择了文件: ${selectedFiles}`);
      // }
    } else {
      // 用户取消了选择
      console.log("用户取消了文件选择");
    }
  } catch (error) {
    // 处理可能发生的错误
    console.error("打开文件对话框时出错:", error);
    alert("打开文件选择器失败");
  }
};
</script>

<style scoped>
.dashboard-container {
  /* padding: 20px; */
  display: flex; /* 使用 flex 布局 */
  flex-direction: column; /* 垂直排列 */
  height: 100%; /* 确保容器占满父容器高度 */
  box-sizing: border-box; /* padding 不会撑大容器 */
}

.action-buttons {
  margin-bottom: 20px; /* 按钮区域和文件列表之间的间距 */
  flex-shrink: 0; /* 防止按钮区域被压缩 */
}

.upload-button {
  padding: 8px 15px;
  background-color: #28a745; /* 绿色背景 */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9em;
  transition: background-color 0.2s;
}

.upload-button:hover {
  background-color: #218838; /* 悬停时深一点的绿色 */
}

/* 让 File 组件占据剩余空间并可滚动 */
.dashboard-container > :last-child { /* 选中 File 组件 */
  flex-grow: 1; /* 占据剩余空间 */
  overflow-y: auto; /* 内容超出时显示滚动条 */
}

h1 {
  color: #333;
  margin-bottom: 30px;
  flex-shrink: 0; /* 防止标题被压缩 */
}
</style>
