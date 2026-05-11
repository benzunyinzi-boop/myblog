<template>
  <section class="section-gap container photo-page">
    <header class="section-headline">
      <div>
        <span class="section-chip">visuals / lens</span>
        <h1 class="page-title">透过取景器</h1>
        <p class="page-copy">在写代码的间隙看世界。用相机抓住那些无法被系统日志记录的瞬间。</p>
      </div>
    </header>

    <div class="photo-filters">
      <button
        v-for="f in photoFilters"
        :key="f.value"
        class="photo-filter"
        :class="{ active: current === f.value }"
        @click="current = f.value"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="visible.length" class="photo-masonry">
      <figure v-for="p in visible" :key="p.id" class="photo-item">
        <img :src="p.src" :alt="p.title" loading="lazy" />
        <figcaption class="photo-meta">
          <strong>{{ p.title }}</strong>
          <span>{{ p.exif }} · {{ p.takenAt }}</span>
        </figcaption>
        <span class="film-hole film-hole-tl" />
        <span class="film-hole film-hole-tr" />
        <span class="film-hole film-hole-bl" />
        <span class="film-hole film-hole-br" />
      </figure>
    </div>
    <div v-else class="state-box glass-card">这个分类下还没有照片。</div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { photoFilters, photos, type Photo, type PhotoTag } from '../../data/photos'

const current = ref<PhotoTag | 'all'>('all')

const visible = computed<Photo[]>(() =>
  current.value === 'all' ? photos : photos.filter((p) => p.tag === current.value)
)
</script>
