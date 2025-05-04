<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import {
  OpenImageDialog,
  GetImageDimensions,
  ResizeWindow,
  RestoreWindowSize,
  GetImageBase64,
} from "../wailsjs/go/main/App";
import { Quit } from "../wailsjs/runtime/runtime";

// State variables
const imageSrc = ref("");
const isViewingImage = ref(false);
const isDragging = ref(false);
const startPos = ref({ x: 0, y: 0 });

// Add keyboard event listener when component is mounted
onMounted(() => {
  document.addEventListener("keydown", handleKeyPress);
});

// Remove event listener when component is unmounted
onBeforeUnmount(() => {
  document.removeEventListener("keydown", handleKeyPress);
});

// Select image
async function selectImage() {
  try {
    const filePath = await OpenImageDialog();
    if (!filePath) return;

    // Get image dimensions
    const dimensions = await GetImageDimensions(filePath);
    if (dimensions.width <= 0 || dimensions.height <= 0) {
      console.error("Unable to get image dimensions");
      return;
    }

    console.log(`Image dimensions: ${dimensions.width}x${dimensions.height}`);

    // Resize window to match image dimensions
    await ResizeWindow(dimensions.width, dimensions.height);

    // Get image data and display
    const base64Image = await GetImageBase64(filePath);
    if (base64Image) {
      imageSrc.value = base64Image;
      isViewingImage.value = true;
    }
  } catch (error) {
    console.error("Error processing image:", error);
  }
}

// Return to selection screen
function backToSelection() {
  isViewingImage.value = false;
  imageSrc.value = "";
  RestoreWindowSize();
}

// Handle keyboard events
function handleKeyPress(event: KeyboardEvent) {
  if (event.key === "Escape" && isViewingImage.value) {
    backToSelection();
  }
}

// Handle window drag start
function startDrag(event: MouseEvent) {
  if (isViewingImage.value) return; // Disable dragging when viewing image

  isDragging.value = true;
  startPos.value = { x: event.clientX, y: event.clientY };

  document.addEventListener("mousemove", onDrag);
  document.addEventListener("mouseup", stopDrag);
}

// Handle window dragging
function onDrag(event: MouseEvent) {
  if (!isDragging.value) return;

  const { runtime } = window as any;
  if (runtime && runtime.WindowMove) {
    runtime.WindowMove();
  }
}

// Stop window dragging
function stopDrag() {
  isDragging.value = false;
  document.removeEventListener("mousemove", onDrag);
  document.removeEventListener("mouseup", stopDrag);
}

// Exit application
function exit() {
  Quit();
}
</script>

<template>
  <div class="container" style="--wails-draggable: drag">
    <!-- Image selection screen -->
    <div v-if="!isViewingImage" class="selection-screen">
      <div class="title-bar">
        <div class="title">Image Window</div>
        <div class="close-btn" style="--wails-draggable: no-drag" @click="exit">
          ✕
        </div>
      </div>
      <div class="content" style="--wails-draggable: no-drag">
        <div class="icon">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="64"
            height="64"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <circle cx="8.5" cy="8.5" r="1.5"></circle>
            <polyline points="21 15 16 10 5 21"></polyline>
          </svg>
        </div>

        <button class="select-btn" @click="selectImage">Select Image</button>

        <div class="hint">Press ESC to return after selecting an image</div>
      </div>
    </div>

    <!-- Image viewing screen -->
    <div v-else class="image-viewer">
      <img
        :src="imageSrc"
        alt="Selected Image"
        class="full-image"
        style="--wails-draggable: drag"
      />
    </div>
  </div>
</template>

<style scoped>
.container {
  height: 100%;
  width: 100%;
  overflow: hidden;
  font-family: "Segoe UI", sans-serif;
  position: relative;
  cursor: move;
}

/* Selection screen styles */
.selection-screen {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e293b;
  color: #fff;
  border-radius: 8px; /* Rounded corners */
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.title-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #0f172a;
  -webkit-app-region: drag;
}

.title {
  font-size: 14px;
  font-weight: 500;
}

.close-btn {
  -webkit-app-region: no-drag;
  cursor: pointer;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.close-btn:hover {
  background-color: #f43f5e;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.icon {
  color: #94a3b8;
  margin-bottom: 20px;
}

.select-btn {
  -webkit-app-region: no-drag;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 16px;
}

.select-btn:hover {
  background-color: #2563eb;
}

.hint {
  font-size: 14px;
  color: #94a3b8;
  text-align: center;
}

/* Image viewer styles */
.image-viewer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #000;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.full-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>

<style>
html,
body,
#app {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
  background: transparent;
}

* {
  box-sizing: border-box;
}
</style>
